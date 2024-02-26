package services

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"io"
		"io/ioutil"
		"os"
		"path/filepath"
		"strings"
		"sync"
)

type group struct {
		Type string
		App  string
}

type InitDataService interface {
		Load(file string, loader ...string) error
		SetLoader(name string, handler func(data []byte, filename string) bool, extras ...map[string]interface{}) InitDataService
}

type LoaderDataService interface {
		SetInit(path ...string) LoaderDataService
		Init()
}

type InitLoaderService interface {
		InitDataService
		LoaderDataService
}

type initDataServiceImpl struct {
		initPaths map[string]*struct {
				IsLoad bool
				Hash   string
		}
		loaders map[string]*struct {
				Handler func([]byte, string) bool
				Extras  map[string]interface{}
		}
}

var (
		_initDataServiceLock sync.Once
		_initDataService     *initDataServiceImpl
)

func GetInitDataServiceInstance() InitLoaderService {
		if _initDataService == nil {
				_initDataServiceLock.Do(newInitDataService)
		}
		return _initDataService
}

func newInitDataService() {
		_initDataService = new(initDataServiceImpl)
		_initDataService.init()
}

func (this *initDataServiceImpl) init() {
		this.initPaths = make(map[string]*struct {
				IsLoad bool
				Hash   string
		})
		this.loaders = make(map[string]*struct {
				Handler func([]byte, string) bool
				Extras  map[string]interface{}
		})
}

func (this *initDataServiceImpl) Init() {
		for fs, b := range this.initPaths {
				if b.IsLoad {
						continue
				}
				if err := this.Load(fs); err == nil {
						b.IsLoad = true
						this.initPaths[fs] = b
				}
		}
}

func (this *initDataServiceImpl) SetInit(path ...string) LoaderDataService {
		for _, fs := range path {
				state, err := os.Stat(fs)
				if err != nil {
						logs.Error(err)
						continue
				}
				if state.IsDir() {
						p, _ := filepath.Abs(fs)
						arr := this.readDir(p)
						for _, f := range arr {
								this.append(f)
						}
						continue
				}
				if f, err := filepath.Abs(fs); err == nil {
						this.append(f)
				}
		}
		return this
}

func (this *initDataServiceImpl) readDir(dir string) []string {
		var (
				err   error
				files []string
		)
		dir, err = filepath.Abs(dir)
		if err != nil {
				return []string{}
		}
		err = filepath.Walk(dir, func(pa string, info os.FileInfo, err error) error {
				if dir == pa {
						return nil
				}
				if err == nil {
						if info.IsDir() {
								files = append(files, this.readDir(pa)...)
						} else {
								files = append(files, pa)
						}
				}
				return err
		})
		if err == nil && len(files) > 0 {
				return files
		}
		return []string{}
}

func (this *initDataServiceImpl) append(fs string) {
		if _, ok := this.initPaths[fs]; ok {
				return
		}
		this.initPaths[fs] = &struct {
				IsLoad bool
				Hash   string
		}{
				IsLoad: false,
				Hash:   libs.FileHash(fs),
		}
}

func (this *initDataServiceImpl) Load(file string, loaderNames ...string) error {
		logs.Info("loading database file " + file)
		if len(loaderNames) == 0 {
				ext := filepath.Ext(file)
				name := filepath.Base(file)
				if name != "" && ext != "" {
						name = strings.Replace(name, ext, "", 1)
				}
				loader := this.getLoader(name)
				if ext != "" {
						ext = string([]rune(ext)[1:])
				}
				if loader == nil {
						loader = this.getLoader(ext)
				}
				if loader == nil {
						loader = this.getLoader(name + "." + ext)
				}
				if loader == nil {
						return fmt.Errorf(file + " loader not found! ")
				}
				data, err := ioutil.ReadFile(file)
				if err != nil {
						return err
				}
				if !loader.Handler(data, file) {
						return fmt.Errorf("loader " + file + " failed")
				}
		} else {
				data, err := ioutil.ReadFile(file)
				if err != nil {
						return err
				}
				for _, name := range loaderNames {
						obj := this.getLoader(name)
						if obj == nil {
								return fmt.Errorf("loader:%s not found", name)
						}
						if obj.Handler(data, file) {
								return fmt.Errorf("loader:%s  failed", name)
						}
				}
		}
		return nil
}

func (this *initDataServiceImpl) Storage(data []byte, writers ...io.Writer) error {
		if len(data) == 0 {
				return fmt.Errorf("empty stoarge call")
		}
		for _, w := range writers {
				_, err := w.Write(data)
				if err != nil {
						return err
				}
		}
		return nil
}

func (this *initDataServiceImpl) SetLoader(name string, handler func(data []byte, filename string) bool, extras ...map[string]interface{}) InitDataService {
		if _, ok := this.loaders[name]; ok {
				return this
		}
		if len(extras) == 0 {
				extras = append(extras, map[string]interface{}{})
		}
		this.loaders[name] = &struct {
				Handler func([]byte, string) bool
				Extras  map[string]interface{}
		}{
				Handler: handler,
				Extras:  extras[0],
		}
		return this
}

func (this *initDataServiceImpl) getLoader(name string) *struct {
		Handler func([]byte, string) bool
		Extras  map[string]interface{}
} {
		if name == "" {
				return nil
		}
		// 自然查询
		if loader, ok := this.loaders[name]; ok && loader != nil {
				return loader
		}
		low := strings.ToLower(name)
		// 全小写查询
		if low != name {
				if loader, ok := this.loaders[low]; ok && loader != nil {
						return loader
				}
		}
		tmp := []rune(name)
		if len(tmp) < 2 {
				return nil
		}
		// 首字母大写
		name = strings.ToUpper(string(tmp[0])) + string(tmp[1:])
		if loader, ok := this.loaders[name]; ok && loader != nil {
				return loader
		}
		return nil
}

func init() {
		var instance = GetInitDataServiceInstance()
		instance.SetLoader("sms", smsDataLoader)
		instance.SetLoader("tags", tagsDataLoader)
		instance.SetLoader("app_info", appInfoDataLoader)
		instance.SetLoader("roles", userRoleConfigDataLoader)
		instance.SetLoader("init", initLoader)
		instance.SetLoader("configs", configsLoader)
		instance.SetLoader("sensitives", sensitivesLoader)
}

// 短信数据数据模版加载器
func smsDataLoader(data []byte, filename string) bool {
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		if err := models.MessageTemplateModelOf().Adds(jsonArr); err != nil {
				return false
		}
		return true
}

// tag 自动 更新|添加|初始
func tagsDataLoader(data []byte, filename string) bool {
		if !strings.Contains(filename, "tags") {
				return false
		}
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		err := models.TagsModelOf().Adds(jsonArr)
		if err != nil {
				return true
		}
		return true
}

// app info 自动添加
func appInfoDataLoader(data []byte, filename string) bool {
		if !strings.Contains(filename, "app_info") {
				return false
		}
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		err := models.AppModelOf().Adds(jsonArr)
		if err != nil {
				return true
		}
		return true
}

// 用户角色 配置表
func userRoleConfigDataLoader(data []byte, filename string) bool {
		if !strings.Contains(filename, "roles") {
				return false
		}
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		err := models.UserRolesConfigModelOf().Adds(jsonArr)
		if err != nil {
				return true
		}
		return true
}

// 敏感词加载
func sensitivesLoader(data []byte, filename string) bool {
		if !strings.Contains(filename, "sensitives") {
				return false
		}
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		var items = getSensitivesByGroups(jsonArr)
		for g, words := range items {
				err := models.SensitiveWordsModelOf().Adds(words, g.Type, g.App)
				if err == nil {
						continue
				} else {
						logs.Error(err)
				}
		}
		return true
}

func getSensitivesByGroups(items []map[string]interface{}) map[group][]string {
		var (
				data    = make(map[group][]string)
				appName = models.NewAppInfo().GetAppName()
		)
		for _, it := range items {
				word, ok1 := it["word"]
				if !ok1 || word == nil || word == "" {
						continue
				}
				typ, ok2 := it["type"]
				if !ok2 || typ == nil || typ == "" {
						typ = models.SensitiveWordsTypeGlobal
				}
				app, ok3 := it["app"]
				if !ok3 || app == nil || app == "" {
						app = appName
				}
				groupKey := group{
						Type: typ.(string),
						App:  app.(string),
				}
				arr, ok := data[groupKey]
				if !ok {
						arr = []string{}
				}
				arr = append(arr, word.(string))
				data[groupKey] = arr
		}
		return data
}

// 加载json数据
func loaderJsons(data []byte) []map[string]interface{} {
		var arr []map[string]interface{}
		err := json.Unmarshal(data, &arr)
		if err != nil {
				return nil
		}
		return arr
}

// 加载系统初始数据
func initLoader(data []byte, filename string) bool {
		if !strings.Contains(filename, "init") {
				return false
		}
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		for _, item := range jsonArr {
				switch item["system"] {
				case "code":
						newAppService().InitCode()
				}
		}
		return true
}

// 配置信息载入初始化
func configsLoader(data []byte, filename string) bool {
		if !strings.Contains(filename, "configs") {
				return false
		}
		var jsonArr = loaderJsons(data)
		if jsonArr == nil {
				return false
		}
		if err := ConfigServiceOf().Adds(jsonArr); err == nil {
				return true
		}
		return false
}
