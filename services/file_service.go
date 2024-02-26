package services

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/weblfe/travel-app/libs"
		"io"
		"io/ioutil"
		"mime/multipart"
		"os"
		"path"
		"path/filepath"
		"strings"
		"sync"
		"time"
)

type FileService interface {
		Get(file string, disk ...string) *os.File
		Exits(file string, disk ...string) bool
		Content(file string, disk ...string) []byte
		Save(file string, data []byte, disk ...string) (string, interface{})
		GetReader(file string, disk ...string) (io.Reader, error)
		GetWriter(file string, disk ...string) (io.Writer, error)
		SaveByReader(reader io.ReadCloser, extras beego.M) (beego.M, bool)
}

type FileSchemaWrapper interface {
		Register(name string, handler FileService)
}

type FileSystem interface {
		FileService
		FileSchemaWrapper
		Config() map[string]*FsConfig
		AddDisk(name string, root string, fn ...func(string) string) FileSystem
}

type FsConfig struct {
		Root string `json:"root"`
		Name string `json:"name"`
		Path func(string) string
}

var (
		_instanceLock sync.Once
		_fileSystem   *fileSystemServiceImpl
)

const (
		FileSystemGroupDiv             = ","
		FileSystemRootTpl              = "filesystem.%s.root"
		FileSystemJsonTpl              = "filesystem.%s.conf"
		AutoLoadFileSystemDiskGroupKey = "filesystem_groups"
)

type fileSystemServiceImpl struct {
		config  map[string]*FsConfig
		schemas map[string]FileService
}

func GetFileSystem() FileSystem {
		if _fileSystem == nil {
				_instanceLock.Do(newFileSystem)
		}
		return _fileSystem
}

func newFileSystem() {
		_fileSystem = new(fileSystemServiceImpl)
		_fileSystem.init()
}

func (this *fileSystemServiceImpl) init() {
		this.config = make(map[string]*FsConfig)
		this.schemas = make(map[string]FileService)
}

func (this *fileSystemServiceImpl) Exits(file string, disk ...string) bool {
		if len(disk) == 0 {
				disk = append(disk, "default")
		}
		fn, server := this.resolver(file)
		if fn != nil {
				fs := fn(file)
				if fs == "" {
						return false
				}
				if _, err := os.Stat(fs); err != nil {
						if os.IsExist(err) || os.IsNotExist(err) {
								return false
						}
				}
				return true
		}
		if server != nil {
				return server.Exits(file, disk...)
		}
		return false
}

func (this *fileSystemServiceImpl) Get(file string, disk ...string) *os.File {
		if len(disk) == 0 {
				disk = append(disk, "default")
		}
		fn, server := this.resolver(file)
		if fn != nil {
				if fs, err := os.Open(fn(file)); err == nil {
						return fs
				}
				return nil
		}
		if server != nil {
				return server.Get(file, disk...)
		}
		return nil
}

func (this *fileSystemServiceImpl) Content(file string, disk ...string) []byte {
		if len(disk) == 0 {
				disk = append(disk, "default")
		}
		var (
				tmp    = make([]byte, 1024)
				buffer []byte
		)
		reader, err := this.GetReader(file, disk...)
		if err == nil && reader != nil {
				for {
						tmp = tmp[0:0]
						n, err := reader.Read(tmp)
						if n > 0 && err != io.EOF {
								buffer = append(buffer, tmp...)
								continue
						}
						break
				}
		}
		return buffer
}

func (this *fileSystemServiceImpl) Save(file string, data []byte, disk ...string) (string, interface{}) {
		if len(disk) == 0 {
				disk = append(disk, "default")
		}
		fn, service := this.resolver(disk[0])
		if fn != nil && service == nil {
				file = fn(file)
				if err := ioutil.WriteFile(file, data, os.ModePerm|os.ModeAppend); err == nil {
						return file, true
				}
		}
		if fn == nil && service != nil {
				return service.Save(file, data, disk...)
		}
		return "", false
}

func (this *fileSystemServiceImpl) GetReader(file string, disk ...string) (io.Reader, error) {
		if len(disk) == 0 {
				disk = append(disk, "default")
		}
		fn, server := this.resolver(file)
		if fn != nil {
				return os.Open(fn(file))
		}
		if server != nil {
				return server.GetReader(file, disk...)
		}
		return nil, fmt.Errorf("get reader faild,disk %s nont exits", disk[0])
}

func (this *fileSystemServiceImpl) GetWriter(file string, disk ...string) (io.Writer, error) {
		if len(disk) == 0 {
				disk = append(disk, "default")
		}
		fn, server := this.resolver(file)
		if fn != nil {
				return os.Open(fn(file))
		}
		if server != nil {
				return server.GetWriter(file, disk...)
		}
		return nil, fmt.Errorf("get writer failed,disk %s nont exits", disk[0])
}

func (this *fileSystemServiceImpl) Register(name string, handler FileService) {
		if handler == nil {
				return
		}
		if _, ok := this.schemas[name]; ok {
				return
		}
		this.schemas[name] = handler
}

func (this *fileSystemServiceImpl) Config() map[string]*FsConfig {
		return this.config
}

func (this *fileSystemServiceImpl) resolver(disk string) (func(string) string, FileService) {
		if ser, ok := this.schemas[disk]; ok && ser != nil {
				return nil, ser
		}
		if m, ok := this.config[disk]; ok && m != nil {
				return this.path(m), nil
		}
		return nil, nil
}

func (this *fileSystemServiceImpl) path(conf *FsConfig) func(string) string {
		// 路径解析器
		if conf.Path != nil {
				return conf.Path
		}
		// 路径拼接
		if conf.Root != "" {
				conf.Path = func(s string) string {
						return path.Join(conf.Root, s)
				}
				return conf.Path
		}
		// 绝对路径
		return func(s string) string {
				if p, err := filepath.Abs(s); err == nil {
						return p
				}
				return ""
		}
}

// 自动载入相关disk 配置
func (this *fileSystemServiceImpl) load() {
		var groups = beego.AppConfig.String(AutoLoadFileSystemDiskGroupKey)
		if groups == "" {
				return
		}
		var arr = strings.SplitN(groups, FileSystemGroupDiv, -1)
		for _, disk := range arr {
				disk = strings.TrimSpace(disk)
				// @todo ${AppPath} 动态解析
				root := beego.AppConfig.String(fmt.Sprintf(FileSystemRootTpl, disk))
				if root != "" {
						this.AddDisk(disk, root, nil)
						continue
				}
				conf := beego.AppConfig.String(fmt.Sprintf(FileSystemJsonTpl, disk))
				if conf != "" {
						cnf := &FsConfig{}
						if err := json.Unmarshal([]byte(conf), cnf); err == nil {
								this.config[disk] = cnf
						}
				}
		}
}

func (this *fileSystemServiceImpl) AddDisk(disk string, root string, fn ...func(string) string) FileSystem {
		if _, ok := this.config[disk]; ok {
				return this
		}
		if len(fn) == 0 {
				fn = append(fn, nil)
		}
		this.config[disk] = &FsConfig{
				Name: disk,
				Root: root,
				Path: fn[0],
		}
		return this
}

func (this *fileSystemServiceImpl) SaveByReader(reader io.ReadCloser, extras beego.M) (beego.M, bool) {
		var (
				savePath string
				filePath = extras["path"]
		)
		if filePath == "" || filePath == nil {
				filePath = this.getFileSavePath()
		}
		filename := extras["filename"]
		if filename == "" || filename == nil {
				filename = this.getFileNameByMapper(extras)
		}
		root := filePath.(string)
		name := filename.(string)
		if strings.Contains(name, root) {
				savePath = name
		} else {
				savePath = filepath.Join(root, name)
		}
		// 文件是否存在
		savePath = libs.UniqueFile(savePath)
		fs, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		defer func() {
				_ = fs.Close()
		}()
		if err != nil {
				logs.Error(err)
				return extras, false
		}
		n, err := io.Copy(fs, reader)
		if err == nil && n > 0 {
				info, _ := fs.Stat()
				extras["size"] = info.Size()
				extras["filePath"] = savePath
				extras["hash"] = libs.FileHash(savePath)
				extras["filename"] = info.Name()
				return extras, true
		}
		return nil, false
}

func (this *fileSystemServiceImpl) getFileSavePath() string {
		var (
				t       = time.Now()
				y, m, d = t.Date()
				date    = fmt.Sprintf("%d-%d-%d", y, m, d)
				dir     = fmt.Sprintf("%s/%s", PathsServiceOf().StoragePath(), date)
		)
		_ = os.MkdirAll(dir, os.ModePerm)
		return dir
}

func (this *fileSystemServiceImpl) getFileNameByMapper(m beego.M) string {
		if fileInfo, ok := m["fileInfo"]; ok {
				if info, ok := fileInfo.(os.FileInfo); ok {
						return info.Name()
				}
				if info, ok := fileInfo.(*multipart.FileHeader); ok {
						return info.Filename
				}
		}
		return fmt.Sprintf("tmp_%d", time.Now().Unix())
}
