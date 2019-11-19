package util

import (
	"time"

	"github.com/balinor2017/interview-order/config"
)

var (
	timezoneLoc *time.Location
)

func InitTimeZoneLocation() {
	serverMode := config.MustGetString("server.mode")
	timezoneLocString := config.MustGetString(serverMode + ".timezone")
	SetTimeZoneLocation(timezoneLocString)
}

func GetTimeZoneLocation() *time.Location {
	return timezoneLoc
}

func SetTimeZoneLocation(timezoneLocString string) {
	timezoneLoc, _ = time.LoadLocation(timezoneLocString)
}

func ToTimeString(time time.Time) (resDate string) {
	return time.In(timezoneLoc).Format("2006-01-02 15:04:05")
}
