package components

import (
		"context"
		"errors"
		"github.com/coreos/etcd/clientv3"
		"github.com/spf13/viper"
		"log"
		"os"
		"strings"
		"time"
)

const (
		EntryPoints    = "entryPoints"
		Username       = "username"
		Password       = "password"
		TimeOut        = "timeout"
		HostName       = "host"
		HostSchema     = "127.0.0.1:2380"
		DotDiv         = ","
		CnfFile        = "file"
		DefaultTimeOut = 10 * time.Second
		WatchEventPrefix = "watcher."
)

type Option interface {
		Key() string
		V() string
		Value() interface{}
		Copy() Option
}

type Registry interface {
		Storage
		Boot() error
		SetFile(file string) Registry
		SetEntryPort(host string, options ...Option)
		Save(file string, handlers ...func(fs *os.File, storage Storage) error) error
}

type Storage interface {
		Del(keys []string) error
		SetOptions(options ...Option)
		Get(key string) (string, error)
		Set(key string, value string) error
		GetWatcher() clientv3.Watcher
		Pull(keys []string) (map[string]string, error)
}

type RegistryImpl struct {
		ctx       map[string]interface{}
		entryPort string
		storage   Storage
		container *viper.Viper
		handlers  map[string]func(*viper.Viper,Storage)error
}

func (this *RegistryImpl) GetWatcher() clientv3.Watcher {
		return this.storage.GetWatcher()
}

func (this *RegistryImpl) SetOptions(options ...Option) {
		if len(options) == 0 {
				return
		}
		for _, opt := range options {
				this.ctx[opt.Key()] = opt.V()
		}
}

type EtcdStorageImpl struct {
		ApiClient *clientv3.Client
		Options   []Option
		Ctx       context.Context
}

type OptionImpl struct {
		K, Val string
}

func (this *OptionImpl) Key() string {
		return this.K
}

func (this *OptionImpl) V() string {
		return this.Val
}

func (this *OptionImpl) Value() interface{} {
		return this.Val
}

func (this *OptionImpl) Copy() Option {
		return NewOption(this.K, this.Val)
}

func (this *OptionImpl) String() string {
		return "{\"" + this.Key() + "\":\"" + this.V() + "\"}"
}

func NewOption(key, value string) Option {
		var opt = new(OptionImpl)
		opt.K = key
		opt.Val = value
		return opt
}

func RegistryOf(options ...Option) Registry {
		var registry = NewRegistry()
		return registry.options(options)
}

func NewEtcdStorage(option ...Option) Storage {
		var storageIns = new(EtcdStorageImpl)
		storageIns.init()
		if len(option) == 0 {
				return storageIns
		}
		for _, opt := range option {
				storageIns.Options = append(storageIns.Options, opt)
		}
		return storageIns
}

func (this *EtcdStorageImpl) init() {
		this.Ctx = nil
		this.ApiClient = nil
		this.Options = make([]Option, 2)
}

func (this *EtcdStorageImpl) GetWatcher()clientv3.Watcher {
		return this.ApiClient
}

func (this *EtcdStorageImpl) SetOptions(options ...Option) {
		for _, opt := range options {
				this.Options = append(this.Options, opt)
		}
}

func (this *EtcdStorageImpl) getClient() *clientv3.Client {
		var err error
		if this.ApiClient == nil {
				this.ApiClient, err = clientv3.New(this.getCfg())
				if err == nil {
						log.Fatal(err)
						return nil
				}
		}
		return this.ApiClient
}

func (this *EtcdStorageImpl) getContext() context.Context {
		return this.Ctx
}

func (this *EtcdStorageImpl) SetContext(ctx context.Context) Storage {
		this.Ctx = ctx
		return this
}

func (this *EtcdStorageImpl) getCfg() clientv3.Config {
		return clientv3.Config{
				Endpoints: this.getEndpoints(),
				Username:  this.getUserName(),
				Password:  this.getPassword(),
				Context:   this.getContext(),
		}
}

func (this *EtcdStorageImpl) getEndpoints() []string {
		return strings.Split(this.get(EntryPoints), DotDiv)
}

func (this *EtcdStorageImpl) getUserName() string {
		return this.get(Username)
}

func (this *EtcdStorageImpl) getPassword() string {
		return this.get(Password)
}

func (this *EtcdStorageImpl) get(key string) string {
		for _, opt := range this.Options {
				if opt.Key() == key {
						return opt.V()
				}
		}
		return ""
}

func (this *EtcdStorageImpl) setOptions(key, value string) Storage {
		this.Options = append(this.Options, NewOption(key, value))
		return this
}

func (this *EtcdStorageImpl) Del(keys []string) error {
		for _, key := range keys {
				ctx, fn := this.getRequestCtx()
				resp, err := this.getClient().Delete(ctx, key)
				if err != nil {
						return err
				}
				fn()
				log.Printf("%v\n", resp.OpResponse())
		}
		return nil
}

func (this *EtcdStorageImpl) getRequestCtx() (context.Context, context.CancelFunc) {
		return context.WithTimeout(context.Background(), this.getTimeOut())
}

func (this *EtcdStorageImpl) getTimeOut() time.Duration {
		var (
				t      = this.get(TimeOut)
				d, err = time.ParseDuration(t)
		)
		if err == nil {
				return d
		}
		return DefaultTimeOut
}

func (this *EtcdStorageImpl) Get(key string) (string, error) {
		panic("implement me")
}

func (this *EtcdStorageImpl) Set(key string, value string) error {
		panic("implement me")
}

func (this *EtcdStorageImpl) Pull(keys []string) (map[string]string, error) {
		panic("implement me")
}

func NewRegistry() *RegistryImpl {
		var registry = new(RegistryImpl)
		registry.init()
		return registry
}

func (this *RegistryImpl) options(options []Option) Registry {
		var (
				host = ""
				opts = make([]Option, 2)
		)
		// 设置
		for _, opt := range options {
				if opt.Key() == EntryPoints {
						host = opt.V()
						continue
				}
				opts = append(opts, opt)
		}
		// 所选项
		if host == "" && 0 != len(opts) {
				this.SetEntryPort(host, opts...)
		}
		return this
}

func (this *RegistryImpl) init() {
		this.ctx = make(map[string]interface{}, 3)
		this.entryPort = HostSchema
		this.storage = nil
}

func (this *RegistryImpl) getStorageProvider() Storage {
		var storageImpl = NewEtcdStorage(this.getStorageOptions()...)
		return storageImpl
}

func (this *RegistryImpl) getStorageOptions() []Option {
		return []Option{
				NewOption(EntryPoints, this.getOptionValue(EntryPoints)),
				NewOption(Username, this.getOptionValue(Username)),
				NewOption(Password, this.getOptionValue(Password)),
		}
}

func (this *RegistryImpl) getOptionValue(key string) string {
		if v, ok := this.ctx[key]; ok {
				return v.(string)
		}
		return ""
}

func (this *RegistryImpl) Boot() error {
		if this.storage == nil {
				this.storage = this.getStorageProvider()
		}
		this.handlers = make(map[string]func(*viper.Viper,Storage)error)
		return nil
}

func (this *RegistryImpl) SetFile(file string) Registry {
		info, err := os.Stat(file)
		if err != nil {
				return this
		}
		if info.IsDir() {
				return this
		}
		this.ctx[CnfFile] = file
		return this
}

func (this *RegistryImpl) SetEntryPort(host string, options ...Option) {
		if len(options) == 0 {
				options = append(options, NewOption(HostName, host))
		}
		this.storage.SetOptions(options...)
}

func (this *RegistryImpl)Watch()  {
		for  {
				this.storage.GetWatcher()
				select {

				}
		}
}

func (this *RegistryImpl)On(name string,handler func(config *viper.Viper,storage Storage)error)  {
		this.handlers[name] = handler
}

func (this *RegistryImpl) Del(keys []string) error {
		return this.storage.Del(keys)
}

func (this *RegistryImpl) Get(key string) (string, error) {
		return this.storage.Get(key)
}

func (this *RegistryImpl) Set(key string, value string) error {
		return this.storage.Set(key, value)
}

func (this *RegistryImpl) Pull(keys []string) (map[string]string, error) {
		return this.storage.Pull(keys)
}

func (this *RegistryImpl) Save(file string, handlers ...func(fs *os.File, storage Storage) error) error {
		var (
				fs         *os.File
				state, err = os.Stat(file)
		)
		if err != nil {
				if err != os.ErrExist {
						return err
				}
				fs, err = os.Create(file)
				if err != nil {
						return err
				}
		}
		if state.IsDir() {
				return errors.New("is not normal file")
		}
		if fs == nil {
				fs, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModePerm)
		}
		if len(handlers) == 0 {
				handlers = append(handlers,this.updateConfigFile)
		}
		for _,handler:=range handlers{
				err = handler(fs,this.storage)
				if err!=nil {
						return err
				}
		}
		return nil
}

func (this *RegistryImpl)updateConfigFile(fs *os.File,storage Storage) error  {

		return nil
}