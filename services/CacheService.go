package services

import (
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/cache"
		_ "github.com/astaxie/beego/cache/memcache"
		_ "github.com/astaxie/beego/cache/redis"
		_ "github.com/astaxie/beego/cache/ssdb"
		"strings"
		"sync"
)

type CacheService interface {
		Get(string) cache.Cache
		Add(string, *CacheConfig) CacheService
}

type cacheServiceImpl struct {
		config    map[string]*CacheConfig
		instances map[string]cache.Cache
		locker    *sync.Mutex
}

const (
		GroupDot            = ","
		AutoLoadCacheInsKey = "cache_groups"
		CacheDriverTpl      = "%s_cache_driver"
		CacheConfigTpl      = "%s_cache_config"
		CacheDefaultDriver  = "redis"
)

// 配置
type CacheConfig struct {
		Driver string `json:"driver"`
		Config string `json:"config"`
}

var (
		once                 sync.Once
		cacheServiceInstance *cacheServiceImpl
)

// 获取
func GetCacheService() CacheService {
		if cacheServiceInstance == nil {
				once.Do(newCacheService)
		}
		return cacheServiceInstance
}

func newCacheService() {
		cacheServiceInstance = new(cacheServiceImpl)
		cacheServiceInstance.init()
}

func (this *cacheServiceImpl) init() {
		this.locker = &sync.Mutex{}
		this.instances = make(map[string]cache.Cache)
		this.config = make(map[string]*CacheConfig)
		this.load()
}

func (this *cacheServiceImpl) load() {
		var insArr = beego.AppConfig.String(AutoLoadCacheInsKey)
		if insArr == "" {
				return
		}
		var groups = strings.SplitN(insArr, GroupDot, -1)
		for _, name := range groups {
				k := strings.TrimSpace(name)
				if k == "" {
						continue
				}
				if _, ok := this.config[k]; ok {
						continue
				}
				config := beego.AppConfig.String(fmt.Sprintf(CacheConfigTpl, k))
				driver := beego.AppConfig.DefaultString(fmt.Sprintf(CacheDriverTpl, k), CacheDefaultDriver)
				if config == "" && driver == "" {
						continue
				}
				this.config[name] = &CacheConfig{
						Driver: driver,
						Config: config,
				}
		}
}

func (this *cacheServiceImpl) Add(name string, config *CacheConfig) CacheService {
		if config == nil {
				return this
		}
		if _, ok := this.config[name]; ok {
				return this
		}
		this.config[name] = config
		return this
}

func (this *cacheServiceImpl) Get(name string) cache.Cache {
		if ins, ok := this.instances[name]; ok && ins != nil {
				return ins
		}
		return this.invoker(name)
}

func (this *cacheServiceImpl) invoker(name string) cache.Cache {
		this.locker.Lock()
		defer this.locker.Unlock()
		if m, ok := this.config[name]; ok {

				ins, err := cache.NewCache(m.Driver, m.Config)
				if err == nil && ins != nil {
						this.instances[name] = ins
				}
				return ins
		}
		return nil
}
