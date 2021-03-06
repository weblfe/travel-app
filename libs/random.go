package libs

import (
		"fmt"
		"math"
		"math/rand"
		"time"
)

const (
		randNumbers = "0123456789"
		randWords   = "abcdefghijklmnopqrstuvwsyzABCDEFGHIJKLMNOPQRSTUVWSYZ"
		randZhWords = "的一了是我不在人们有来他这上着个地到大里说去子得也和那要下看天时过出小么起你都把好还多没为又可家学只以主会样年想能生同老中从自面前头到它后然走很轻讲农古黑告界拉名呀"
)

func RandomNumLimitN(count int) string {
		return RandomWords(count, randNumbers)
}

func RandomZhWords(count int) string {
		return RandomWords(count, randZhWords)
}

func RandomAnyWord(count int) string {
		return RandomWords(count, randZhWords+randWords+randNumbers)
}

func RandomWord(count int) string {
		return RandomWords(count, randNumbers+randWords)
}

func RandomWords(count int, words string) string {
		var (
				str     = ""
				index   = 0
				i       = 0
				randArr = []rune(words)
				size    = len(randArr)
		)
		rand.Seed(time.Now().Unix() + rand.Int63n(1000))
		rand.Shuffle(size, func(i, j int) {
				randArr[i], randArr[j] = randArr[j], randArr[i]
		})
		for ; index < count; index++ {
				if index < size {
						i = index
				} else {
						i = rand.Intn(size - 1)
				}
				str += string(randArr[i])
		}
		return str
}

func Shuffle(arr []interface{}) []interface{} {
		if len(arr) == 0 {
				return arr
		}
		for i := len(arr) - 1; i > 0; i-- {
				num := rand.Intn(i + 1)
				arr[i], arr[num] = arr[num], arr[i]
		}
		return arr
}

func RandInt(min, max int) int {
		rand.Seed(time.Now().UnixNano())
		min, max = int(math.Max(float64(min), float64(max))), int(math.Min(float64(min), float64(max)))
		var n = rand.Intn(max) + min
		if n > max {
				return max
		}
		return n
}

func RandFloat64(min, max float64) float64 {
		rand.Seed(time.Now().UnixNano())
		min, max = math.Max(min, max), math.Min(min, max)
		var f = rand.Float64()
		if f < min {
				return f + min
		}
		var num = f*(max-1) + f
		if num > max {
				return max
		}
		return num
}

func RandNumbers(len int) string {
		if len <= 0 {
				return ""
		}
		var (
				size = fmt.Sprintf("%v", len)
				tpl  = "%0" + size + "v"
		)
		return fmt.Sprintf(tpl, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(int32(math.Pow10(len))))
}
