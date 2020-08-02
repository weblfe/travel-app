package repositories

import (
		"github.com/astaxie/beego"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
		"github.com/weblfe/travel-app/transports"
		"time"
)

type CommentRepository interface {
		Create() common.ResponseJson
		Detail() common.ResponseJson
		Lists() common.ResponseJson
}

type commentRepository struct {
		ctx     common.BaseRequestContext
		service services.CommentService
}

func NewCommentRepository(ctx common.BaseRequestContext) CommentRepository {
		var repository = new(commentRepository)
		repository.ctx = ctx
		repository.init()
		return repository
}

func (this *commentRepository) init() {
		this.service = services.CommentServiceOf()
}

func (this *commentRepository) GetDto() *DtoRepository {
		return GetDtoRepository()
}

func (this *commentRepository) Create() common.ResponseJson {
		var (
				err     error
				comment *models.Comment
				data    = transports.NewComment(this.ctx.GetInput())
		)
		if data == nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError, common.InvalidParametersError), "发布评论失败")
		}
		comment = data.Decode()
		if comment.Content != "" {
				comment.UserId = getUserId(this.ctx)
		}
		// 自动过滤敏感词
		comment.Content = models.GetDfaInstance().ChangeSensitiveWords(comment.Content)
		comment.Status = models.StatusAuditPass
		err = this.service.Commit(comment)
		if err != nil {
				return common.NewErrorResp(common.NewErrors(err, common.ServiceFailed), "发布评论失败")
		}
		return common.NewSuccessResp(beego.M{"timestamp": time.Now().Unix(), "id": comment.Id.Hex()}, "发布成功")
}

func (this *commentRepository) Detail() common.ResponseJson {
		var (
				comment *models.Comment
				id, ok  = this.ctx.GetParam(":id")
		)
		if !ok {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError, common.InvalidParametersError), "评论Id缺失")
		}
		comment = this.service.GetById(id.(string))
		if comment == nil {
				return common.NewErrorResp(common.NewErrors(common.RecordNotFound, common.RecordNotFoundError), "评论不存在")
		}
		return common.NewSuccessResp(comment.M(this.getTransports()), "获取评论详情成功")
}

func (this *commentRepository) Lists() common.ResponseJson {
		var (
				meta     *models.Meta
				comments []*models.Comment
				typ, id  = this.ctx.GetString("targetType"), this.ctx.GetString("targetId")
				page     = models.NewListParam(this.ctx.GetInt("page", 1), this.ctx.GetInt("count", 20))
		)
		if id == "" {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError, common.InvalidParametersError), "参数缺失")
		}
		// 默认未 游记评论
		if typ == "" {
				typ = "post"
		}
		comments, meta = this.service.Lists(typ, id, page)
		if comments == nil || meta == nil {
				return common.NewErrorResp(common.NewErrors(common.InvalidParametersError, common.InvalidParametersError), "获取列表失败")
		}
		return common.NewSuccessResp(beego.M{"items": this.each(comments), "meta": meta}, "发布成功")
}

// 评论转换
func (this *commentRepository) each(items []*models.Comment) []beego.M {
		var (
				lists []beego.M
				trans = this.getTransports()
		)
		for _, it := range items {
				data := it.M(trans)
				this.addReviews(&data, it)
				lists = append(lists, data)
		}
		return lists
}

// 添加回复
func (this *commentRepository) addReviews(data *beego.M, it *models.Comment) {
		(*data)["reviews"] = []*models.Comment{}
		if it.ReviewNum <= 0 {
				return
		}
		var reviews, count = this.service.GetReviews(it.Id.Hex())
		if count > 0 {
				(*data)["reviews"] = this.each(reviews)
		}
}

func (this *commentRepository) getTransports() func(m beego.M) beego.M {
		var dto = this.GetDto()
		return func(m beego.M) beego.M {
				m = transforms.FieldsFilter([]string{"deletedAt", "updatedAt", "tags", "refersIds"})(m)
				m = this.appendUser(m, dto)
				return m
		}
}

func (this *commentRepository) appendUser(m beego.M, dto *DtoRepository) beego.M {
		if value, ok := m["userId"]; ok {
				var (
						v    interface{}
						key  string
						user *BaseUser
						id   = value.(string)
				)
				if id != "" {
						key = dto.Key(key)
				}
				if key == "" {
						return m
				}
				v = dto.Get(id)
				if v != nil {
						user = v.(*BaseUser)
				}
				if user == nil {
						user = dto.GetUserById(id)
						if user != nil {
								dto.Cache(key, user)
						}
				}
				if user != nil {
						m["userInfo"] = user
				}
		}
		return m
}
