package services

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"github.com/weblfe/travel-app/libs"
	"os"
	"time"
)

type TicketService interface {
	Remove(string) bool
	Expired(string) bool
	GetStorageProvider() cache.Cache
	CreateTicket(data map[string]interface{}) string
	GetTicketInfo(string) (map[string]interface{}, error)
}

type ticketServiceImpl struct {
	storage cache.Cache
	expire  time.Duration
	BaseService
}

const (
	defaultExpire = 17520 * time.Hour // 2 year
)

func TicketServiceOf(storage ...CacheService) TicketService {
	var service = new(ticketServiceImpl)
	if len(storage) > 0 {
		service.storage = storage[0].Get("default")
	}
	service.Init()
	return service
}

func (this *ticketServiceImpl) Init() {
	this.init()
	if this.storage == nil {
		this.storage = GetCacheService().Get("default")
	}
}

func (this *ticketServiceImpl) Expired(s string) bool {
	if !this.storage.IsExist(s) {
		return false
	}
	return true
}

func (this *ticketServiceImpl) ticket() string {
	return libs.Md5(fmt.Sprintf("%v", time.Now().UnixNano()))
}

func (this *ticketServiceImpl) CreateTicket(data map[string]interface{}) string {
	var expire time.Duration
	if v, ok := data["expire"]; ok {
		expire, ok = v.(time.Duration)
		if !ok {
			expire = this.getExpire()
		}
	}
	data["expired"] = time.Now().Add(expire).Unix()
	var (
		err      error
		ticket   = this.ticket()
		_data, _ = json.Marshal(data)
	)
	err = this.storage.Put(ticket, string(_data), expire)
	if err == nil {
		return ticket
	}
	return ""
}

func (this *ticketServiceImpl) GetTicketInfo(s string) (map[string]interface{}, error) {
	var (
		err  error
		v    = this.storage.Get(s)
		data = map[string]interface{}{}
	)
	if v == nil || v == "" {
		return nil, fmt.Errorf("not found")
	}
	err = json.Unmarshal(v.([]byte), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *ticketServiceImpl) Remove(ticket string) bool {
	if err := this.storage.Delete(ticket); err != nil {
		return false
	}
	return true
}

func (this *ticketServiceImpl) GetStorageProvider() cache.Cache {
	return this.storage
}

func (this *ticketServiceImpl) SetExpire(e time.Duration) bool {
	if this.expire > 0 {
		return false
	}
	this.expire = e
	return true
}

func (this *ticketServiceImpl) getExpire() time.Duration {
	if this.expire > 0 {
		return this.expire
	}
	var du = os.Getenv("TICKET_EXPIRE_DURATION")
	if du == "" {
		return defaultExpire
	}
	if d, err := time.ParseDuration(du); err == nil {
		this.expire = d
		return d
	}
	return defaultExpire
}
