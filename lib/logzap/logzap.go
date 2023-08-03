package logzap

import (
	"github.com/pkg/errors"
	"gnet/lib/utils"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"gnet/game/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// DebugLevel level
	DebugLevel = zapcore.Level(zap.DebugLevel)
	// InfoLevel level
	InfoLevel = zapcore.Level(zap.InfoLevel)
	// WarnLevel level
	WarnLevel = zapcore.Level(zap.WarnLevel)
	// ErrorLevel level
	ErrorLevel = zapcore.Level(zap.ErrorLevel)
	// PanicLevel level
	PanicLevel = zapcore.Level(zap.PanicLevel)
	// FatalLevel level
	FatalLevel = zapcore.Level(zap.FatalLevel)
)

var (
	cfg      *zap.Config
	logger   *zap.Logger
	sugar    *zap.SugaredLogger
	source   string
	curLevel zapcore.Level // 日志等级
	display  bool          // 终端是否打印
)

// 日志未按日期、大小分文件打印
//func init() {
//	if conf.Debug {
//		curLevel = DebugLevel
//		cfg = &zap.Config{
//			Level:       zap.NewAtomicLevelAt(DebugLevel),
//			Development: true,
//			Sampling: &zap.SamplingConfig{
//				Initial:    100,
//				Thereafter: 100,
//			},
//			Encoding: "json",
//			EncoderConfig: zapcore.EncoderConfig{
//				TimeKey:     "ts",
//				LevelKey:    "lv",
//				NameKey:     "logger",
//				CallerKey:   "caller",
//				FunctionKey: zapcore.OmitKey,
//				MessageKey:  "msg",
//				//StacktraceKey:  "stack",
//				LineEnding:     zapcore.DefaultLineEnding,
//				EncodeLevel:    zapcore.LowercaseLevelEncoder,
//				EncodeTime:     zapcore.EpochTimeEncoder,
//				EncodeDuration: zapcore.StringDurationEncoder,
//				EncodeCaller:   zapcore.ShortCallerEncoder,
//			},
//			OutputPaths:      []string{"stderr"},
//			ErrorOutputPaths: []string{"stderr"},
//		}
//	} else {
//		curLevel = InfoLevel
//		cfg = &zap.Config{
//			Level:       zap.NewAtomicLevelAt(InfoLevel),
//			Development: false,
//			Sampling: &zap.SamplingConfig{
//				Initial:    100,
//				Thereafter: 100,
//			},
//			Encoding: "json",
//			EncoderConfig: zapcore.EncoderConfig{
//				//TimeKey:     "ts",
//				LevelKey:    "lv",
//				NameKey:     "logger",
//				CallerKey:   "caller",
//				FunctionKey: zapcore.OmitKey,
//				MessageKey:  "msg",
//				//StacktraceKey:  "stack",
//				LineEnding:     zapcore.DefaultLineEnding,
//				EncodeLevel:    zapcore.LowercaseLevelEncoder,
//				EncodeTime:     zapcore.EpochTimeEncoder,
//				EncodeDuration: zapcore.SecondsDurationEncoder,
//				EncodeCaller:   zapcore.ShortCallerEncoder,
//			},
//			OutputPaths:      []string{"stderr"},
//			ErrorOutputPaths: []string{"stderr"},
//		}
//	}
//	buildLogger()
//}

// 日志文件按日期、大小分文件滚动打印, 日志文件保留7天
func init() {
	fileDir := conf.LogsConf.FileDir
	if fileDir[0] == '.' {
		fileDir = filepath.Join(utils.GetCurDir(), fileDir) // '.'开头认为是相对路径
	}
	if fileDir != "" {
		err := utils.CreateDir(fileDir)
		if err != nil {
			panic(errors.Wrap(err, `logzap init err`))
		}
	}
	fileName := conf.LogsConf.FileName + ".%Y%m%d"
	maxSize := conf.LogsConf.MaxSize
	if maxSize <= 0 {
		panic(errors.New(`logzap init err maxSize`))
	}
	maxAge := conf.LogsConf.MaxAge
	if maxAge <= 0 {
		panic(errors.New(`logzap init err maxAge`))
	}
	rotateTime := time.Hour * 24
	if conf.Debug {
		curLevel = DebugLevel
		display = true
	} else {
		curLevel = InfoLevel
		display = false
	}

	cfg = &zap.Config{
		Level:       zap.NewAtomicLevelAt(curLevel),
		Development: conf.Debug,
		Sampling: &zap.SamplingConfig{
			Initial:    20,
			Thereafter: 10,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "ts",
			LevelKey:    "lv",
			NameKey:     "logger",
			CallerKey:   "caller",
			FunctionKey: zapcore.OmitKey,
			MessageKey:  "msg",
			//StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	encoder := zapcore.NewJSONEncoder(cfg.EncoderConfig)

	var opts []zap.Option
	if cfg.Development {
		opts = append(opts, zap.Development())
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(ErrorLevel))
	}
	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}
	if cell := cfg.Sampling; cell != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if cell.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(cell.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				cfg.Sampling.Initial,
				cfg.Sampling.Thereafter,
				samplerOpts...,
			)
		}))
	}

	var writers []zapcore.WriteSyncer
	if display {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	var divHook io.Writer
	divHook = newDivWriter(fileDir, fileName, maxSize, maxAge, rotateTime)
	writers = append(writers, zapcore.AddSync(divHook))

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), curLevel),
	)
	logger = zap.New(core, opts...)
	if source != "" {
		logger = logger.With(zap.String("source", source))
	}
	sugar = logger.Sugar()
}

func buildLogger() {
	if newLogger, err := cfg.Build(); err == nil {
		if logger != nil {
			logger.Sync()
		}
		logger = newLogger
		if source != "" {
			logger = logger.With(zap.String("source", source))
		}
		sugar = logger.Sugar()
	} else {
		panic(err)
	}
}

// SetSource sets the component name (dispatcher/gate/game) of module
func SetSource(source_ string) {
	source = source_
	buildLogger()
}

// SetLevel sets the log level
func SetLevel(lv zapcore.Level) {
	curLevel = lv
	cfg.Level.SetLevel(lv)
}

// GetLevel get the current log level
func GetLevel() zapcore.Level {
	return curLevel
}

// TraceError prints the stack and error
func TraceError(format string, args ...interface{}) {
	Error(string(debug.Stack()))
	Errorf(format, args...)
}

// SetOutput sets the output writer
func SetOutput(outputs []string) {
	cfg.OutputPaths = outputs
	buildLogger()
}

// ParseLevel converts string to Levels
func ParseLevel(s string) zapcore.Level {
	if strings.ToLower(s) == "debug" {
		return DebugLevel
	} else if strings.ToLower(s) == "info" {
		return InfoLevel
	} else if strings.ToLower(s) == "warn" || strings.ToLower(s) == "warning" {
		return WarnLevel
	} else if strings.ToLower(s) == "error" {
		return ErrorLevel
	} else if strings.ToLower(s) == "panic" {
		return PanicLevel
	} else if strings.ToLower(s) == "fatal" {
		return FatalLevel
	}
	Errorf("ParseLevel: unknown level: %s", s)
	return DebugLevel
}

func Debugw(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Debugw(format, args...)
}

func Debugf(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Debugf(format, args...)
}

func Infow(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Infow(format, args...)
}

func Infof(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Infof(format, args...)
}

func Warnw(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Warnw(format, args...)
}

func Warnf(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Warnf(format, args...)
}

func Error(args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Error(args...)
}

func Errorw(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Errorw(format, args...)
}

func Errorf(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Errorf(format, args...)
}

func Panic(args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Panic(args...)
}

func Panicw(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Panicw(format, args...)
}

func Panicf(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Fatal(args...)
}

func Fatalw(format string, args ...interface{}) {
	debug.PrintStack()
	sugar.With(zap.Time("ts", time.Now())).Fatalw(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	debug.PrintStack()
	sugar.With(zap.Time("ts", time.Now())).Fatalf(format, args...)
}
