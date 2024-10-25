package CoreUtils

import (
	"fmt"
	"strconv"
	"time"
)

func TimeStamp() int64 {
	return time.Now().UnixMilli()
}
func TimeStampString() string {
	return strconv.FormatInt(TimeStamp(), 10)
}

func MillisecondsToTime(ts int64) string {
	ms := ts % 1000
	ts -= ms
	ts /= 1000
	seconds := ts % 60
	ts -= seconds
	ts /= 60
	minutes := ts % 60
	ts -= minutes
	ts /= 60
	hours := ts % 24
	ts -= hours
	days := ts / 24
	//hours := ts / 60
	return fmt.Sprintf("%2d Day %2d Hour %2d Minutes %2d Second %3d MillSecond", days, hours, minutes, seconds, ms)
}
