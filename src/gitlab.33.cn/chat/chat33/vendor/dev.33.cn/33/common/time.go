package common

import (
	"time"
)

var loc = local()

func NowMillionSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func NowSecond() int64 {
	return time.Now().Unix()
}

func local() *time.Location {
	loc, _ := time.LoadLocation("Asia/Chongqing")
	return loc
}

func ToCstTime(layout string, timeStr string) time.Time {
	t, _ := time.ParseInLocation(layout, timeStr, loc)
	return t
}

func CstTime(timeStr string) time.Time {
	// 2017-10-26T10:02:56.205Z      UTC time
	// 2017-10-26T17:15:46.711+08:00 CST time
	t, _ := time.Parse(time.RFC3339, timeStr) // UTC
	return t.In(loc)
}

func Time2YYMMDDhhmmss(t time.Time) string {
	return t.In(loc).Format("2006-01-02 15:04:05")
}

func Sec2YYMMDDhhmmss(sec int64) string {
	return time.Unix(sec, 0).In(loc).Format("2006-01-02 15:04:05")
}

func Sec2Time(sec int64) time.Time {
	return time.Unix(sec, 0).In(loc)
}

func MillionSecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func BeginSecOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndSecOfDay(t time.Time) time.Time {
	bt := BeginSecOfDay(t)
	return bt.Add(time.Hour * 24).Add(-time.Second)
}
