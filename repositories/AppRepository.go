package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"time"
)

type AppRepository interface {
		GetConfig(typ string) common.ResponseJson
		Apply() common.ResponseJson
}

type appRepository struct {
		ctx             common.BaseRequestContext
		service         services.AppService
		appApplyService services.ApplyService
}

func NewAppRepository(ctx common.BaseRequestContext) AppRepository {
		var repository = new(appRepository)
		repository.ctx = ctx
		repository.Init()
		return repository
}

func (this *appRepository) Init() {
		this.service = services.AppServiceOf()
		this.appApplyService = services.ApplyServiceOf()
}

// 获取配置
func (this *appRepository) GetConfig(driver string) common.ResponseJson {
		var items = this.service.GetAppInfos(driver)
		if len(items) == 0 {
				return common.NewErrorResp(common.NewErrors(common.NotFound, "config empty"), "配置为空")
		}
		return common.NewSuccessResp(items, "获取配置成功")
}

// 申请，举报
func (this *appRepository) Apply() common.ResponseJson {
		var (
				err    error
				userId = getUserId(this.ctx)
				data   = models.NewApplyInfo()
		)
		data.UserId = userId
		data.Date = models.GetDate()
		data.Content = this.ctx.GetString("content")
		data.Type = this.ctx.GetString("type", models.ApplyTypeReport)
		data.Extras = this.ctx.GetJsonData()
		err = this.appApplyService.Commit(data)
		if err == nil {
				return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "提交成功")
		}
		return common.NewFailedResp(common.ServiceFailed, "今日反馈到底上限")
}
