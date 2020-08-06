package services

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/models"
		"time"
)

type CommentService interface {
		GetById(string) *models.Comment
		Commit(data *models.Comment) error
		IncrThumbsUp(id string, incr int) error
		GetReviews(string) ([]*models.Comment, int)
		Lists(typ, id string, page models.ListsParams, extras ...beego.M) ([]*models.Comment, *models.Meta)
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
		data.Defaults()
		if data.TargetId == "" {
				return errors.New("empty targetId")
		}
		if data.Content == "" {
				return errors.New("empty content")
		}
		if data.UserId == "" {
				return errors.New("empty userId")
		}
		if data.TargetType =="" {
				return errors.New("empty targetType")
		}
		// 检查是否自己给自己评论，回复
		/*if err:=this.check(data);err!=nil {
				return err
		}*/
		var err = this.model.Add(data)
		if err == nil {
				go this.addSuccessAfter(data)
		}
		return err
}

func (this *commentServiceImpl)check(data *models.Comment) error {
		if PostServiceOf().Exists(bson.M{"userId":data.UserId,"_id":bson.ObjectIdHex(data.TargetId)}) {
				return errors.New("comment limit for self")
		}
		if this.Exists(bson.M{"userId":data.UserId,"_id":bson.ObjectIdHex(data.TargetId)}) {
				return errors.New("review limit for self")
		}
		return nil
}

func (this *commentServiceImpl)Exists(query bson.M) bool {
		return this.model.Exists(query)
}

// 更新评论数
func (this *commentServiceImpl) addSuccessAfter(comment *models.Comment) {
		_ = this.resolverRefersId(comment)
		_ = this.IncrCommentForPost(comment)
}

// 评论数据增加更新
func (this *commentServiceImpl) IncrCommentForPost(comment *models.Comment) error {
		if len(comment.RefersIds) >= 2 && comment.Status == models.StatusAuditPass {
				// 评论 的评论数统计
				if comment.TargetType == "comment" {
						_ = this.IncrCommentReviews(comment.RefersIds[1])
				}
		}
		// 作品评论数据统计
		if comment.TargetType == "post" && comment.Status == models.StatusAuditPass {
				if len(comment.RefersIds) >0 {
					_ = PostServiceOf().IncrComment(comment.RefersIds[0])
				}
		}
		return nil
}

// 更新评论 回复数
func (this *commentServiceImpl) IncrCommentReviews(id string) error {
		var data = this.GetById(id)
		if data == nil {
				return errors.New("not found comment id:" + id)
		}
		data.ReviewNum++
		return this.model.Update(beego.M{"_id": data.Id}, beego.M{"reviewNum": data.ReviewNum, "updatedAt": time.Now().Local()})
}

// 构建关联
func (this *commentServiceImpl) resolverRefersId(comment *models.Comment) error {
		if comment.TargetType == "comment" && comment.TargetId != "" {
				var data = this.GetById(comment.TargetId)
				if data == nil {
						return errors.New("not found refers id")
				}
				if data.RefersIds == nil {
						data.RefersIds = []string{}
				}
				if data.TargetType == "post" {
						comment.RefersIds = append([]string{data.TargetId}, comment.RefersIds...)
				}
				if data.TargetType == "comment" && len(data.RefersIds) >= 1 {
						comment.RefersIds = append([]string{data.RefersIds[0]}, comment.RefersIds...)
				}
				return this.model.Update(bson.M{"_id": comment.Id}, bson.M{"refersIds": comment.RefersIds, "updatedAt": time.Now().Local()})
		}
		return nil
}

// 评论列表
func (this *commentServiceImpl) Lists(targetType, targetId string, page models.ListsParams, extras ...beego.M) ([]*models.Comment, *models.Meta) {
		var (
				err   error
				query = beego.M{}
				meta  = new(models.Meta)
				items []*models.Comment
		)
		if targetType != "" {
				query["targetType"] = targetType
		}
		if targetId != "" {
				query["targetId"] = targetId
		}
		if len(extras) > 0 {
				for _, it := range extras {
						query = models.Merger(query, it)
				}
		}
		meta.Count = page.Count()
		meta.Page = page.Page()
		this.model.UseSoftDelete()
		listQuery := this.model.ListsQuery(query, page)
		// desc createdAt
		err = listQuery.Sort("+createdAt").All(&items)
		if err == nil {
				meta.Size = len(items)
				meta.Total, _ = this.model.ListsQuery(query, nil).Count()
				return items, meta
		}
		return nil, meta
}

// 更新评论 点赞数量
func (this *commentServiceImpl) IncrThumbsUp(id string, incr int) error {
		return this.model.IncrBy(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"thumbsUpNum": incr})
}

// 通过ID获取
func (this *commentServiceImpl) GetById(id string) *models.Comment {
		var (
				err     error
				comment = models.NewComment()
		)
		this.model.UseSoftDelete()
		err = this.model.GetById(id, comment)
		if err == nil {
				return comment
		}
		return nil
}

// 获取评论的所有回复
func (this *commentServiceImpl) GetReviews(id string) ([]*models.Comment, int) {
		this.model.UseSoftDelete()
		var (
				err     error
				reviews = make([]*models.Comment, 2)
				query   = this.model.NewQuery(bson.M{"targetId": id, "targetType": "comment", "status": models.StatusAuditPass})
		)
		reviews = reviews[:0]
		err = query.Sort("+createdAt").All(&reviews)
		if err == nil {
				return reviews, len(reviews)
		}
		return reviews, 0
}
