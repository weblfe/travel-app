package services

import "github.com/astaxie/beego/cache"

type TicketService interface {
	Expired(string) bool
	CreateTicket(expire int64, extras ...map[string]interface{}) string
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
	return  true
}

func (this *ticketServiceImpl) CreateTicket(expire int64, extras ...map[string]interface{}) string {
	panic("implement me")
}

func (this *ticketServiceImpl) GetTicketInfo(s string) (map[string]interface{}, error) {
	panic("implement me")
}


