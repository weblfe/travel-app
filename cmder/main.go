package main

import (
	"fmt"
	"github.com/subosito/gotenv"
	"github.com/weblfe/travel-app/cmder/cmd"
	"os"
)

func main() {
	cmd.Execute()
}

func init() {
	if _, err := os.Stat(".env"); err == nil {
		if err = gotenv.Load(".env"); err != nil {
			fmt.Println("error:", err.Error())
		}
	}
}
