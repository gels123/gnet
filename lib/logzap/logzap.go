/*
日志打印   eg. 语法糖模式=>logzap.Debugw("=sdfsdf=", "num=", 100) 原生模式(性能更佳)=>logzap.Debug("=sdfsdf=", zap.Bool("ok", true), zap.String("name", "lord1"))
*/
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
	// debugLevel level
	debugLevel = zapcore.Level(zap.DebugLevel)
	// infoLevel level
	infoLevel = zapcore.Level(zap.InfoLevel)
	// warnLevel level
	warnLevel = zapcore.Level(zap.WarnLevel)
	// errorLevel level
	errorLevel = zapcore.Level(zap.ErrorLevel)
	// panicLevel level
	panicLevel = zapcore.Level(zap.PanicLevel)
	// fatalLevel level
	fatalLevel = zapcore.Level(zap.FatalLevel)
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
//
//	func init() {
//		if conf.Debug {
//			curLevel = debugLevel
//			cfg = &zap.Config{
//				Level:       zap.NewAtomicLevelAt(debugLevel),
//				Development: true,
//				Sampling: &zap.SamplingConfig{
//					Initial:    100,
//					Thereafter: 100,
//				},
//				Encoding: "json",
//				EncoderConfig: zapcore.EncoderConfig{
//					TimeKey:     "ts",
//					LevelKey:    "lv",
//					NameKey:     "logger",
//					CallerKey:   "caller",
//					FunctionKey: zapcore.OmitKey,
//					MessageKey:  "msg",
//					//StacktraceKey:  "stack",
//					LineEnding:     zapcore.DefaultLineEnding,
//					EncodeLevel:    zapcore.LowercaseLevelEncoder,
//					EncodeTime:     zapcore.EpochTimeEncoder,
//					EncodeDuration: zapcore.StringDurationEncoder,
//					EncodeCaller:   zapcore.ShortCallerEncoder,
//				},
//				OutputPaths:      []string{"stderr"},
//				ErrorOutputPaths: []string{"stderr"},
//			}
//		} else {
//			curLevel = infoLevel
//			cfg = &zap.Config{
//				Level:       zap.NewAtomicLevelAt(infoLevel),
//				Development: false,
//				Sampling: &zap.SamplingConfig{
//					Initial:    100,
//					Thereafter: 100,
//				},
//				Encoding: "json",
//				EncoderConfig: zapcore.EncoderConfig{
//					//TimeKey:     "ts",
//					LevelKey:    "lv",
//					NameKey:     "logger",
//					CallerKey:   "caller",
//					FunctionKey: zapcore.OmitKey,
//					MessageKey:  "msg",
//					//StacktraceKey:  "stack",
//					LineEnding:     zapcore.DefaultLineEnding,
//					EncodeLevel:    zapcore.LowercaseLevelEncoder,
//					EncodeTime:     zapcore.EpochTimeEncoder,
//					EncodeDuration: zapcore.SecondsDurationEncoder,
//					EncodeCaller:   zapcore.ShortCallerEncoder,
//				},
//				OutputPaths:      []string{"stderr"},
//				ErrorOutputPaths: []string{"stderr"},
//			}
//		}
//		buildLogger()
//	}
//
//func buildLogger() {
//	if newLogger, err := cfg.Build(); err == nil {
//		if logger != nil {
//			logger.Sync()
//		}
//		logger = newLogger
//		if source != "" {
//			logger = logger.With(zap.String("source", source))
//		}
//		sugar = logger.Sugar()
//	} else {
//		panic(err)
//	}
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
		curLevel = debugLevel
		display = true
	} else {
		curLevel = infoLevel
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
			TimeKey:        "ts",
			LevelKey:       "lv",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stack",
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
		opts = append(opts, zap.AddStacktrace(errorLevel))
	}
	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}
	opts = append(opts, zap.AddCallerSkip(1))
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

	var cores []zapcore.Core
	var divHook io.Writer
	divHook = newDivWriter(fileDir, fileName, maxSize, maxAge, rotateTime)
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(divHook), curLevel))
	if display {
		encoderConfigc := cfg.EncoderConfig
		encoderConfigc.EncodeLevel = zapcore.LowercaseColorLevelEncoder
		encoderConfigc.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderc := zapcore.NewConsoleEncoder(encoderConfigc)
		cores = append(cores, zapcore.NewCore(encoderc, zapcore.AddSync(os.Stdout), curLevel))
	}
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, opts...)
	logger.Sync()
	if source != "" {
		logger = logger.With(zap.String("source", source))
	}
	sugar = logger.Sugar()
}

// SetSource sets the component name (dispatcher/gate/game) of module
func SetSource(source_ string) {
	source = source_
	if source != "" && logger != nil {
		logger = logger.With(zap.String("source", source))
		sugar = logger.Sugar()
	}
}

// ParseLevel converts string to Levels
func ParseLevel(s string) zapcore.Level {
	if strings.ToLower(s) == "debug" {
		return debugLevel
	} else if strings.ToLower(s) == "info" {
		return infoLevel
	} else if strings.ToLower(s) == "warn" || strings.ToLower(s) == "warning" {
		return warnLevel
	} else if strings.ToLower(s) == "error" {
		return errorLevel
	} else if strings.ToLower(s) == "panic" {
		return panicLevel
	} else if strings.ToLower(s) == "fatal" {
		return fatalLevel
	}
	Errorf("ParseLevel: unknown level: %s", s)
	return debugLevel
}

// GetLevel get the current log level
func GetLevel() zapcore.Level {
	return curLevel
}

func Debug(msg string, fields ...zap.Field) {
	if curLevel <= debugLevel {
		logger.Debug(msg, fields...)
	}
}

func Debugw(format string, args ...interface{}) {
	if curLevel <= debugLevel {
		sugar.With(zap.Time("ts", time.Now())).Debugw(format, args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if curLevel <= debugLevel {
		sugar.With(zap.Time("ts", time.Now())).Debugf(format, args...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if curLevel <= infoLevel {
		logger.Info(msg, fields...)
	}
}

func Infow(format string, args ...interface{}) {
	if curLevel <= infoLevel {
		sugar.With(zap.Time("ts", time.Now())).Infow(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if curLevel <= infoLevel {
		sugar.With(zap.Time("ts", time.Now())).Infof(format, args...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if curLevel <= warnLevel {
		logger.Warn(msg, fields...)
	}
}

func Warnw(format string, args ...interface{}) {
	if curLevel <= warnLevel {
		sugar.With(zap.Time("ts", time.Now())).Warnw(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if curLevel <= warnLevel {
		sugar.With(zap.Time("ts", time.Now())).Warnf(format, args...)
	}
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Errorw(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Errorw(format, args...)
}

func Errorf(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Errorf(format, args...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

func Panicw(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Panicw(format, args...)
}

func Panicf(format string, args ...interface{}) {
	sugar.With(zap.Time("ts", time.Now())).Panicf(format, args...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Fatalw(format string, args ...interface{}) {
	debug.PrintStack()
	sugar.With(zap.Time("ts", time.Now())).Fatalw(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	debug.PrintStack()
	sugar.With(zap.Time("ts", time.Now())).Fatalf(format, args...)
}

// TraceError prints the stack and error
func TraceError(format string, args ...interface{}) {
	Error(string(debug.Stack()))
	Errorf(format, args...)
}
