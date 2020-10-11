package kernel

import (
		"fmt"
		"github.com/coreos/etcd/clientv3"
		docker "github.com/fsouza/go-dockerclient"
		"log"
		"strconv"
		"strings"
		"time"
)

// docker 服务自动发信
type DockerServiceRegister interface {
		Get(name string) string
		Host(name string) string
		DelService(serName string)
		Register(host string, args ...Args)
		CreateService(host string, serName string)
		Query(args Args) (map[string]string, error)
		Loader(loaders...func(serviceRegister DockerServiceRegister,info *DockerInfo))
}

// docker 注册信息
type DockerInfo struct {
		data  EntryArr
		Hosts *HostInfo
}

// 状态值
type State int

const (
		Off            State = -1 // 不可用
		Unknown        State = 0  // 未知
		On             State = 1  // 可用
		FlagDockerHost       = "docker_host"
		FlagService          = "service"
)

// host 信息
type HostInfo struct {
		HostName string
		Ipv4     string
		Ipv6     string
		Port     int
		Status   State
}

// 配置键
type Entry struct {
		Key   string
		Value string
}

type EntryArr []*Entry

// args
type Args struct {
		K string
		V interface{}
}

type dockerServiceRegisterImpl struct {
		client   *docker.Client
		info     *DockerInfo
		registry *clientv3.Client
}

func NewHostInfo() *HostInfo {
		var info = new(HostInfo)
		return info
}

func NewDockerInfo() *DockerInfo {
		var info = new(DockerInfo)
		info.init()
		return info
}

func (this *dockerServiceRegisterImpl) Get(name string) string {
		panic("implement me")
}

func (this *dockerServiceRegisterImpl) Host(name string) string {
		panic("implement me")
}

func (this *dockerServiceRegisterImpl) DelService(serName string) {
		panic("implement me")
}

func (this *dockerServiceRegisterImpl) Register(host string, args ...Args) {
		panic("implement me")
}

func (this *dockerServiceRegisterImpl) CreateService(host string, serName string) {
		panic("implement me")
}

func (this *dockerServiceRegisterImpl) Query(args Args) (map[string]string, error) {
		panic("implement me")
}


func (this EntryArr) GetInt(key string, defaults int) int {
		for _, it := range this {
				if it == nil || it.Key != key {
						continue
				}
				n, err := strconv.Atoi(strings.TrimSpace(it.Value))
				if err == nil {
						return n
				}
		}
		return defaults
}

func (this EntryArr) Get(key string, defaults string) string {
		for _, it := range this {
				if it == nil || it.Key != key {
						continue
				}
				return it.Value
		}
		return defaults
}

func (this *DockerInfo) init() {
		this.Hosts = NewHostInfo()
		this.data = make(EntryArr, 10)
		this.data = this.data[:0]
}

func (this *DockerInfo) Init(handlers ...func(data *EntryArr, info *HostInfo)) {
		if len(handlers) == 0 {
				return
		}
		for _, handler := range handlers {
				handler(&this.data, this.Hosts)
		}
}

func (this *HostInfo) getEndPoint() string {
		return this.GetIp() + this.GetPort()
}

func (this *HostInfo) GetIp() string {
		if this.Ipv6 == "" {
				if this.Ipv4 == "" {
						return "127.0.0.1"
				}
				return this.Ipv4
		}
		return this.Ipv6
}

func (this *HostInfo) GetPort() string {
		if this.Port == 0 {
				return "2375"
		}
		return fmt.Sprintf("%v", this.Port)
}

func newDockerServiceRegister() *dockerServiceRegisterImpl {
		var ins = new(dockerServiceRegisterImpl)
		return ins.init()
}

func (this *dockerServiceRegisterImpl) init() *dockerServiceRegisterImpl {
		this.info = NewDockerInfo()
		return this
}

func (this *dockerServiceRegisterImpl) Boot() {
		this.GetDocker()
		this.GetRegistry()
		this.Check()
}

func (this *dockerServiceRegisterImpl) GetRegistry() *clientv3.Client {
		var err error
		if this.registry == nil {
				this.registry, err = clientv3.New(this.getRegistryCnf())
		}
		if err != nil {
				log.Fatal(err)
		}
		return this.registry
}

func (this *dockerServiceRegisterImpl) Loader(loaders...func(serviceRegister DockerServiceRegister,info *DockerInfo)) {
		for _, loader := range loaders {
				loader(this,this.info)
		}
}

func (this *dockerServiceRegisterImpl) GetDocker() *docker.Client {
		var err error
		if this.client == nil {
				this.client, err = docker.NewClient(this.info.Hosts.getEndPoint())
		}
		if err != nil {
				log.Fatal(err)
		}
		return this.client
}

func (this *dockerServiceRegisterImpl) getRegistryCnf() clientv3.Config {
		return clientv3.Config{
				Endpoints:   strings.Split(this.info.data.Get("registry_endpoints", "127.0.0.1:2379"), ","),
				DialTimeout: time.Duration(int64(this.info.data.GetInt("dial_timeout", 5)) * int64(time.Second)),
		}
}

// 检查 docker 宿主机
func (this *dockerServiceRegisterImpl) Check() {

}
