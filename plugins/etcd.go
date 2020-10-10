package plugins

import (
		"context"
		"fmt"
		"github.com/astaxie/beego/logs"
		"github.com/coreos/etcd/clientv3"
		"github.com/coreos/etcd/mvcc/mvccpb"
		"github.com/pkg/errors"
		"log"
		"os"
		"strings"
		"sync"
		"time"
)

const (
		EtcdProvider        = "etcd"
		EtcdEndpointsEnvKey = "ETCD_ENDPOINTS"
)

func GetEtcdClient() (*clientv3.Client, error) {
		var cli, err = clientv3.New(getEtcdCnf())
		if err != nil {
				return nil, err
		}
		return cli, nil
}

func getEtcdCnf() clientv3.Config {
		var Endpoints = os.Getenv(EtcdEndpointsEnvKey)
		if Endpoints == "" {
				Endpoints = "localhost:2379"
		}
		return clientv3.Config{
				Endpoints:   strings.Split(Endpoints, ","),
				DialTimeout: 5 * time.Second,
		}
}

type ConfigureCentreRepository interface {
		Del(keys []string, appId ...string) int
		Put(key string, value interface{}, appId ...string) error
		Get(key string, appId ...string) (map[string]string, error)
		Pull(appId string) (map[string]string, error)
}

type EtcdConfigureCentreRepository interface {
		ConfigureCentreRepository
		Watch(args WatchArgs)
		Count(key string) (int64, error)
		Keys(key ...string) ([]string, error)
}

type WatchArgs struct {
		Key     string
		Ctx     context.Context
		Options []clientv3.OpOption
		Handler func(watchChan clientv3.WatchChan)
}

type EtcdWatcherHandler func(eventType mvccpb.Event_EventType, kv mvccpb.KeyValue, event *clientv3.Event)

type ConfigBootStrap interface {
		Boot()
}

type configureCentreRepositoryMangerImpl struct {
		Providers map[string]ConfigureCentreRepository
}

type etcdConfigureCentreRepositoryImpl struct {
		ProviderType string
		client       *clientv3.Client
		properties   map[string]interface{}
}

var (
		configureCentreRepositoryMangerIns *configureCentreRepositoryMangerImpl
		_InsLocker                         sync.Once
)

func GetConfigureCentreRepositoryMangerInstance() *configureCentreRepositoryMangerImpl {
		if configureCentreRepositoryMangerIns == nil {
				_InsLocker.Do(func() {
						configureCentreRepositoryMangerIns = new(configureCentreRepositoryMangerImpl)
						configureCentreRepositoryMangerIns.init()
				})
		}
		return configureCentreRepositoryMangerIns
}

func (this *configureCentreRepositoryMangerImpl) init() {
		if this.Providers == nil {
				this.Providers = make(map[string]ConfigureCentreRepository)
		}
}

func (this *configureCentreRepositoryMangerImpl) Get(key string) ConfigureCentreRepository {
		if pro, ok := this.Providers[key]; ok {
				return pro
		}
		return nil
}

func (this *configureCentreRepositoryMangerImpl) Register(key string, provider interface{}) *configureCentreRepositoryMangerImpl {
		if pro, ok := provider.(ConfigureCentreRepository); ok {
				_, ok1 := this.Providers[key]
				if ok1 {
						return this
				}
				this.Providers[key] = pro
				return this
		}
		if factory, ok := provider.(func() ConfigureCentreRepository); ok {
				_, ok1 := this.Providers[key]
				if ok1 {
						return this
				}
				this.Providers[key] = factory()
				return this
		}
		return this
}

func (this *configureCentreRepositoryMangerImpl) InitDef() *configureCentreRepositoryMangerImpl {
		this.Register(EtcdProvider,EtcdConfigureCentreRepositoryOf())
		return this
}

func (this *configureCentreRepositoryMangerImpl) Boot() {
		for _, provider := range this.Providers {
				if provider == nil {
						continue
				}
				if bootStrap, ok := provider.(ConfigBootStrap); ok {
						bootStrap.Boot()
				}
		}
}

func newEtcdConfigureCentreRepository() ConfigureCentreRepository {
		var rep = new(etcdConfigureCentreRepositoryImpl)
		rep.init()
		return rep
}

func EtcdConfigureCentreRepositoryOf() ConfigureCentreRepository {
		return newEtcdConfigureCentreRepository()
}

func (this *etcdConfigureCentreRepositoryImpl) init() {
		this.ProviderType = EtcdProvider
		this.properties = make(map[string]interface{})
}

func (this *etcdConfigureCentreRepositoryImpl) Boot() {
		var err error
		if this.client == nil {
				this.client, err = GetEtcdClient()
		}
		if err != nil {
				this.client = nil
				logs.Error(err)
				log.Fatal(err)
				return
		}
}

func (this *etcdConfigureCentreRepositoryImpl) getClient() *clientv3.Client {
		if this.client == nil {
				this.Boot()
		}
		return this.client
}

func (this *etcdConfigureCentreRepositoryImpl) getCmdCtx() (context.Context, context.CancelFunc) {
		var (
				ctx context.Context
				fn  context.CancelFunc
		)
		ctx, fn = context.WithTimeout(context.Background(), this.get("context_timeout", time.Second).(time.Duration))
		return ctx, fn
}

func (this *etcdConfigureCentreRepositoryImpl) Get(key string, appId ...string) (map[string]string, error) {
		var ctx, fn = this.getCmdCtx()
		if ctx == nil {
				return nil, errors.New("getCmdCtx failed!")
		}
		defer fn()
		var (
				resp, err = this.getClient().Get(ctx, etcdKey(key, arrStrFirst(appId)))
		)
		if resp == nil {
				return nil, errors.New("Get Value failed!")
		}
		if err != nil {
				logs.Error(err)
				return nil, err
		}
		var data = make(map[string]string)
		for _, val := range resp.Kvs {
				data[string(val.Key)] = string(val.Value)
		}
		return data, err
}

func (this *etcdConfigureCentreRepositoryImpl) Del(keys []string, appId ...string) int {
		var (
				count   = 0
				ctx, fn = this.getCmdCtx()
				client  = this.getClient()
		)
		defer fn()
		for _, key := range keys {
				resp, err := client.Delete(ctx, etcdKey(key, arrStrFirst(appId)))
				if err != nil {
						logs.Error(err)
						continue
				}
				if resp.Deleted > 0 {
						count++
				}
		}
		return count
}

func (this *etcdConfigureCentreRepositoryImpl) Put(key string, value interface{}, appId ...string) error {
		var (
				ctx, fn = this.getCmdCtx()
				client  = this.getClient()
		)
		defer fn()
		_, err := client.Put(ctx, etcdKey(key, arrStrFirst(appId)), fmt.Sprintf("%v", value))
		return err
}

func (this *etcdConfigureCentreRepositoryImpl) Pull(appId string) (map[string]string, error) {
		var (
				ctx, fn = this.getCmdCtx()
				client  = this.getClient()
		)
		defer fn()
		resp, err := client.Get(ctx, appId, clientv3.WithPrefix())
		if err != nil {
				return nil, err
		}
		var data = make(map[string]string)
		for _, kvs := range resp.Kvs {
				data[string(kvs.Key)] = string(kvs.Value)
		}
		return data, nil
}

func (this *etcdConfigureCentreRepositoryImpl) Watch(args WatchArgs) {
		var (
				client = this.getClient()
				wch    clientv3.WatchChan
		)
		if args.Key == "" {
				log.Fatal(errors.New("miss watch key"))
				return
		}
		if len(args.Options) == 0 {
				wch = client.Watch(args.Ctx, args.Key)
		} else {
				wch = client.Watch(args.Ctx, args.Key, args.Options...)
		}
		if args.Handler != nil {
				args.Handler(wch)
		}
}

func (this *etcdConfigureCentreRepositoryImpl) Keys(key ...string) ([]string, error) {
		var (
				keys    []string
				client  = this.getClient()
				ctx, fn = this.getCmdCtx()
				resp    *clientv3.GetResponse
				err     error
		)
		defer fn()
		if len(key) == 0 {
				resp, err = client.Get(ctx, "/", clientv3.WithKeysOnly(), clientv3.WithPrefix())
		} else {
				resp, err = client.Get(ctx, arrStrFirst(key), clientv3.WithKeysOnly(), clientv3.WithPrefix())
		}
		if err != nil {
				return nil, err
		}

		for _, kvs := range resp.Kvs {
				keys = append(keys, string(kvs.Key))
		}
		return keys, nil
}

func (this *etcdConfigureCentreRepositoryImpl) Count(key string) (int64, error) {
		var (
				client  = this.getClient()
				ctx, fn = this.getCmdCtx()
				resp    *clientv3.GetResponse
				err     error
		)
		defer fn()
		if len(key) == 0 {
				resp, err = client.Get(ctx, "/", clientv3.WithCountOnly(), clientv3.WithPrefix())
		} else {
				resp, err = client.Get(ctx, key, clientv3.WithCountOnly(), clientv3.WithPrefix())
		}
		if err != nil {
				return 0, err
		}
		return resp.Count, nil
}

func (this *etcdConfigureCentreRepositoryImpl) get(key string, defaults ...interface{}) interface{} {
		if v, ok := this.properties[key]; ok {
				return v
		}
		if len(defaults) == 0 {
				return nil
		}
		return defaults[0]
}

func (this *etcdConfigureCentreRepositoryImpl) getCmdOption(args ...string) clientv3.OpOption {
		return func(op *clientv3.Op) {
				if len(args) <= 0 {
						return
				}
				for _, v := range args {
						argValue(v)
				}
		}
}

func (this *etcdConfigureCentreRepositoryImpl) SetProperty(key string, value interface{}) *etcdConfigureCentreRepositoryImpl {
		this.properties[key] = value
		return this
}

func (this *etcdConfigureCentreRepositoryImpl) Close() {
		if this.client == nil {
				return
		}
		if err := this.client.Close(); err != nil {
				logs.Error(err)
		}
		this.client = nil
}

func arrStrFirst(arr []string) string {
		if len(arr) == 0 {
				return ""
		}
		return arr[0]
}

func argOption(key, value string) string {
		return key + "->" + value
}

func argValue(data string) (key, value string) {
		if data == "" || len(data) == 0 {
				return "", ""
		}
		if !strings.Contains(data, "->") {
				return "", ""
		}
		return arrStrExplode(strings.SplitN(data, "->", 2))
}

func arrStrExplode(arr []string) (val1, val2 string) {
		var n = len(arr)
		if n >= 2 {
				return arr[0], arr[1]
		}
		if n == 0 {
				return "", ""
		}
		return arr[0], ""
}

func etcdKey(key string, root ...string) string {
		if len(root) == 0 {
				return etcRootKey(key)
		}
		var prefix []string
		for _, k := range root {
				if k == "" || strings.TrimSpace(k) == "" {
						continue
				}
				prefix = append(prefix, k)
		}
		if len(prefix) == 0 {
				return key
		}
		if strings.Index(strings.Join(prefix, "/"), key) >= 1 {
				return key
		}
		return etcRootKey(strings.Join(append(root, key), "/"))
}

func etcRootKey(key string) string {
		if key[0] == '/' {
				return key
		}
		return "/" + key
}

func NewWatchHandler(handler EtcdWatcherHandler) func(watchChan clientv3.WatchChan) {
		return func(watchChan clientv3.WatchChan) {
				for arg := range watchChan {
						for _, event := range arg.Events {
								handler(event.Type, *event.Kv, event)
						}
				}
		}
}
