package services

import (
		"fmt"
		"github.com/astaxie/beego/config/env"
		"github.com/weblfe/travel-app/models"
		"strings"
		"sync"
)

type AvatarService interface {
		GetDefaultAvatar(...int) *models.Attachment
		GetAvatarUrlById(string) string
		GetAvatarUrlDefault(...int) string
		GetAvatarById(string, ...int) *AvatarInfo
}

type avatarServerImpl struct {
		BaseService
}

var (
		_avatarLock   sync.Once
		_avatarServer *avatarServerImpl
)

type AvatarInfo struct {
		Url    string `json:"url"`
		Id     string `json:"id"`
		Gender int    `json:"gender"`
}

func AvatarServerOf() AvatarService {
		if _avatarServer == nil {
				_avatarLock.Do(newAvatarServer)
		}
		return _avatarServer
}

func newAvatarServer() {
		_avatarServer = new(avatarServerImpl)
		_avatarServer.init()
		_avatarServer.Init()
}

func (this *avatarServerImpl) Init() {
		this.Constructor = func(args ...interface{}) interface{} {
				return AvatarServerOf()
		}
}

func (this *avatarServerImpl) GetDefaultAvatar(gender ...int) *models.Attachment {
		if len(gender) == 0 {
				gender = append(gender, 0)
		}
		gen := models.GetGenderKey(gender[0])
		id := env.Get(fmt.Sprintf("USER_%s_AVATAR_MEDIA_ID", strings.ToUpper(gen)), "")
		if id == "" {
				return nil
		}
		m := AttachmentServiceOf().Get(id)
		if m == nil {
				return nil
		}
		return m
}

func (this *avatarServerImpl) GetAvatarById(id string, gender ...int) *AvatarInfo {
		var (
				info = new(AvatarInfo)
				url  = this.GetAvatarUrlById(id)
		)
		if len(gender) == 0 {
				gender = append(gender, 0)
		}
		if url == "" {
				url = this.GetAvatarUrlDefault(gender...)
		}
		info.Url = url
		info.Gender = gender[0]
		info.Id = id
		return info
}

func (this *avatarServerImpl) GetAvatarUrlById(id string) string {
		var attach = AttachmentServiceOf().Get(id)
		return UrlTicketServiceOf().GetTicketUrlByAttach(attach)
}

func (this *avatarServerImpl) GetAvatarUrlDefault(gender ...int) string {
		var attach = this.GetDefaultAvatar(gender...)
		return UrlTicketServiceOf().GetTicketUrlByAttach(attach)
}
