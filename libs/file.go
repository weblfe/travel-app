package libs

import "os"

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