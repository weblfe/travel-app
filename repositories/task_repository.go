package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/services"
)

type TaskRepository interface {
		SyncAssetsTask() common.ResponseJson
}

type taskRepositoryImpl struct {
		ctx common.BaseRequestContext
}


func (this *taskRepositoryImpl) SyncAssetsTask() common.ResponseJson {
		var count = services.AttachmentServiceOf().SyncOssTask()
		return common.NewSuccessResp(beego.M{"count": count}, "任务执行中...")
}

func NewTaskRepository(ctx common.BaseRequestContext) TaskRepository {
		var repository = new(taskRepositoryImpl)
		repository.ctx = ctx
		return repository
}
