package repositories

import "github.com/weblfe/travel-app/common"

type AuthenticateRepository interface {
		QQ() common.ResponseJson
		WeChat() common.ResponseJson
		Apple() common.ResponseJson
}

type authenticateRepositoryImpl struct {
		ctx common.BaseRequestContext
}

func (this *authenticateRepositoryImpl) QQ() common.ResponseJson {
		panic("implement me")
}

func (this *authenticateRepositoryImpl) WeChat() common.ResponseJson {
		panic("implement me")
}

func (this *authenticateRepositoryImpl) Apple() common.ResponseJson {
		panic("implement me")
}

func NewAuthenticateRepository(ctx common.BaseRequestContext) AuthenticateRepository {
		var repository = new(authenticateRepositoryImpl)
		repository.ctx = ctx
		return repository
}
