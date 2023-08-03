package utils

import (
	"time"
)

// 获取当前时间戳
func GetTime() int64 {
	return time.Now().Unix()
}

// 获取传入(当前)时间的0点的时间戳
func GetTimeH0(timestamp ...int64) int64 {
	if len(timestamp) <= 0 {
		ti := time.Now().Unix()
		return ti - ti%86400
	} else {
		return timestamp[0] - timestamp[0]%86400
	}
}

// 获取传入(当前)时间的周一0点的时间戳
func GetTimeW1H0(timestamp ...int64) int64 {
	if len(timestamp) <= 0 {
		ti := time.Now().Unix()
		return ti - ti%(86400*7)
	} else {
		return timestamp[0] - timestamp[0]%(86400*7)
	}
}

// 获取传入(当前)时间是周几
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

// 获取当前时间的格式化字符串
func GetTimeFormat(layout ...string) string {
	if len(layout) > 0 {
		return time.Now().Format(layout[0])
	} else {
		return time.Now().Format(`2000-01-01 00:00:00`)
	}
}
