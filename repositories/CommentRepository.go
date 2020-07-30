package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
)

type CommentRepository interface {
		 Create() common.ResponseJson
}

type commentRepository struct {
		ctx               common.BaseRequestContext
		service services.CommentService
}

func NewCommentRepository(ctx common.BaseRequestContext) CommentRepository  {
		var repository = new(commentRepository)
		repository.ctx = ctx
		return repository
}

func (this *commentRepository)Create() common.ResponseJson  {
		var (
				err error
			data = beego.M{}
			comment = models.NewComment()
		)
		err =this.ctx.JsonDecode(&data)
		if len(data) == 0 || err != nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError,common.InvalidParametersError),"发布评论失败")
		}
		comment.Load(data)
		err= this.service.Commit(comment)
		if err != nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError,common.InvalidParametersError),"发布评论失败")
		}
		return common.NewErrorResp(common.NewErrors("创建失败",-1))
}