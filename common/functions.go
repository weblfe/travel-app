package common

import (
		"os"
)

func SmsDebugOn() {
		_ = os.Setenv("sms_debug_on", "1")
}

func SmsDebugOff() {
		_ = os.Setenv("sms_debug_on", "0")
}
