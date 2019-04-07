package util

import (
	"time"
)

//TimeConvert 时间操作类
var TimeConvert = timeconvert{}

type timeconvert struct {
}

//GetCurrTime 当前UTC时间
func (timeconvert) GetCurrTime() time.Time {
	return time.Now().UTC()
}

//GetCurrTimeSecond 当前UTC时间精确到秒
func (timeconvert) GetCurrTimeSecond() time.Time {
	result := TimeConvert.GetCurrTime()
	result = time.Date(
		result.Year(),
		result.Month(),
		result.Day(),
		result.Hour(),
		result.Minute(),
		result.Second(),
		0,
		time.UTC)
	return result
}

//GetCurrDate 当前时间的日期
func (timeconvert) GetCurrDate() time.Time {
	result := TimeConvert.GetCurrTime()
	result = time.Date(
		result.Year(),
		result.Month(),
		result.Day(),
		0,
		0,
		0,
		0,
		time.UTC)
	return result
}

func (timeconvert) GetDate(d time.Time) time.Time {
	result := time.Date(
		d.Year(),
		d.Month(),
		d.Day(),
		0,
		0,
		0,
		0,
		time.UTC)
	return result
}
func (timeconvert) GetMinDateTime() time.Time {
	result := time.Date(
		1970,
		1,
		1,
		0,
		0,
		0,
		0,
		time.UTC)
	return result
}
