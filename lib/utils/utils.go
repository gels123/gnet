/*
 * 通用函数
 */
package utils

import (
	"gnet/lib/logzap"
	"regexp"
	"runtime/debug"

	"go.uber.org/zap"
)

var (
	panicReg *regexp.Regexp
	logReg   *regexp.Regexp
)

func init() {
	panicReg = regexp.MustCompile(`(?m)^panic\(.*\)$`)
	logReg = regexp.MustCompile(`(?m)^lib/log/log\.go:.*$`)
}

func PanicWhen(b bool, s string) {
	if b {
		panic(s)
	}
}

// GetStack return calling stack as string in which messages before log and panic are tripped.
func GetStack() string {
	stack := string(debug.Stack())
	for {
		find := panicReg.FindStringIndex(stack)
		if find == nil {
			break
		}
		stack = stack[find[1]:]
	}
	for {
		find := logReg.FindStringIndex(stack)
		if find == nil {
			break
		}
		stack = stack[find[1]:]
	}
	return stack
}

// CatchPanic calls a function and returns the error if function paniced
func CatchPanic(f func()) (err interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			logzap.Errorf("%v panic, stack: %v\n, %v", f, GetStack(), err)
		}
	}()
	f()
	return
}

// RunPanicless calls a function panic-freely
func RunPanicless(f func()) (panicless bool) {
	defer func() {
		err := recover()
		panicless = (err == nil)
		if err != nil {
			logzap.Errorf("%v panic, stack: %v\n, %v", f, GetStack(), err)
		}
	}()
	f()
	return
}

// RepeatUntilPanicless runs the function repeatly until there is no panic
func RepeatUntilPanicless(f func()) {
	for !RunPanicless(f) {
	}
}

// NextLargerKey finds the next key that is larger than the specified key,
// but smaller than any other keys that is larger than the specified key
func NextLargerKey(key string) string {
	return key + "\x00" // the next string that is larger than key, but smaller than any other keys > key
}

// 安全启动协程执行任务
func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logzap.Error("recover error", zap.String("stack", GetStack()))
			}
			return
		}()
		f()
	}()
}
