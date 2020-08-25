package services

import (
		"github.com/astaxie/beego"
		"path"
		"strings"
		"sync"
)

type pathsService struct {
		BaseService
		paths   map[string]*map[string]string
		parsers map[string]func(...string) string
}

var (
		_pathService *pathsService
		_pathLock    sync.Once
)

func newPathService() {
		_pathService = new(pathsService)
		_pathService.Init()
}

func PathsServiceOf() *pathsService {
		if _pathService == nil {
				_pathLock.Do(newPathService)
		}
		return _pathService
}

func (this *pathsService) Init() {
		this.paths = make(map[string]*map[string]string)
		this.Constructor = func(args ...interface{}) interface{} {
				return PathsServiceOf()
		}
}

func (this *pathsService) Resolve(name string, key ...string) string {
		if v, ok := this.paths[name]; ok {
				if len(key) == 0 {
						return (*v)["root"]
				}
				k := strings.Join(key, ".")
				return (*v)[k]
		}
		if parser, ok := this.parsers[name]; ok {
				return parser(key...)
		}
		return ""
}

func (this *pathsService) BasePath(sub ...string) string {
		root := this.Resolve("base_path")
		if root == "" {
				root = beego.AppPath
				if root == "" {
						return path.Join(sub...)
				}
		}
		arr := []string{root}
		arr = append(arr, sub...)
		return path.Join(arr...)
}

func (this *pathsService) StoragePath(sub ...string) string {
		root := this.Resolve("storage_path")
		if root == "" {
				root = this.BasePath("/static/storage")
		}
		arr := []string{root}
		arr = append(arr, sub...)
		return path.Join(arr...)
}

func (this *pathsService) LogPath(sub ...string) string {
		root := this.Resolve("log_path")
		if root == "" {
				root = this.BasePath("/static/logs")
		}
		arr := []string{root}
		arr = append(arr, sub...)
		return path.Join(arr...)
}

func (this *pathsService) CachePath(sub ...string) string {
		root := this.Resolve("cache_path")
		if root == "" {
				root = this.BasePath("/static/cache")
		}
		arr := []string{root}
		arr = append(arr, sub...)
		return path.Join(arr...)
}

func (this *pathsService) Register(name string, v interface{}) *pathsService {
		if m, ok := v.(map[string]string); ok {
				list, ok := this.paths[name]
				if ok {
						for k, v := range m {
								(*list)[k] = v
						}
				} else {
						this.paths[name] = &m
				}
		}
		if fn, ok := v.(func(...string) string); ok {
				_, ok := this.parsers[name]
				if ok {
						return this
				}
				this.parsers[name] = fn
		}
		return this
}
