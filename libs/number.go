package libs

import (
		"fmt"
		"math"
		"reflect"
		"strconv"
		"strings"
)

const (
		BigNumberK      = 1000
		BigNumberKUnit  = "k"
		BigNumberW      = BigNumberK * 10
		BigNumberWUnit  = "w"
		BigNumberKW     = BigNumberW * BigNumberK
		BigNumberKWUnit = "kw"
		BigNumberWW     = BigNumberKW * BigNumberK
		BigNumberWWUnit = "ww"
)

var (
		BigNumberMapper = []struct {
				Limiter int64
				Unit    string
		}{
				{BigNumberK, BigNumberKUnit},
				{BigNumberW, BigNumberWUnit},
				{BigNumberKW, BigNumberKWUnit},
				{BigNumberWW, BigNumberWWUnit},
		}
)

func IsNumber(v interface{}) bool {
		var (
				getValue = reflect.ValueOf(v)
				kindName = getValue.Kind()
		)
		switch kindName {
		case reflect.Int:
				fallthrough
		case reflect.Int8:
				fallthrough
		case reflect.Int16:
				fallthrough
		case reflect.Int32:
				fallthrough
		case reflect.Int64:
				fallthrough
		case reflect.Uint:
				fallthrough
		case reflect.Uint8:
				fallthrough
		case reflect.Uint16:
				fallthrough
		case reflect.Uint32:
				fallthrough
		case reflect.Uint64:
				fallthrough
		case reflect.Uintptr:
				fallthrough
		case reflect.Float32:
				fallthrough
		case reflect.Float64:
				return true
		}
		return false
}

// 大数格式化
func BigNumberStringer(num int64) string {
		var (
				times      int64 = 0
				before     int64 = BigNumberK
				beforeUnit       = BigNumberKUnit
				total            = len(BigNumberMapper)
				n                = int64(math.Abs(float64(num)))
		)
		for i, mapper := range BigNumberMapper {
				times = 0
				if total > i+1 {
						times = before / BigNumberMapper[i+1].Limiter
				}
				if times > 1 && n/times > 1 {
						continue
				}
				if mapper.Limiter >= n && n/mapper.Limiter >= 1 {
						before = mapper.Limiter
						beforeUnit = mapper.Unit
						break
				}

		}
		if before == BigNumberK {
				if BigNumberK/2 > num {
						return fmt.Sprintf("%d", num)
				}
				if BigNumberK == num {
						return fmt.Sprintf("%d%s", int(float64(num)/float64(before)), beforeUnit)
				}
				return fmt.Sprintf("%s%s", DecimalText(float64(num)/float64(before)), beforeUnit)
		}
		return fmt.Sprintf("%d%s", num/before, beforeUnit)
}

// 保留 小数位数
func DecimalText(f float64, n ...int) string {
		if len(n) == 0 {
				n = append(n, 2)
		}
		if n[0] < 0 {
				n[0] = 2
		}
		var (
				str  = fmt.Sprintf("%v", Decimal(f*math.Pow10(n[0]), n...))
				nums = strings.SplitN(str, ".", -1)
		)
		if len(nums) >= 1 {
				inst, err := strconv.ParseFloat(nums[0], 64)
				if err != nil {
						return "0." + strings.Repeat("0", n[0])
				}
				inst = inst / math.Pow10(n[0])
				return fmt.Sprintf("%v", inst)
		}
		return "0.000"
}

func Decimal(f float64, n ...int) float64 {
		if len(n) == 0 {
				n = append(n, 2)
		}
		if n[0] < 0 {
				n[0] = 2
		}
		n10 := math.Pow10(n[0])
		return math.Trunc((f+0.5/n10)*n10) / n10
}
