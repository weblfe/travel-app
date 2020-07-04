package libs

import (
		"fmt"
		"math"
		"os"
		"strconv"
		"strings"
		"unicode"
)

func IsExits(file string) bool {
		if _, err := os.Stat(file); err != nil {
				if os.IsExist(err) || os.IsNotExist(err) {
						return false
				}
				if os.IsPermission(err) {
						return false
				}
		}
		return true
}

// 文件大小格式化
func FormatFileSize(fileSize int64) string {
		if fileSize <= 0 {
				return "0B"
		}
		for i, unit := range _fileSizeUnitMap {
				if fileSize < int64(math.Pow(1024, float64(i+1))) {
						return fmt.Sprintf("%.2f%s", float64(fileSize)/math.Pow(1024, float64(i)), unit)
				}
		}
		return "xXB"
}

var (
		_fileSizeUnitMap = []string{
				"B", "KB", "MB", "GB", "TB", "EB", "ZB", "YB",
		}
)

type FileSize int64

func (this FileSize) String() string {
		return FormatFileSize(int64(this))
}

func (this FileSize) Parse(size string) int64 {
		var (
				unit string
				data = []rune(size)
				num  = len(data)
		)
		if unicode.IsNumber(data[num-2]) {
				unit = string(data[num-2:])
		}
		if unit == "" && unicode.IsNumber(data[num-1]) {
				unit = string(data[num-1:])
		}
		if unit == "" && unicode.IsNumber(data[num-3]) {
				unit = string(data[num-3:])
		}
		if unit == "" {
				return 0
		}
		return FileSizeTans(strings.Replace(size, unit, "", 1), unit)
}

// 文件大小字符串转换
func FileSizeTans(size string, unit string) int64 {
		num, err := strconv.ParseFloat(size, 64)
		if err != nil {
				return 0
		}
		for i, u := range _fileSizeUnitMap {
				if strings.EqualFold(u, unit) {
						return int64(num * math.Pow(1024, float64(i)))
				}
		}
		return 0
}
