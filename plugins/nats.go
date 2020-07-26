package plugins

import (
		"encoding/json"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/nats-io/nats.go"
		"github.com/nats-io/nkeys"
		"io/ioutil"
		"os"
		"sync"
		"time"
)

// 配置
type NatsConfig struct {
		Name     string `json:"name"`
		ConnUrl  string `json:"connUrl"`
		Nkey     string `json:"nkey,omitempty"`
		User     string `json:"user,omitempty"`
		NkeySeed string `json:"nkeySeed,omitempty"`
		Password string `json:"password,omitempty"`
		Token    string `json:"token,omitempty"`
}

const (
		_NatsPluginName = "nats_plugin"
		NatsUserKey     = "NATS_USER"
		NatsName        = "NATS_NAME"
		NatsPassword    = "NATS_PASSWORD"
		NatsConnUrl     = "NATS_CONN_URL"
		NatsNkey        = "NATS_NKEY"
		NatsSeedKey     = "NATS_SEED_KEY"
		NatsTokenKey    = "NATS_TOKEN"
)

var (
		_NatsPlugInstance *NatsPlugin
		_lock             = sync.Mutex{}
		_syncLocks        = make(map[string]sync.Once)
)

// 插件
type NatsPlugin struct {
		Name   string      `json:"name"`
		Config *NatsConfig `json:"config"`
		conn   *nats.Conn  `json:",omitempty"`
}

func GetNatsPlugin() *NatsPlugin {
		if _NatsPlugInstance == nil {
				var lock = getLock(_NatsPluginName)
				lock.Do(func() {
						_NatsPlugInstance = newNatsPlugin()
						_NatsPlugInstance.Init()
				})
		}
		return _NatsPlugInstance
}

func getLock(name string) sync.Once {
		_lock.Lock()
		defer _lock.Unlock()
		if lock, ok := _syncLocks[name]; ok {
				return lock
		}
		var lock = sync.Once{}
		_syncLocks[name] = lock
		return lock
}

func newNatsPlugin() *NatsPlugin {
		var plugin = new(NatsPlugin)
		return plugin
}

func getString(key string) string {
		var v = beego.AppConfig.String(key)
		if v != "" {
				return v
		}
		return os.Getenv(key)
}

func (this *NatsPlugin) PluginName() string {
		return _NatsPluginName
}

func (this *NatsPlugin) Register() {
		Plugin(this.PluginName(), this)
}

func (this *NatsPlugin) Init() {
		if this.Config == nil {
				this.Config = NewNatsConfig()
				this.Config.init()
		}
		if this.Name == "" {
				this.Name = _NatsPluginName
		}
}

func (this *NatsPlugin) getConfig() *NatsConfig {
		if this.Config == nil {
				this.Config = NewNatsConfig()
		}
		return this.Config
}

func NewNatsConfig() *NatsConfig {
		var config = new(NatsConfig)
		return config.init()
}

func NewNatsConfigByFile(file string) *NatsConfig {
		var config = new(NatsConfig)
		return config.parse(file)
}

// 初始化配置
func (this *NatsConfig) init() *NatsConfig {
		this.Nkey = getString(NatsNkey)
		this.User = getString(NatsUserKey)
		this.ConnUrl = getString(NatsConnUrl)
		this.Token = getString(NatsTokenKey)
		this.NkeySeed = getString(NatsSeedKey)
		this.Password = getString(NatsPassword)
		this.Name = getString(NatsName)
		if this.Name == "" {
				this.Name = fmt.Sprintf("%s_%d", _NatsPluginName, time.Now().Unix())
		}
		return this
}

func (this *NatsConfig) getSignCb() func(bytes []byte) ([]byte, error) {
		return func(bytes []byte) ([]byte, error) {
				var akp, err = nkeys.FromSeed([]byte(this.NkeySeed))
				if err != nil {
						return nil, err
				}
				return akp.Sign(bytes)
		}
}

// 解析配置文件
func (this *NatsConfig) parse(file string) *NatsConfig {
		var state, err = os.Stat(file)
		if err == nil || state == nil {
				return this
		}
		if state.IsDir() {
				return this
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
				return this
		}
		if err := json.Unmarshal(data, this); err == nil {
				return this
		}
		return this
}

// 获取conn
func (this *NatsPlugin) GetConn() *nats.Conn {
		var nc, err = nats.Connect(this.getConnUrl(), this.getOptions()...)
		if err == nil {
				this.conn = nc
		}
		return this.conn
}

func (this *NatsPlugin) getOptions() []nats.Option {
		var options []nats.Option
		options = append(options, this.getOption())
		return options
}

// 获取option
func (this *NatsPlugin) getOption() func(*nats.Options) error {
		var config = this.getConfig()
		return func(options *nats.Options) error {
				if options.Token == "" && config.Token != "" {
						options.Token = config.Token
				}
				if options.Password == "" && config.Password != "" {
						options.Password = config.Password
				}
				if options.User == "" && config.User != "" {
						options.User = config.User
				}
				if options.Url == "" && config.ConnUrl != "" {
						options.Url = config.ConnUrl
				}
				if options.Name == "" && config.Name != "" {
						options.Name = config.Name
				}
				if options.SignatureCB == nil && config.NkeySeed != "" {
						options.SignatureCB = config.getSignCb()
				}
				if options.Nkey == "" && config.Nkey != "" {
						options.Nkey = config.Nkey
				}
				return nil
		}
}

func (this *NatsPlugin) getConnUrl() string {
		var url = this.getConfig().ConnUrl
		if url == "" {
				return nats.DefaultURL
		}
		return url
}
