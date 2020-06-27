package main

import (
		"github.com/astaxie/beego"
		_ "github.com/astaxie/beego/session/redis"
		_ "github.com/go-sql-driver/mysql"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
)

func init() {
		bootstrap()
}

// 引导逻辑
func bootstrap() {
		// swagger
		initSwagger()
		// session
		initSession()
		// database
		initDatabase()
		// middleware
		initMiddleware()
}

// 配置 session
func initSession() {
		if libs.InArray(beego.AppConfig.String("session_on"), []string{"on", "1", "true", "yes"}) {
				beego.BConfig.WebConfig.Session.SessionOn = true
				beego.BConfig.WebConfig.Session.SessionProvider = beego.AppConfig.DefaultString("session_driver", "redis")
				beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.DefaultString("session_config", "127.0.0.1:6379")
		}
}

// 初始 swagger
func initSwagger() {
		if beego.BConfig.RunMode == "dev" {
				beego.BConfig.WebConfig.DirectoryIndex = true
				beego.BConfig.WebConfig.StaticDir["/static/swagger"] = "swagger"
		}
}

// 初始化数据库
func initDatabase() {
		mode := beego.BConfig.RunMode
		if database, err := beego.AppConfig.GetSection(mode + ".database"); err == nil {
				if driver, ok := database["db_driver"]; ok && driver == "mongodb" {
						initMongodb(database)
				}
		}

}

// 初始 mongodb
func initMongodb(data map[string]string) {
		for key, v := range data {
				models.SetProfile(key, v)
		}
}

// 初始middleware
func initMiddleware()  {
		 manger:=middlewares.GetMiddlewareManger()
		 // 注册路由中间件
		 manger.Router(middlewares.AuthMiddlewareName,"/user/info",beego.BeforeRouter)
		 // 启用中间
		 manger.Boot()
}