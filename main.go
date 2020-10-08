package main

import (
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		_ "github.com/weblfe/travel-app/routers"
)

func main() {
		logs.Info("api server start....")
		beego.Run()
}