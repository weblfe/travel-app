package services

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"time"
)

type ApplyService interface {
		Commit(info *models.ApplyInfo) error
}

type applyServiceImpl struct {
		BaseService
		model *models.ApplyInfoModel
}

type ApplyHistory struct {
		Timestamp int64
		Content   string
}

func ApplyServiceOf() ApplyService {
		var service = new(applyServiceImpl)
		service.Init()
		return service
}

func (this *applyServiceImpl) Init() *applyServiceImpl {
		this.init()
		this.model = models.ApplyInfoModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return ApplyServiceOf()
		}
		return this
}

// Commit 提交举报 ｜ 反馈
func (this *applyServiceImpl) Commit(info *models.ApplyInfo) error {
		if info.Date == 0 {
				info.Date = models.GetDate()
		}
		if info.Type == "" {
				info.Type = models.ApplyTypeReport
		}
		if info.Content == "" {
				return common.NewErrors(common.InvalidParametersCode, common.InvalidParametersError)
		}
		var data = this.model.GetByUnique(info.M())
		// 新增
		if data == nil {
				info.InitDefault()
				return this.model.Add(info)
		}
		// 更新记录
		if data.Content != info.Content {
				data.Extras = models.Merger(data.Extras, info.Extras)
				var arr, ok = data.Extras["contents"]
				if !ok {
						arr = []ApplyHistory{}
				}
				data.Extras["contents"] = append(arr.([]ApplyHistory), ApplyHistory{
						Timestamp: data.UpdatedAt.Unix(), Content: data.Content,
				})
		}
		data.Id = info.Id
		data.CreatedAt = info.CreatedAt
		data.UpdatedAt = time.Now().Local()
		return this.model.UpdateById(info.Id.Hex(), data)
}
