package services

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"github.com/weblfe/travel-app/libs"
	"time"
)

type TicketService interface {
	Remove(string)bool
	Expired(string) bool
	GetStorageProvider() cache.Cache
	CreateTicket(expire time.Duration, extras ...map[string]interface{}) string
	GetTicketInfo(string) (map[string]interface{}, error)
}

type ticketServiceImpl struct {
	storage cache.Cache
	BaseService
}

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
	return libs.Md5(time.Now().String())
}

func (this *ticketServiceImpl) CreateTicket(expire time.Duration, extras ...map[string]interface{}) string {
	if len(extras) == 0 {
		extras = append(extras, map[string]interface{}{})
	}
	extras[0]["expired"] = time.Now().Add(expire).Unix()
	var (
		err     error
		ticket  = this.ticket()
		data, _ = json.Marshal(extras[0])
	)
	err = this.storage.Put(ticket, string(data), expire)
	if err == nil {
		return ticket
	}
	return ""
}

func (this *ticketServiceImpl) GetTicketInfo(s string) (map[string]interface{}, error) {
	 var (
	 	err error
	 	v = this.storage.Get(s)
	 	data = map[string]interface{}{}
	 )
	 if v == nil || v == "" {
	 	return nil,fmt.Errorf("not found")
	 }
	 err=json.Unmarshal(v.([]byte),&data)
	 if err !=nil {
	 	return nil, err
	 }
	 return data,nil
}

func (this *ticketServiceImpl)Remove(ticket string)bool  {
	if err:=this.storage.Delete(ticket);err!=nil {
		return false
	}
	return true
}

func (this *ticketServiceImpl)GetStorageProvider() cache.Cache  {
	return this.storage
}