package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/weblfe/travel-app/bootstrap"
	_ "github.com/weblfe/travel-app/routers"
	"os"
)

func main() {
	logs.Info("api server start....")
	path, _ := os.Getwd()
	bootstrap.StartUp(path)
	beego.Run()
}
