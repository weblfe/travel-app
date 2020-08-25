package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transports"
		"time"
)

type ThumbsUpRepository interface {
		Up() common.ResponseJson
		Down() common.ResponseJson
		Count() common.ResponseJson
}

type thumbsUpRepositoryImpl struct {
		service services.ThumbsUpService
		ctx     common.BaseRequestContext
}

const (
		ActionUp   = 1
		ActionDown = 0
)

func NewThumbsUpRepository(ctx common.BaseRequestContext) ThumbsUpRepository {
		var repository = new(thumbsUpRepositoryImpl)
		repository.init()
		repository.ctx = ctx
		return repository
}

func (this *thumbsUpRepositoryImpl) init() {
		this.service = services.ThumbsUpServiceOf()
}

// 用户点赞
func (this *thumbsUpRepositoryImpl) Up() common.ResponseJson {
		return this.Action(ActionUp)
}

// 点赞行为
func (this *thumbsUpRepositoryImpl) Action(action int) common.ResponseJson {
		var (
				request = transports.NewThumbUpRequest()
				userId  = this.ctx.Session(middlewares.AuthUserId)
		)
		if userId == nil || userId == "" {
				return common.NewUnLoginResp("请求先登录!")
		}
		err := request.Load(this.ctx.GetInput().RequestBody)
		if err != nil || request.IsEmpty() {
				err = request.ParseFrom(this.ctx.GetInput())
				if err != nil {
						return common.NewErrorResp(common.NewErrors(common.EmptyParamCode, err), "点赞失败")
				}
		}
		if this.ctx.Method() == "Delete" {
				action = ActionDown
		}
		var uid = userId.(string)
		if uid == "" {
				return common.NewUnLoginResp("请求先登录!")
		}
		switch action {
		case ActionUp:
				// 用户点赞
				var total = this.service.Up(request.GetType(), request.Id, uid)
				return common.NewSuccessResp(beego.M{"type": request.GetType(), "id": request.Id, "thumbsUpTotal": total, "timestamp": time.Now().Unix()}, "点赞成功")
		case ActionDown:
				// 用户取消点赞
				var total = this.service.Down(request.GetType(), request.Id, uid)
				return common.NewSuccessResp(beego.M{"type": request.GetType(), "id": request.Id, "thumbsUpTotal": total, "timestamp": time.Now().Unix()}, "请求成功")
		}
		return common.NewErrorResp(common.NewErrors(common.InvalidParametersCode, "请求失败,未知操作"))
}

// 用户取消点赞
func (this *thumbsUpRepositoryImpl) Down() common.ResponseJson {
		return this.Action(ActionDown)
}

func (this *thumbsUpRepositoryImpl) Count() common.ResponseJson {
		return common.NewInDevResp(this.ctx.GetActionId())
}
