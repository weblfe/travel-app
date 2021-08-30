package common

import (
		"errors"
		"strconv"
		"strings"
)

// Version 版本号协议 ：https://semver.org/lang/zh-CN/
//  X.Y.Z
// 主版本号、次版本号及修订号以数值比较，例如：1.0.0 < 2.0.0 < 2.1.0 < 2.1.1
// 范例：1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0。
// <major> "." <minor> "." <patch>
type Version string // 版本号

func (this Version) Check(ver string, compare string) bool {
		if ver == "" || compare == "" {
				return false
		}
		if v, err := this.toIntVer(ver); err == nil {
				return this.compare(v, compare)
		}
		return false
}

func (this Version) toIntVer(ver string) (int, error) {
		var (
				arr = strings.SplitN(ver, ".", -1)
				num = len(arr)
		)
		if num < 3 {
				for ; num < 3; num++ {
						arr = append(arr, "0")
				}
		}
		var (
				major, minor, patch          = []rune(arr[0]), []rune(arr[1]), []rune(arr[2])
				majorLen, minorLen, patchLen = len(major), len(minor), len(patch)
		)
		if majorLen < 4 {
				major = append([]rune(strings.Repeat("0", 4-majorLen)), major...)
		}
		if minorLen < 2 {
				minor = append(minor, []rune(strings.Repeat("0", 2-majorLen))...)
		}
		if patchLen < 2 {
				patch = append(patch, []rune(strings.Repeat("0", 2-patchLen))...)
		}
		var majorStr, minorStr, patchStr = string(major), string(minor), string(patch)
		newVer := strings.Join([]string{majorStr, minorStr, patchStr}, "")
		if newVer != "" {
				return strconv.Atoi(newVer)
		}
		return 0, errors.New("error version")
}

func (this Version) IsVersion() bool {
		if n, err := this.toIntVer(string(this)); err == nil && n > 0 {
				return true
		}
		return false
}

func (this Version) compare(ver int, compare string) bool {
		var (
				v1, err = this.toIntVer(string(this))
		)
		if err != nil || v1 == 0 || compare == "" {
				return false
		}
		switch compare {
		case "==":
				if v1 == ver {
						return true
				}
		case "!=":
				if v1 != ver {
						return true
				}
		case ">":
				if v1 > ver {
						return true
				}
		case ">=":
				if v1 >= ver {
						return true
				}
		case "<":
				if v1 < ver {
						return true
				}
		case "<=":
				if v1 <= ver {
						return true
				}
		}
		return false
}
