package utils

import (
	"time"
)

// 今日毫秒数
func TodayInMillis() int64 {
	now := time.Now()
	y, m, d := now.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, now.Location())
	return date.UnixNano() / int64(time.Millisecond)
}

// 当前时间毫秒数
func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// 天数转为毫秒数
func DaysToMillis(days int) int64 {
	return int64(time.Hour * 24 * time.Duration(days) / time.Millisecond)
}
