package services

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
)

type CommentService interface {
		Commit(data *models.Comment) error
		IncrThumbsUp(id string, incr int) error
}

type commentServiceImpl struct {
		BaseService
		model *models.CommentModel
}

func CommentServiceOf() CommentService {
		var service = new(commentServiceImpl)
		service.Init()
		return service
}

func (this *commentServiceImpl) Init() {
		this.init()
		this.model = models.CommentModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return CommentServiceOf()
		}
}

// 提交保持评论
func (this *commentServiceImpl) Commit(data *models.Comment) error {
		var err = this.model.Add(data.Defaults())
		if err == nil {
				go this.addSuccessAfter(data)
		}
		return err
}

// 更新评论数
func (this *commentServiceImpl) addSuccessAfter(comment *models.Comment) {
		if comment.PostId != "" {
				_ = PostServiceOf().IncrComment(comment.PostId.Hex())
		}
}

// 评论列表
func (this *commentServiceImpl) Lists(query beego.M, limit models.ListsParams) ([]*models.Comment, *models.Meta) {
		var (
				err   error
				meta  = new(models.Meta)
				items []*models.Comment
		)
		meta.Count = limit.Count()
		meta.Page = limit.Page()
		meta.Total, err = this.model.Lists(query, &items, limit)
		if err == nil {
				meta.Boot()
				return items, meta
		}
		return nil, meta
}

// 更新点赞数量
func (this *commentServiceImpl) IncrThumbsUp(id string, incr int) error {
		return this.model.IncrBy(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"thumbsUpNum": incr})
}
