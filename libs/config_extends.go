package libs

import (
		"errors"
		"github.com/astaxie/beego/logs"
		"github.com/spf13/viper"
		"net/url"
		"sync"
		"time"
)

var (
		_viper          *viper.Viper
		_cache          = sync.Map{}
		_configSyncLock = sync.Once{}
		_loaders        = initLoaders()
)

// UrlLoader url-loader
type UrlLoader struct {
		Name    string
		Match   func(*url.URL) bool
		Handler func(*url.URL) error
}

// 配置器
func newViper() {
		_viper = viper.New()
}

// 初始化加载器
func initLoaders() []*UrlLoader {
		var loaders = make([]*UrlLoader, 2)
		loaders = loaders[:0]
		return loaders
}

// Import 加载
func Import(source string, force ...bool) bool {
		if len(force) == 0 {
				force = append(force, false)
		}
		// 强制
		if !force[0] {
				if _, ok := _cache.Load(source); ok {
						return true
				}
		}
		var schema, err = url.Parse(source)
		if err != nil {
				return false
		}
		if err := checkUrl(schema); err != nil {
				logs.Error(err)
				return false
		}
		return loaderBySchema(schema)
}

// 检查url
func checkUrl(obj *url.URL) error {
		if obj.Scheme == "" {
				return errors.New("miss scheme")
		}
		if obj.Scheme != "file" {
				if obj.Host == "" {
						return errors.New("miss host")
				}
				if obj.User == nil {
						return errors.New("miss user")
				}
				var (
						user    = obj.User.Username()
						pass, _ = obj.User.Password()
				)
				if pass == "" {
						return errors.New("miss user pass")
				}
				if user == "" {
						return errors.New("miss user username")
				}
				return nil
		}
		if obj.Path == "" {
				return errors.New("miss path")
		}
		return nil
}

// GetConfigApp 获取 config App
func GetConfigApp() *viper.Viper {
		if _viper == nil {
				_configSyncLock.Do(newViper)
		}
		return _viper
}

// 通过 url 协议加载
func loaderBySchema(schema *url.URL) bool {
		var loader = getLoader(schema)
		if loader == nil {
				return false
		}
		if err := loader(schema); err != nil {
				logs.Error(err)
				return false
		}
		_cache.Store(schema.String(), time.Now().Unix())
		return true
}

// 获取加载处理器
func getLoader(obj *url.URL) func(obj *url.URL) error {
		if obj == nil || obj.Scheme == "" {
				return nil
		}
		var loaders []func(*url.URL) error
		for _, loader := range _loaders {
				if loader.Name != obj.Scheme {
						continue
				}
				if loader.Handler == nil {
						continue
				}
				if loader.Match != nil && !loader.Match(obj) {
						return loader.Handler
				}
				if loader.Match == nil {
						loaders = append(loaders, loader.Handler)
						continue
				}
		}
		if len(loaders) == 0 {
				return nil
		}
		return loaders[0]
}

// 注册加载处理器
func registerLoader(name string, match func(*url.URL) bool, handler func(*url.URL) error) {
		if handler == nil {
				return
		}
		for _, loader := range _loaders {
				if loader.Name != name {
						continue
				}

				if loader.Handler == nil {
						loader.Handler = handler
				}
				if loader.Match == nil {
						loader.Match = match
				}
				return
		}
		_loaders = append(_loaders, &UrlLoader{
				Name: name, Match: match, Handler: handler,
		})
}

// 注册
func init() {

		registerLoader("file", func(u *url.URL) bool {
				return true
		}, func(u *url.URL) error {

				return nil
		})

		registerLoader("etcd", func(u *url.URL) bool {
				return true
		}, func(u *url.URL) error {

				return nil
		})

		registerLoader("mysql", func(u *url.URL) bool {
				return true
		}, func(u *url.URL) error {

				return nil
		})

		registerLoader("mongodb", func(u *url.URL) bool {
				return true
		}, func(u *url.URL) error {

				return nil
		})
}
