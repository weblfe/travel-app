package services

import (
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"strings"
		"time"
)

type UrlTicketService interface {
		TicketService
		GetUrl(string) string
		Incr(string) int
		GetAccessUrl(string) string
		GetTicketUrlByAttach(*models.Attachment) string
		GetTicketInfoToSimple(ticket string) *SimpleUrlAttach
}

const (
		_IncrPrefix = "incr_"
)

type urlTicketServiceImpl struct {
		ticketServiceImpl
		attachmentService AttachmentService
}

func UrlTicketServiceOf() UrlTicketService {
		return newUrlTicketService()
}

func newUrlTicketService() *urlTicketServiceImpl {
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
		if err := this.GetStorageProvider().Incr(this.IncrKey(key)); err == nil {
				return 1
		}
		return 0
}

// 统计访问次数key
func (this *urlTicketServiceImpl)IncrKey(key string) string  {
		if strings.Contains(key,_IncrPrefix) {
				return key
		}
		return _IncrPrefix + key
}

// 获取url hash 值
func (this *urlTicketServiceImpl) getUrlHash(ticket string) string {
		var data = this.GetTicketInfoToSimple(ticket)
		if data == nil {
				return ""
		}
		// 使用 mediaId 作为唯一访问统计 incr
		if data.MediaId != "" {
				return data.MediaId
		}
		// 未记录的资源使用 url hash
		return libs.Md5(data.Url)
}

// 获取访问url
func (this *urlTicketServiceImpl) GetAccessUrl(ticket string) string {
		var url = this.GetUrl(ticket)
		if url != "" {
				go this.Incr(ticket)
		}
		return url
}

// 获取访问的mediaId
func (this *urlTicketServiceImpl) GetAccessMediaId(ticket string) string {
		var url = this.GetUrl(ticket)
		if url != "" {
				go this.Incr(ticket)
		}
		return url
}

// 是否过期
func (this *urlTicketServiceImpl) Expired(s string) bool {
		if !this.ticketServiceImpl.Expired(s) {
				return false
		}
		this.Incr(s)
		return true
}

// 获取media
func (this *urlTicketServiceImpl) GetMediaId(ticket string) string {
		var data, err = this.GetTicketInfo(ticket)
		if err != nil {
				return ""
		}
		return data["mediaId"].(string)
}

// 获取media ticket
func (this *urlTicketServiceImpl) GetMediaTicket(mediaId string) string {
		var attach = this.attachmentService.Get(mediaId)
		return this.GetTicketUrlByAttach(attach)
}

// 获取附件 服务
func (this *urlTicketServiceImpl) getAttachmentService() AttachmentService {
		if this.attachmentService == nil {
				this.attachmentService = AttachmentServiceOf()
		}
		return this.attachmentService
}

// 获取 访问ticket
func (this *urlTicketServiceImpl) GetTicketUrlByAttach(attach *models.Attachment) string {
		if attach == nil || attach.DeletedAt != 0 {
				return ""
		}
		var url = attach.CdnUrl
		if url == "" && attach.Url != "" {
				url = attach.Url
		}
		var ticket = this.CreateTicket(30*time.Minute, map[string]interface{}{"url": url, "mediaId": attach.Id.Hex(), "path": attach.Path})
		return strings.Replace(url, attach.Id.Hex(), ticket, -1)
}

// 获取简单数据对象通过 ticket
func (this *urlTicketServiceImpl) GetTicketInfoToSimple(ticket string) *SimpleUrlAttach {
		var data, err = this.GetTicketInfo(ticket)
		if err != nil {
				return nil
		}
		return NewSimpleUrlAttach(data)
}

type SimpleUrlAttach struct {
		Url     string `json:"url"`
		MediaId string `json:"mediaId"`
		Path    string `json:"path"`
}

func NewSimpleUrlAttach(data ...map[string]interface{}) *SimpleUrlAttach {
		var attach = new(SimpleUrlAttach)
		attach.Load(data[0])
		return attach
}

func (this *SimpleUrlAttach) Load(data map[string]interface{}) *SimpleUrlAttach {
		if data == nil {
				return this
		}
		for key, v := range data {
				this.Set(key, v)
		}
		return this
}

func (this *SimpleUrlAttach) Set(key string, v interface{}) *SimpleUrlAttach {
		switch key {
		case "url":
				this.Url = v.(string)
		case "mediaId":
				this.MediaId = v.(string)
		case "path":
				this.Path = v.(string)
		}
		return this
}
