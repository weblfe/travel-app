package main

import (
	"github.com/weblfe/travel-app/bootstrap"
	_ "github.com/weblfe/travel-app/routers"
)

func main() {
	bootstrap.StartUp()
	bootstrap.Run()
}
