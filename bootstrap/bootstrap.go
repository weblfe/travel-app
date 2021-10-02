package bootstrap

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
	"path/filepath"
	"strings"
	"time"
)

type appBootstrap struct {
	appPath string
}

var _starter = &appBootstrap{}

// Startup 启动
func (starter *appBootstrap) Startup() {
	starter.init()
}

func (starter *appBootstrap) SetAppPath(appPath string) *appBootstrap {
	if starter.appPath != "" {
		return starter
	}
	starter.appPath = appPath
	if !filepath.IsAbs(starter.appPath) {
		starter.appPath, _ = filepath.Abs(starter.appPath)
	}
	if !strings.HasSuffix(starter.appPath, "/") {
		starter.appPath = starter.appPath + "/"
	}
	return starter
}

// 载入引导逻辑
func (starter *appBootstrap) init() {
	logs.Info("bootstrap start")
	starter.bootstrap()
	logs.Info("bootstrap end...")
}

// 引导逻辑
func (starter *appBootstrap) bootstrap() {
	// 环境注册
	starter.initRegisterEnv()
	// 注册结构体
	starter.registerGob()
	// swagger
	starter.initSwagger()
	// session
	starter.initSession()
	// database
	starter.initDatabase()
	// middleware
	starter.initMiddleware()
	// 注册插件
	starter.registerPlugins()
	// 注册全局服务
	starter.registerServices()
}

// 配置 session
func (starter *appBootstrap) initSession() {
	if libs.InArray(beego.AppConfig.String("session_on"), []string{"on", "1", "true", "yes"}) {
		logs.Info("init session")
		beego.BConfig.WebConfig.Session.SessionOn = true
		beego.BConfig.WebConfig.Session.SessionProvider = beego.AppConfig.DefaultString("session_driver", "redis")
		beego.BConfig.WebConfig.Session.SessionProviderConfig = beego.AppConfig.DefaultString("session_config", "127.0.0.1:6379")
	}
}

// 初始 swagger
func (starter *appBootstrap) initSwagger() {
	if beego.BConfig.RunMode == "dev" {
		logs.Info("init swagger")
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/static/swagger"] = "swagger"
	}
}

// 初始化数据库
func (starter *appBootstrap) initDatabase() {
	mode := beego.BConfig.RunMode
	database, err := beego.AppConfig.GetSection(mode + ".database")
	if err == nil {
		if driver, ok := database["db_driver"]; ok && driver == "mongodb" {
			starter.initMongodb(database)
		}
	}
	starter.initMigration()
	logs.Info("init database")
}

// 初始化数据迁移
func (starter *appBootstrap) initMigration() {
	service := services.GetInitDataServiceInstance()
	service.SetInit(starter.resolvePath("./static/database"))
	service.Init()
	logs.Info("init database migration")
}

// 初始 mongodb
func (starter *appBootstrap) initMongodb(data map[string]string) {
	for key, v := range data {
		models.SetProfile(key, v)
	}
	logs.Info("init database profiles")
	starter.initMongoIndex()
	logs.Info("init database index create")
}

// 初始化 数据索引
func (starter *appBootstrap) initMongoIndex() {
	models.UserModelOf().CreateIndex(true)
	models.PostsModelOf().CreateIndex(true)
	models.ConfigModelOf().CreateIndex(true)
	models.ThumbsUpModelOf().CreateIndex(true)
	models.SensitiveWordsModelOf().CreateIndex(true)
	models.UserRolesConfigModelOf().CreateIndex(true)
	logs.Info("init database index")
}

// 注册环境变量
func (starter *appBootstrap) initRegisterEnv() {
	logs.Info("init env")
	var pwd string
	if starter.appPath == "" {
		pwd, _ = os.Getwd()
	} else {
		pwd = starter.appPath
	}
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
	logs.Info("load env files: %v", arr)
	_ = godotenv.Load(arr...)
	// 重新载入env
	starter.reloadEnv()
	// 加载 主配置
	_ = beego.LoadAppConfig("ini", "conf/main.conf")
	// 重新载入配置
	starter.reloadConfig()
}

// 重新载入配置
func (starter *appBootstrap) reloadConfig() {
	beego.BConfig.RunMode = beego.AppConfig.DefaultString(starter.setConfGlobalScope("runmode"), beego.BConfig.RunMode)
	beego.BConfig.AppName = beego.AppConfig.DefaultString(starter.setConfGlobalScope("appname"), beego.BConfig.AppName)
	beego.BConfig.ServerName = beego.AppConfig.DefaultString(starter.setConfGlobalScope("servername"), beego.BConfig.ServerName)
	beego.BConfig.Listen.HTTPPort = beego.AppConfig.DefaultInt(starter.setConfGlobalScope("httpport"), beego.BConfig.Listen.HTTPPort)
	beego.BConfig.Listen.HTTPSPort = beego.AppConfig.DefaultInt(starter.setConfGlobalScope("httpsport"), beego.BConfig.Listen.HTTPSPort)
}

// 重新载入env
func (starter *appBootstrap) reloadEnv() {
	for _, e := range os.Environ() {
		splits := strings.Split(e, "=")
		env.Set(splits[0], os.Getenv(splits[0]))
	}

	var (
		on       = env.Get("CONFIGURE_CENTER_ON", "")
		provider = env.Get("CONFIGURE_CENTER_PROVIDER", "etcd")
		appId    = env.Get("CONFIGURE_CENTER_APPID", "travel-app")
	)
	logs.Info("CONFIGURE_CENTER_ON:" + on)
	// 是否加载 配置中心
	if on != "" && on != "0" && on != "false" {
		var manger = plugins.GetConfigureCentreRepositoryMangerInstance()
		manger.InitDef().Boot()
		providerIns := manger.Get(provider)
		if providerIns == nil {
			logs.Error("ConfigureCentre Provider :" + provider + " miss")
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
			ticket = ticket + 10*time.Second
		}
		return
	}
}

// SetConfGlobalScope 全局
func (starter *appBootstrap) setConfGlobalScope(key string) string {
	if strings.Contains(key, "::") {
		return key
	}
	return "default::" + key
}

// 初始middleware
func (starter *appBootstrap) initMiddleware() {
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
func (starter *appBootstrap) registerGob() {
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
func (starter *appBootstrap) registerPlugins() {
	logs.Info("init plugins")
	plugins.GetOSS().Register()
	plugins.GetQrcode().Register()
	plugins.GetNatsPlugin().Register()
	plugins.GetLimiter().Register()
}

// 注册全局服务
func (starter *appBootstrap) registerServices() {
	logs.Info("init services")
	// 注册
	services.RegisterUrlService()
}

// 解析 path
func (starter *appBootstrap) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if strings.Contains(path, "./") {
		return strings.Replace(path, "./", starter.appPath, 1)
	}
	return filepath.Join(starter.appPath, path)
}

// 随机
func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	min, max = int(math.Max(float64(min), float64(max))), int(math.Min(float64(min), float64(max)))
	var n = rand.Intn(max) + min
	if n > max {
		return max
	}
	return n
}

func StartUp(appPath ...string) {
	_path, _ := os.Getwd()
	appPath = append(appPath, _path)
	_starter.SetAppPath(appPath[0]).Startup()
}

func Run() {
	logs.Info("api server start....")
	beego.Run()
}
