package log

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
)

// levels
const (
	DEBUG_LEVEL = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
	FATAL_LEVEL
	LEVEL_MAX
)

const (
	DEBUG_LEVEL_DESC = "[debug] "
	INFO_LEVEL_DESC  = "[info] "
	WARN_LEVEL_DESC  = "[warn] "
	ERROR_LEVEL_DESC = "[error] "
	FATAL_LEVEL_DESC = "[fatal] "
)

type Msg struct {
	level     int
	levelDesc string
	msg       string
}

type Logger interface {
	DoPrintf(level int, levelDesc, msg string)
	SetColored(colored bool)
	Close()
}

var glogger Logger
var gloggerMut sync.Mutex

var HasCallerPos bool

func do(level int, desc, format string, param ...interface{}) {
	if glogger == nil {
		log.Fatal("log is not init, please call log.Init first.")
		return
	}
	m := &Msg{level, desc, fmt.Sprintf(format, param...)}
	gloggerMut.Lock()
	glogger.DoPrintf(m.level, m.levelDesc, m.msg)
	gloggerMut.Unlock()

	if level == FATAL_LEVEL {
		format = desc + format
		panic(fmt.Sprintf(format, param...))
	}
}
func SetLogger(logger Logger) {
	gloggerMut.Lock()
	glogger = logger
	gloggerMut.Unlock()
}

func preProcess(format string, level int) string {
	if level == INFO_LEVEL || level == WARN_LEVEL || level == FATAL_LEVEL {
		return format
	}
	if !HasCallerPos {
		return format
	}
	_, file, line, ok := runtime.Caller(2)
	if ok {
		n := strings.LastIndex(file, "/")
		file = file[n+1:]
		format = fmt.Sprintf("[%v:%v] ", file, line) + format
	}
	return format
}

func Debug(format string, param ...interface{}) {
	format = preProcess(format, DEBUG_LEVEL)
	do(DEBUG_LEVEL, DEBUG_LEVEL_DESC, format, param...)
}

func Info(format string, param ...interface{}) {
	format = preProcess(format, INFO_LEVEL)
	do(INFO_LEVEL, INFO_LEVEL_DESC, format, param...)
}

func Warn(format string, param ...interface{}) {
	format = preProcess(format, WARN_LEVEL)
	do(WARN_LEVEL, WARN_LEVEL_DESC, format, param...)
}

func Error(format string, param ...interface{}) {
	format = preProcess(format, ERROR_LEVEL)
	do(ERROR_LEVEL, ERROR_LEVEL_DESC, format, param...)
}

func Fatal(format string, param ...interface{}) {
	format = preProcess(format, FATAL_LEVEL)
	do(FATAL_LEVEL, FATAL_LEVEL_DESC, format, param...)
}

func Close() {
	if glogger == nil {
		return
	}
	gloggerMut.Lock()
	glogger.Close()
	gloggerMut.Unlock()
}

//init log with SimpleLogger
func Init(path string, filename string, fileLevel, shellLevel, maxLine, bufSize int) Logger {
	logger := CreateLogger(path, filename, fileLevel, shellLevel, maxLine, bufSize)
	SetLogger(logger)
	return logger
}

//init log with default logger
func init() {
	HasCallerPos = true
	SetLogger(DefaultLogger)
}
