package helper

import (
	"regexp"
	"runtime/debug"
)

var (
	panicReg *regexp.Regexp
	logReg   *regexp.Regexp
)

func init() {
	regStr := `(?m)^panic\(.*\)$`
	panicReg = regexp.MustCompile(regStr)

	logReg = regexp.MustCompile(`(?m)^.*lotou/log/log\.go:.*$`)
}

func PanicWhen(b bool, s string) {
	if b {
		panic(s)
	}
}

//GetStack return calling stack as string in which messages before log and panic are tripped.
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
