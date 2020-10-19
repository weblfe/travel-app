package kernel

import (
		"context"
		"errors"
		"github.com/coreos/etcd/clientv3"
		"github.com/spf13/viper"
		"log"
		"os"
		"path/filepath"
		"strings"
		"time"
)

// 配置逻辑参数
type ConfigureArgs struct {
		File     string   // 配置文件
		Prefix   string   // 配置前缀
		Excludes []string // 排除
		EtcdCnf  Etcd     // etcd 配置
}

// etcd 配置
type Etcd struct {
		clientv3.Config
}

var loaders []func(cmder ConfigureCmder, args *ConfigureArgs)

type ConfigureCmder interface {
		Init() ConfigureCmder
		Boot()
		Exec() error
		Loader(loaders ...func(cmder ConfigureCmder, args *ConfigureArgs))
}

type configureCmderImpl struct {
		property *ConfigureArgs
		client   *clientv3.Client
		data     map[string]string
}

func newConfigureCmderImpl() *configureCmderImpl {
		var ins = new(configureCmderImpl)
		ins.Init()
		return ins
}

func AddLoader(loader func(cmder ConfigureCmder, args *ConfigureArgs)) {
		loaders = append(loaders, loader)
}

func GetLoaders() []func(cmder ConfigureCmder, args *ConfigureArgs) {
		return loaders
}

func GetConfigureIns() *configureCmderImpl {
		return newConfigureCmderImpl()
}

func InvokerConfigure(file string, prefix string, excludes []string, endpoints string, timeout int64) *configureCmderImpl {
		var ins = GetConfigureIns()
		AddLoader(func(cmder ConfigureCmder, args *ConfigureArgs) {
				args.File = file
				args.Prefix = prefix
				args.Excludes = excludes
				args.EtcdCnf.Endpoints = strings.Split(endpoints, ",")
				args.EtcdCnf.DialTimeout = time.Second * time.Duration(timeout)
		})
		ins.Boot()
		return ins
}

func (this *configureCmderImpl) Init() ConfigureCmder {
		if this.property == nil {
				this.property = newConfigureArgs()
		}
		if this.data == nil {
				this.data = make(map[string]string)
		}
		return this
}

func (this *configureCmderImpl) Boot() {
		for _, loader := range GetLoaders() {
				loader(this, this.GetProperty())
		}
		this.property.Boot()
}

func (this *configureCmderImpl) Exec() error {
		var ins = this.getProvider()
		// 推送配置
		data := this.getData()
		if len(data) == 0 {
				log.Println("empty configure data push!")
				return errors.New("empty configure data push")
		}
		if err := push(*ins, this.getData()); err != nil {
				log.Fatal(err)
				return err
		} else {
				log.Println("push config file <" + this.GetProperty().File + "> success")
				return errors.New("push config file <" + this.GetProperty().File + "> success")
		}
}

func (this *configureCmderImpl) getData() map[string]string {
		if len(this.data) == 0 {
				this.data = this.getConfigFileData()
		}
		return this.data
}

func (this *configureCmderImpl) getProvider() *clientv3.Client {
		var err error
		if this.client == nil {
				this.client, err = clientv3.New(this.GetProperty().EtcdCnf.Row())
		}
		if err != nil {
				log.Fatal(err)
		}
		return this.client
}

func (this *configureCmderImpl) GetProperty() *ConfigureArgs {
		return this.property
}

func newConfigureArgs() *ConfigureArgs {
		var args = new(ConfigureArgs)
		return args.init()
}

func (this *ConfigureArgs) Boot() {

}

func (this *configureCmderImpl) getConfigFileData() map[string]string {
		var (
				file   = this.GetProperty().File
				prefix = this.GetProperty().Prefix
		)
		stat, err := os.Stat(file)
		if err != nil {
				return map[string]string{}
		}
		if stat.IsDir() {
				this.readDir(file, prefix)
		}
		return this.readFile(file, prefix)
}

func (this *configureCmderImpl) Loader(loaders ...func(cmder ConfigureCmder, args *ConfigureArgs)) {
		for _, loader := range loaders {
				loader(this, this.GetProperty())
		}
}

func (this *configureCmderImpl) readFile(filename string, prefix string) map[string]string {
		var (
				data   = map[string]string{}
				loader = this.getLoader()
		)
		fs, err := os.Open(filename)
		if err != nil {
				log.Fatal(err)
		}
		if err := loader.ReadConfig(fs); err != nil {
				log.Fatal(err)
		}
		for _, key := range loader.AllKeys() {
				data[key] = loader.GetString(key)
		}
		return data
}

func (this *configureCmderImpl) getLoader() *viper.Viper {
		return viper.New()
}

func (this *configureCmderImpl) readDir(file string, prefix string) map[string]string {
		var data = map[string]string{}
		filepath.Walk(file, func(path string, info os.FileInfo, err error) error {

				return err
		})
		return data
}

// 初始化
func (this *ConfigureArgs) init() *ConfigureArgs {
		this.etcdInit()
		return this
}

func (this *ConfigureArgs) etcdInit() {
		if len(this.EtcdCnf.Endpoints) == 0 {
				this.EtcdCnf.Endpoints = strings.Split(this.resolve("etcd-url"), ",")
		}
}

func (this *ConfigureArgs) resolve(key string) string {

		return ""
}

func (this Etcd) Row() clientv3.Config {
		return clientv3.Config{
				Endpoints:            this.Endpoints,
				AutoSyncInterval:     this.AutoSyncInterval,
				DialTimeout:          this.DialTimeout,
				DialKeepAliveTime:    this.DialKeepAliveTime,
				DialKeepAliveTimeout: this.DialKeepAliveTimeout,
				MaxCallSendMsgSize:   this.MaxCallSendMsgSize,
				MaxCallRecvMsgSize:   this.MaxCallRecvMsgSize,
				TLS:                  this.TLS,
				Username:             this.Username,
				Password:             this.Password,
				RejectOldCluster:     this.RejectOldCluster,
				DialOptions:          this.DialOptions,
				Context:              this.Context,
		}
}

func push(client clientv3.Client, data map[string]string) error {
		var ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
		var txn = client.Txn(ctx)
		resp, err := txn.If().Then(puts(data)...).Commit()
		if err != nil {
				log.Println(err)
				return err
		}
		if resp.Succeeded {
				return nil
		}
		return errors.New("push failed")
}

func puts(data map[string]string) []clientv3.Op {
		var ops []clientv3.Op
		for key, v := range data {
				ops = append(ops, clientv3.OpPut(key, v))
		}
		return ops
}
