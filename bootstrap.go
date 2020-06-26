package main

import (
		"github.com/astaxie/beego"
		_ "github.com/astaxie/beego/session/redis"
		_ "github.com/go-sql-driver/mysql"
		"github.com/weblfe/travel-app/libs"
)

func init() {
		bootstrap()
}

func bootstrap() {
		// swagger
		initSwagger()
		// session
		initSession()
}

func initSession()  {
		if libs.InArray(beego.AppConfig.String("session_on"), []string{"on", "1", "true", "yes"}) {
				beego.BConfig.WebConfig.Session.SessionOn = true
				beego.BConfig.WebConfig.Session.SessionProvider = beego.AppConfig.DefaultString("session_driver","redis")
				beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.DefaultString("session_config","127.0.0.1:6379")
		}
}

func initSwagger()  {
		if beego.BConfig.RunMode == "dev" {
				beego.BConfig.WebConfig.DirectoryIndex = true
				beego.BConfig.WebConfig.StaticDir["/static/swagger"] = "swagger"
		}
}
