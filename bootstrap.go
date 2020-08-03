package main

import (
		"encoding/gob"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		_ "github.com/astaxie/beego/session/redis"
		_ "github.com/go-sql-driver/mysql"
		"github.com/joho/godotenv"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/plugins"
		"github.com/weblfe/travel-app/services"
		"os"
		"path"
		"strings"
)

func init() {
		bootstrap()
}

// 引导逻辑
func bootstrap() {
		// 环境注册
		initRegisterEnv()
		// 注册结构体
		registerGob()
		// swagger
		initSwagger()
		// session
		initSession()
		// database
		initDatabase()
		// middleware
		initMiddleware()
		// 注册插件
		registerPlugins()

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
		service := services.GetInitDataServiceInstance()
		service.SetInit("./static/database")
		service.Init()
}

// 初始 mongodb
func initMongodb(data map[string]string) {
		for key, v := range data {
				models.SetProfile(key, v)
		}
}

// 数据结构注册
func initRegisterEnv() {
		pwd, _ := os.Getwd()
		envFile := env.Get("ENV_FILE", path.Join(pwd, "/.env"))
		if envFile == "" {
				return
		}
		var arr []string
		envs := strings.SplitN(envFile, ",", -1)

		for _, fs := range envs {
				state, err := os.Stat(fs)
				if err != nil {
						continue
				}
				if state.IsDir() {
						continue
				}
				arr = append(arr, fs)
		}
		if len(arr) < 0 {
				return
		}
		_ = godotenv.Load(arr...)
		// 重新载入env
		reloadEnv()
		// 加载 主配置
		_ = beego.LoadAppConfig("ini", "conf/main.conf")
		// 重新载入配置
		reloadConfig()
}

// 重新载入配置
func reloadConfig() {
		beego.BConfig.RunMode = beego.AppConfig.DefaultString(SetConfGlobalScope("runmode"), beego.BConfig.RunMode)
		beego.BConfig.AppName = beego.AppConfig.DefaultString(SetConfGlobalScope("appname"), beego.BConfig.AppName)
		beego.BConfig.ServerName = beego.AppConfig.DefaultString(SetConfGlobalScope("servername"), beego.BConfig.ServerName)
		beego.BConfig.Listen.HTTPPort = beego.AppConfig.DefaultInt(SetConfGlobalScope("httpport"), beego.BConfig.Listen.HTTPPort)
		beego.BConfig.Listen.HTTPSPort = beego.AppConfig.DefaultInt(SetConfGlobalScope("httpsport"), beego.BConfig.Listen.HTTPSPort)
}

// 重新载入env
func reloadEnv() {
		for _, e := range os.Environ() {
				splits := strings.Split(e, "=")
				env.Set(splits[0], os.Getenv(splits[0]))
		}
}

// 全局
func SetConfGlobalScope(key string) string {
		if strings.Contains(key, "::") {
				return key
		}
		return "default::" + key
}

// 初始middleware
func initMiddleware() {
		manger := middlewares.GetMiddlewareManger()
		// 注册路由中间件
		// 登陆中间
		manger.Router(middlewares.AuthMiddlewareName, "/user/info", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/attachment/*", beego.BeforeRouter)
		manger.Router(middlewares.AttachTicketMiddlewareName, "/attachments/*", beego.BeforeRouter)

		manger.Router(middlewares.TokenMiddleware, "/reset/password", beego.BeforeRouter)
		manger.Router(middlewares.TokenMiddleware, "/posts/lists/my", beego.BeforeRouter)
		manger.Router(middlewares.TokenMiddleware, "/posts/follows", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/posts/create", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/thumbsUp", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/fans", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/comment/create", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/user/friends", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/follow/*", beego.BeforeRouter)

		// 启用中间
		manger.Boot()
}

// 注册结构体
func registerGob() {
		gob.Register(beego.M{})
		gob.Register(models.Tag{})
		gob.Register(models.User{})
		gob.Register(models.Address{})
		gob.Register(models.AppInfo{})
		gob.Register(models.Attachment{})
		gob.Register(models.RequestLog{})
}

// 注册插件
func registerPlugins() {
		plugins.GetNatsPlugin().Register()
}
