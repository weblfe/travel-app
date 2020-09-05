package libs

import (
		"fmt"
		"math"
		"os"
		"path/filepath"
		"strconv"
		"strings"
		"time"
		"unicode"
)

var (
		// 存储单位表
		_fileSizeUnitMap = []string{
				"B", "KB", "MB", "GB", "TB", "EB", "ZB", "YB",
		}
		// 类型表
		_fileTypeMapper = map[string][]string{
				"image": {
						"png", "jpg", "jpeg", "bmp",
						"pcx", "tif", "gif", "tga",
						"exif", "fpx", "svg", "psd",
						"cdr", "pcd", "dxf", "ufo",
						"eps", "ai", "hdri", "raw",
						"wmf", "flic", "emf", "ico",
				},
				"word": {
						"doc", "docx", "xls",
						"xlsx", "ppt", "pptx",
						"pdf",
				},
				"config": {
						"ini", "yml", "json",
						"yaml", "conf", "toml",
						"xml",
				},
				"code": {
						"html", "xhtml", "java",
						"c", "cpp", "js", "php",
						"lua", "rb", "go", "py",
						"sh", "bat", "cmd", "ps1",
				},
				"video": {
						"mp4", "avi", "mov",
						"rmvb", "rm", "flv",
						"3gp", "mpg", "mpe",
						"mpeg", "wmv", "asf",
						"asx", "wvx", "mpa",
				},
				"audio": {
						"mp3", "cda", "wav",
						"wma", "rm", "mid", "ape", "flac",
				},
				"avatar": {"png", "jpg", "jpeg"},
		}
		// 类型索引顺序
		_fileTypeIndex = []string{"image", "word", "config", "code", "video", "audio", "avatar"}
)

type FileSize int64

func IsExits(file string) bool {
		_, err := os.Stat(file)
		if err != nil {
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
						return fmt.Sprintf("%.3f%s", float64(fileSize)/math.Pow(1024, float64(i)), unit)
				}
		}
		return "xXB"
}

func (this FileSize) String() string {
		return FormatFileSize(int64(this))
}

func (this FileSize) Parse(size string) int64 {
		var (
				unit string
				data = []rune(size)
				num  = len(data)
		)
		if unit == "" && unicode.IsNumber(data[num-1]) {
				unit = string(data[num:])
		}
		if unicode.IsNumber(data[num-2]) {
				unit = string(data[num-1:])
		}
		if unit == "" && unicode.IsNumber(data[num-3]) {
				unit = string(data[num-2:])
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

// 获取文件类型
func GetFileType(file string) string {
		var ext = filepath.Ext(file)
		if ext != "" && strings.Contains(ext, ".") {
				ext = ext[1:]
		}
		for _, ty := range _fileTypeIndex {
				items := _fileTypeMapper[ty]
				for _, it := range items {
						if strings.EqualFold(ext, it) {
								return ty
						}
				}
		}
		return ext
}

// 唯一名文件
func UniqueFile(file string, root ...string) string {
		var ext, base = filepath.Ext(file), filepath.Base(file)
		var name = strings.Replace(base,ext,"",-1)
		if IsExits(file) {
				if len(root) > 0 {
						return filepath.Join(root[0], fmt.Sprintf("%s", unique())+"_"+name+ext)
				}
				return strings.Replace(file, base, unique()+ext, 1)
		}
		file = strings.Replace(file, base, unique()+ext, 1)
		if len(root) == 0 {
				return file
		}
		return filepath.Join(root[0], file)
}

func unique() string {
		return fmt.Sprintf("%v", time.Now().UnixNano())
}
