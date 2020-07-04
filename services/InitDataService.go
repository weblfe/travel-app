package services

import (
		"encoding/json"
		"fmt"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"io"
		"io/ioutil"
		"os"
		"path/filepath"
		"strings"
		"sync"
)

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
		GetInitDataServiceInstance().SetLoader("sms", smsDataLoader)
}

// 短信数据数据模版加载器
func smsDataLoader(data []byte, filename string) bool {
		var arr []map[string]interface{}
		err := json.Unmarshal(data, &arr)
		if err != nil {
				return false
		}
		if err := models.MessageTemplateModelOf().Adds(arr); err != nil {
				return false
		}
		return true
}