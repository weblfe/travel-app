package repositories

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

// 用户相关行为
type BehaviorRepository interface {
		FocusOff() common.ResponseJson
		FocusOn() common.ResponseJson
}

type userBehaviorRepositoryImpl struct {
		ctx         common.BaseRequestContext
		userService services.UserService
		dto         *DtoRepository
}

func NewBehaviorRepository(ctx common.BaseRequestContext) BehaviorRepository {
		var repository = newUserBehaviorRepository()
		repository.ctx = ctx
		return repository
}

func newUserBehaviorRepository() *userBehaviorRepositoryImpl {
		var repository = new(userBehaviorRepositoryImpl)
		repository.init()
		return repository
}

func (this *userBehaviorRepositoryImpl) init() {
		this.dto = GetDtoRepository()
		this.userService = services.UserServiceOf()
}

// 取消关注
func (this *userBehaviorRepositoryImpl) FocusOff() common.ResponseJson {
		// var userId, ok = this.ctx.GetParam(":userId")

		panic("implement me")
}

// 关注
func (this *userBehaviorRepositoryImpl) FocusOn() common.ResponseJson {
		// var userId, ok = this.ctx.GetParam(":userId")

		panic("implement me")
}
