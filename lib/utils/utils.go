package utils

import (
	"gowork/lib/helper"
	"gowork/lib/log"
	"time"
)

// CatchPanic calls a function and returns the error if function paniced
func CatchPanic(f func()) (err interface{}) {
	defer func() {
		err := recover()
		if err != nil {
			log.Error("%v panic, stack: %v\n, %v", f, helper.GetStack(), err)
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
			log.Error("%v panic, stack: %v\n, %v", f, helper.GetStack(), err)
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

//获取当前时间戳
func GetTime() int64 {
	return time.Now().Unix()
}

//获取传入(当前)时间的0点的时间戳
func GetTimeH0(timestamp ...int64) int64 {
	if len(timestamp) <= 0 {
		ti := time.Now().Unix()
		return ti - ti%86400
	} else {
		return timestamp[0] - timestamp[0]%86400
	}
}

//获取传入(当前)时间的周一0点的时间戳
func GetTimeW1H0(timestamp ...int64) int64 {
	if len(timestamp) <= 0 {
		ti := time.Now().Unix()
		return ti - ti%(86400*7)
	} else {
		return timestamp[0] - timestamp[0]%(86400*7)
	}
}

//获取传入(当前)时间是周几
func GetTimeWeekDay(timestamp ...int64) int {
	var weekday int
	if len(timestamp) > 0 {
		weekday = int(time.Unix(timestamp[0], 0).Weekday())
	} else {
		weekday = int(time.Now().Weekday())
	}
	if weekday == 0 {
		weekday = 7
	}
	return weekday
}

//获取当前时间的格式化字符串
func GetTimeFormat(layout ...string) string {
	if len(layout) > 0 {
		return time.Now().Format(layout[0])
	} else {
		return time.Now().Format("2006-01-02 15:04:05")
	}
}
