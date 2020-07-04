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
}

type avatarServerImpl struct {
		BaseService
}

var (
		_avatarLock   sync.Once
		_avatarServer *avatarServerImpl
)

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
