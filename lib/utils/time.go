package utils

import "time"

var (
	timeZoneShanghai, _ = time.LoadLocation("Asia/Shanghai")
)

func ShanghaiTimeLocation() *time.Location {
	return timeZoneShanghai
}

func TimestampToTime(timestamp int64) time.Time {
	var sec, nsec int64
	if timestamp/1e18 != 0 {
		sec = timestamp / 1e9
		nsec = timestamp % 1e9
	} else if timestamp/1e12 != 0 {
		sec = timestamp / 1e3
		nsec = timestamp % 1e3 * 1e6 //ms to ns
	} else if timestamp/1e9 != 0 {
		sec = timestamp
		nsec = 0
	}
	return time.Unix(sec, nsec)
}
