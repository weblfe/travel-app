package main

import (
		"encoding/gob"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		"github.com/astaxie/beego/logs"
		_ "github.com/astaxie/beego/session/redis"
		_ "github.com/go-sql-driver/mysql"
		"github.com/joho/godotenv"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/plugins"
		"github.com/weblfe/travel-app/services"
		"math"
		"math/rand"
		"os"
		"path"
		"strings"
		"time"
)

// 载入引导逻辑
func init() {
		logs.Info("bootstrap start")
		bootstrap()
		logs.Info("bootstrap end...")
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
		// 注册全局服务
		registerServices()
}

// 配置 session
func initSession() {
		if libs.InArray(beego.AppConfig.String("session_on"), []string{"on", "1", "true", "yes"}) {
				logs.Info("init session")
				beego.BConfig.WebConfig.Session.SessionOn = true
				beego.BConfig.WebConfig.Session.SessionProvider = beego.AppConfig.DefaultString("session_driver", "redis")
				beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.DefaultString("session_config", "127.0.0.1:6379")
		}
}

// 初始 swagger
func initSwagger() {
		if beego.BConfig.RunMode == "dev" {
				logs.Info("init swagger")
				beego.BConfig.WebConfig.DirectoryIndex = true
				beego.BConfig.WebConfig.StaticDir["/static/swagger"] = "swagger"
		}
}

// 初始化数据库
func initDatabase() {
		mode := beego.BConfig.RunMode
		database, err := beego.AppConfig.GetSection(mode + ".database")
		if err == nil {
				if driver, ok := database["db_driver"]; ok && driver == "mongodb" {
						initMongodb(database)
				}
		}
		initMigration()
		logs.Info("init database")
}

// 初始化数据迁移
func initMigration() {
		service := services.GetInitDataServiceInstance()
		service.SetInit("./static/database")
		service.Init()
		logs.Info("init database migration")
}

// 初始 mongodb
func initMongodb(data map[string]string) {
		for key, v := range data {
				models.SetProfile(key, v)
		}
		logs.Info("init database profiles")
		initMongoIndex()
		logs.Info("init database index create")
}

// 初始化 数据索引
func initMongoIndex() {
		models.UserModelOf().CreateIndex(true)
		models.PostsModelOf().CreateIndex(true)
		models.ConfigModelOf().CreateIndex(true)
		models.ThumbsUpModelOf().CreateIndex(true)
		models.SensitiveWordsModelOf().CreateIndex(true)
		models.UserRolesConfigModelOf().CreateIndex(true)
		logs.Info("init database index")
}

// 注册环境变量
func initRegisterEnv() {
		logs.Info("init env")
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

		var (
				on       = env.Get("CONFIGURE_CENTER_ON", "")
				provider = env.Get("CONFIGURE_CENTER_PROVIDER", "etcd")
				appId    = env.Get("CONFIGURE_CENTER_APPID", "travel-app")
		)
		logs.Info("CONFIGURE_CENTER_ON:"+on)
		// 是否加载 配置中心
		if on != "" && on != "0" && on != "false" {
				var manger = plugins.GetConfigureCentreRepositoryMangerInstance()
				manger.InitDef().Boot()
				providerIns := manger.Get(provider)
				if providerIns == nil {
						logs.Error("ConfigureCentre Provider :"+provider+" miss")
						return
				}
				var ticket = time.Second
				for {
						data, err := providerIns.Pull(appId)
						if err != nil {
								logs.Error(err)
						}
						if len(data) != 0 {
								err = services.ConfigServiceOf().Update(data)
								if err != nil {
										logs.Error(err)
								}
								break
						}
						logs.Info("waiting for  travel app config from configure center")
						time.Sleep(ticket)
						if ticket.Seconds() >= time.Hour.Seconds() {
								logs.Error("time out to waiting configure center pull config")
								break
						}
						ticket = ticket + 10 * time.Second
				}
				return
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
		logs.Info("init middleware")
		manger := middlewares.GetMiddlewareManger()
		// 注册路由中间件
		manger.Router(middlewares.CorsMiddlewareName, "*", beego.BeforeExec)
		manger.Router(middlewares.HeaderMiddleWareName, "*", beego.BeforeExec)
		manger.Router(middlewares.HeaderMiddleWareName, "*", beego.AfterExec)
		// 登陆中间
		manger.Router(middlewares.AuthMiddlewareName, "/user/info", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/attachment/*", beego.BeforeRouter)
		manger.Router(middlewares.AttachTicketMiddlewareName, "/attachments/*", beego.BeforeRouter)
		// 检查认证
		manger.Router(middlewares.TokenMiddleware, "/user/profile", beego.BeforeRouter)
		manger.Router(middlewares.TokenMiddleware, "/reset/password", beego.BeforeRouter)
		manger.Router(middlewares.TokenMiddleware, "/posts/lists/my", beego.BeforeRouter)
		manger.Router(middlewares.TokenMiddleware, "/posts/follows", beego.BeforeRouter)
		manger.Router(middlewares.TokenMiddleware, "/posts/follow", beego.BeforeRouter)

		manger.Router(middlewares.AuthMiddlewareName, "/posts/create", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/logout", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/thumbsUp", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/fans", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/comment/create", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/user/friends", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/follow/*", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/follow", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/posts/audit", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/posts/video/cover", beego.BeforeRouter)
		manger.Router(middlewares.AuthMiddlewareName, "/posts/all", beego.BeforeRouter)

		manger.Router(middlewares.RoleMiddleware, "/posts/audit", beego.BeforeRouter)
		manger.Router(middlewares.RoleMiddleware, "/posts/all", beego.BeforeRouter)
		manger.Router(middlewares.RoleMiddleware, "/posts/video/cover", beego.BeforeRouter)

		// 启用中间
		manger.Boot()
}

// 注册结构体
func registerGob() {
		logs.Info("init gob")
		gob.Register(beego.M{})
		gob.Register(models.Tag{})
		gob.Register(models.User{})
		gob.Register(models.Address{})
		gob.Register(models.AppInfo{})
		gob.Register(models.Attachment{})
		gob.Register(models.RequestLog{})
		gob.Register(models.UserRelation{})
		gob.Register(models.PopularizationChannels{})
}

// 注册插件
func registerPlugins() {
		logs.Info("init plugins")
		plugins.GetOSS().Register()
		plugins.GetQrcode().Register()
		plugins.GetNatsPlugin().Register()
		plugins.GetLimiter().Register()
}

// 注册全局服务
func registerServices() {
		logs.Info("init services")
		// 注册
		services.RegisterUrlService()
}

// 随机
func random(min,max int) int  {
		rand.Seed(time.Now().UnixNano())
		min, max = int(math.Max(float64(min), float64(max))), int(math.Min(float64(min), float64(max)))
		var n = rand.Intn(max) + min
		if n > max {
				return max
		}
		return n
}