package libs_test

import (
	"github.com/weblfe/travel-app/libs"
	"testing"
)

func TestPasswordHash(t *testing.T) {
	var (
		salt   = "71e920133ebb7d0a94b9daed8f6c2d9a"
		pwd    = `18565392186`
		output = libs.PasswordHash(pwd, salt)
	)
	t.Logf("password=%s", output)
}
