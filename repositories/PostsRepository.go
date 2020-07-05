package repositories

import (
		"encoding/json"
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/middlewares"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"time"
)

type PostsRepository interface {
		Create() common.ResponseJson
		Update() common.ResponseJson
		Lists(...string) common.ResponseJson
		GetById(...string) common.ResponseJson
		RemoveId(...string) common.ResponseJson
}

type postRepositoryImpl struct {
		ctx     *beego.Controller
		service services.PostService
}

func NewPostsRepository(ctx *beego.Controller) PostsRepository {
		var repository = new(postRepositoryImpl)
		repository.init()
		repository.ctx = ctx
		return repository
}

func (this *postRepositoryImpl) init() {
		this.service = services.PostServiceOf()
}

func (this *postRepositoryImpl) Create() common.ResponseJson {
		var (
				ctx  = this.ctx.Ctx
				data = new(models.TravelNotes)
				id = this.ctx.GetSession(middlewares.AuthUserId)
		)
		err := json.Unmarshal(ctx.Input.RequestBody, data)
		if err != nil {
				err = this.ctx.ParseForm(data)
				if err != nil {
						return common.NewErrorResp(common.NewErrors(common.EmptyParamCode, err), "参数不足")
				}
		}
		data.UserId = id.(string)
		err = this.service.Create(data.Defaults())
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.ServiceFailed, err), "发布失败")
		}
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "id": data.Id.Hex()}, "发布成功")
}

func (this *postRepositoryImpl) Update() common.ResponseJson {
		var (
				id = this.ctx.GetString(":id")
		)
		if id == "" {
				return common.NewInvalidParametersResp("游记id缺失")
		}
		// @todo 更新
		return common.NewInDevResp(this.ctx.Ctx.Request.URL.String())
}

func (this *postRepositoryImpl) Lists(typ ...string) common.ResponseJson {
		var (
				meta     *models.Meta
				items    []*models.TravelNotes
				page, _  = this.ctx.GetInt("page", 1)
				count, _ = this.ctx.GetInt("count", 20)
				limit    = models.NewListParam(page, count)
				ty       = this.ctx.GetString("type")
		)
		if ty == "" && len(typ) != 0 {
				ty = typ[0]
		}
		switch ty {
		case "my":
				id := this.ctx.GetSession(middlewares.AuthUserId)
				items, meta = this.service.Lists(id.(string), limit)
		case "address":
				items, meta = this.service.ListByAddress(this.ctx.GetString(":address"), limit)
		case "tags":
				address := this.ctx.GetString("tags")
				items, meta = this.service.ListByAddress(address, limit)
		case "search":
				items, meta = this.service.Search(beego.M{}, limit)
		}
		if items != nil && len(items) > 0 && meta != nil {
				var arr []beego.M
				for _,item:=range items{
						arr = append(arr,item.M(filterEmpty))
				}
				return common.NewSuccessResp(beego.M{"items": arr, "meta": meta}, "罗列成功")
		}
		return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
}

func (this *postRepositoryImpl) GetById(id ...string) common.ResponseJson {
		if len(id) == 0 {
				id = append(id, this.ctx.GetString(":id"))
		}
		var data = this.service.GetById(id[0])
		if data == nil || data.DeletedAt != 0 {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		return common.NewSuccessResp(data.M(filterEmpty), "获取成功")
}

func (this *postRepositoryImpl) RemoveId(id ...string) common.ResponseJson {
		if len(id) == 0 {
				id = append(id, this.ctx.GetString(":id"))
		}
		var data = this.service.GetById(id[0])
		if data == nil {
				return common.NewFailedResp(common.RecordNotFound, common.RecordNotFoundError)
		}
		data.DeletedAt = time.Now().Unix()
		err := data.Save()
		if err == nil {
				common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix()}, "删除成功")
		}
		return common.NewFailedResp(common.RecordNotFound, "删除失败")
}
