package services

import (
	"github.com/weblfe/travel-app/libs"
	"time"
)

type UrlTicketService interface {
	TicketService
	GetUrl(string) string
	Incr(string) int
	GetAccessUrl(string) string
}

type urlTicketServiceImpl struct {
	ticketServiceImpl
	attachmentService AttachmentService
}

func UrlTicketServiceOf() UrlTicketService {
	var service = new(urlTicketServiceImpl)
	service.Init()
	return service
}

// 通过 ticket 获取 url
func (this *urlTicketServiceImpl) GetUrl(ticket string) string {
	var data, err = this.GetTicketInfo(ticket)
	if err != nil {
		return ""
	}
	return data["url"].(string)
}

// 统计
func (this *urlTicketServiceImpl) Incr(ticket string) int {
	var key = this.getUrlHash(ticket)
	if key == "" {
		return 0
	}
	if err := this.GetStorageProvider().Incr(key); err == nil {
		return 1
	}
	return 0
}

// 获取url hash 值
func (this *urlTicketServiceImpl) getUrlHash(ticket string) string {
	var url = this.GetUrl(ticket)
	if url == "" {
		return ""
	}
	return libs.Md5(url)
}

// 获取访问url
func (this *urlTicketServiceImpl) GetAccessUrl(ticket string) string {
	var url = this.GetUrl(ticket)
	if url != "" {
		go this.Incr(ticket)
	}
	return url
}

// 获取media ticket
func (this *urlTicketServiceImpl) GetMediaTicket(mediaId string) string {
	var url = this.attachmentService.GetUrl(mediaId)
	if url == "" {
		return ""
	}
	return this.CreateTicket(30*time.Minute, map[string]interface{}{"url": url})
}

// 获取附件 服务
func (this *urlTicketServiceImpl) getAttachmentService() AttachmentService {
	if this.attachmentService == nil {
		this.attachmentService = AttachmentServiceOf()
	}
	return this.attachmentService
}
