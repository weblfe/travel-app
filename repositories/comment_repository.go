package repositories

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/services"
		"github.com/weblfe/travel-app/transforms"
		"github.com/weblfe/travel-app/transports"
		"strings"
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
		if !this.Check(comment) {
				return common.NewErrorResp(common.NewErrors(err, common.ServiceFailed), "评论类型受限")
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

// 是否
func (this *commentRepository) IsPostAuthor(postId, userId string) bool {
		var post = services.PostServiceOf().GetById(postId)
		if post == nil {
				return false
		}
		if post.UserId == userId {
				return true
		}
		return false
}

// 评论检查
func (this *commentRepository) Check(comment *models.Comment) bool {
		if !comment.CheckType() {
				return false
		}
		comment.Content = strings.TrimSpace(comment.Content)
		// 空内容不可以评论
		if comment.Content == "" {
				return false
		}
		// 评论是否存在
		if comment.TargetType != models.CommentTargetTypeReview {
				return services.PostServiceOf().Exists(bson.M{"_id": bson.ObjectIdHex(comment.TargetId)})
		}
		var obj = this.service.GetById(comment.TargetId)
		if obj == nil {
				return false
		}
		// 评论的回复只有一级
		if obj.TargetType != models.CommentTargetTypeComment {
				return false
		}
		// 非作者 不可用评论回复
		if !this.IsPostAuthor(obj.TargetId, comment.UserId) {
				return false
		}
		return true
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
		defer this.GetDto().Flash()
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
				return common.NewErrorResp(common.NewErrors(common.NotFound, common.RecordNotFoundError), "获取列表失败")
		}
		defer this.GetDto().Flash()
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

// 添加
func (this *commentRepository) getTransports() func(m beego.M) beego.M {
		var dto = this.GetDto()
		return func(m beego.M) beego.M {
				m = transforms.FieldsFilter([]string{"deletedAt", "updatedAt", "tags", "refersIds"})(m)
				m = this.appendUser(m, dto)
				return m
		}
}

// 添加用户信息
func (this *commentRepository) appendUser(m beego.M, dto *DtoRepository) beego.M {
		var (
				targetId   = m["targetId"]
				targetType = m["targetType"]
				userId, ok = m["userId"]
		)
		m["isAuthor"] = false
		if ok && userId != nil && userId != "" {
				var (
						v    interface{}
						key  string
						user *BaseUser
						id   = userId.(string)
				)
				if id != "" {
						key = dto.Key(id)
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
						m["isAuthor"] = this.isAuthor(targetId.(string), targetType.(string), user.UserId)
				}
		}
		return m
}

// 是否作者本人
func (this *commentRepository) isAuthor(targetId, targetType, userId string) bool {
		if targetType == "" || targetId == "" || userId == "" {
				return false
		}
		if targetType == models.CommentTargetTypeComment {
				return this.IsPostAuthor(targetId, userId)
		}
		if targetType == models.CommentTargetTypeReview {
				var comment = this.service.GetById(targetId)
				if comment == nil {
						return false
				}
				if comment.TargetType == models.CommentTargetTypeReview {
						return false
				}
				return this.isAuthor(comment.TargetId, comment.TargetType, userId)
		}
		return true
}
