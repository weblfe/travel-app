package libs

import "time"

const (
		TIMEFORMAT       = "20060102150405"
		NORMALTIMEFORMAT = "2006-01-02 15:04:05"
)

// GetTime 当前时间
func GetTime() time.Time {
		return time.Now()
}

// GetTimeString 格式化为:20060102150405
func GetTimeString(t time.Time) string {
		return t.Format(TIMEFORMAT)
}

// GetNormalTimeString 格式化为:2006-01-02 15:04:05
func GetNormalTimeString(t time.Time) string {
		return t.Format(NORMALTIMEFORMAT)
}

// GetTimeUnix 转为时间戳->秒数
func GetTimeUnix(t time.Time) int64 {
		return t.Unix()
}

// GetTimeMills 转为时间戳->毫秒数
func GetTimeMills(t time.Time) int64 {
		return t.UnixNano() / 1e6
}

// GetTimeByInt 时间戳转时间
func GetTimeByInt(t1 int64) time.Time {
		return time.Unix(t1, 0)
}

// GetTimeByString 字符串转时间
func GetTimeByString(timeStr string) (time.Time, error) {
		if timeStr == "" {
				return time.Time{}, nil
		}
		return time.ParseInLocation(TIMEFORMAT, timeStr, time.Local)
}

// GetTimeByNormalString 标准字符串转时间
func GetTimeByNormalString(timeStr string) (time.Time, error) {
		if timeStr == "" {
				return time.Time{}, nil
		}
		return time.ParseInLocation(NORMALTIMEFORMAT, timeStr, time.Local)
}

// CompareTime 比较两个时间大小
func CompareTime(t1, t2 time.Time) bool {
		return t1.Before(t2)
}

// GetNextHourTime n小时后的时间字符串
func GetNextHourTime(s string, n int64) string {
		t2, _ := time.ParseInLocation(TIMEFORMAT, s, time.Local)
		t1 := t2.Add(time.Hour * time.Duration(n))
		return GetTimeString(t1)
}

// GetHourDiffer 计算俩个时间差多少小时
func GetHourDiffer(startTime, endTime string) float32 {
		var hour float32
		t1, err := time.ParseInLocation(TIMEFORMAT, startTime, time.Local)
		t2, err := time.ParseInLocation(TIMEFORMAT, endTime, time.Local)
		if err == nil && CompareTime(t1, t2) {
				diff := GetTimeUnix(t2) - GetTimeUnix(t1)
				hour = float32(diff) / 3600
				return hour
		}
		return hour
}

// CheckHours 判断当前时间是否是整点
func CheckHours() bool {
		_, m, s := GetTime().Clock()
		if m == s && m == 0 && s == 0 {
				return true
		}
		return false
}
