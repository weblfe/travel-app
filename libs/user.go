package libs

import (
		"fmt"
		"time"
)

func RandomNickName(lens ...int) string {
		var (
				now  = time.Now()
				size = ArgsInt(lens, 8)
				uuid = fmt.Sprintf("%d%d", now.YearDay(), now.Nanosecond())
				left = size - len(uuid)
		)
		if left <= 0 {
				return RandomWord(size)
		}
		return RandomWord(left) + uuid
}
