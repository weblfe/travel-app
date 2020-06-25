package libs

import (
		"fmt"
		"io/ioutil"
		"os"
)

func HashCode(any interface{}) string {
		return fmt.Sprintf("%v", any)
}

func FileHash(file string) string {
		if !IsExits(file) {
				return ""
		}
		if data, err := ioutil.ReadFile(file); err == nil {
				return HashCode(data)
		}
		return ""
}

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
