package services

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"time"
)

type PostService interface {
		Audit(string, ...string) bool
		IncrComment(id string) error
		Exists(query bson.M) bool
		Create(notes *models.TravelNotes) error
		GetById(id string) *models.TravelNotes
		IncrThumbsUp(id string, incr int) error
		UpdateById(id string, data beego.M) error
		AutoVideoCoverImageTask(ids []string) int
		ListsQuery(query bson.M, limit models.ListsParams, sort ...string) ([]*models.TravelNotes, *models.Meta)
		GetRankingLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta)
		GetRecommendLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta)
		Lists(userId string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta)
		ListByTags(tags []string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta)
		ListByAddress(address string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta)
		Search(search beego.M, page models.ListsParams) ([]*models.TravelNotes, *models.Meta)
}

type TravelPostServiceImpl struct {
		BaseService
		postModel *models.PostsModel
}

const (
		PostTypeImage = 1
		PostTypeVideo = 2
)

func PostServiceOf() PostService {
		var service = new(TravelPostServiceImpl)
		service.Init()
		return service
}

func (this *TravelPostServiceImpl) Init() {
		this.init()
		this.postModel = models.PostsModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return PostServiceOf()
		}
}

func (this *TravelPostServiceImpl) Lists(userId string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta) {
		extras = append(extras, beego.M{"userId": userId})
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
				query = libs.MapMerge(extras...)
		)
		meta.Page = page.Page()
		meta.Count = page.Count()
		defer this.postModel.Release()
		this.postModel.UseSoftDelete()
		listQuery := this.postModel.ListsQuery(query, page)
		// desc createdAt
		err = listQuery.Sort("-createdAt").All(&lists)
		if err == nil {
				meta.Size = len(lists)
				meta.Total, _ = this.postModel.ListsQuery(query, nil).Count()
				meta.Boot()

				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) ListByTags(tags []string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
				query beego.M
		)
		if len(tags) > 0 {
				extras = append(extras, beego.M{"tags": beego.M{"$in": tags}})
		}
		query = libs.MapMerge(extras...)
		meta.Page = page.Page()
		meta.Size = page.Count()
		meta.Count = page.Count()
		this.postModel.UseSoftDelete()

		listQuery := this.postModel.ListsQuery(query, page)
		// desc createdAt
		err = listQuery.Sort("-createdAt").All(&lists)
		if err == nil {
				meta.Size = len(lists)
				meta.Total, _ = this.postModel.ListsQuery(query, nil).Count()
				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) ListByAddress(address string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta) {
		extras = append(extras, beego.M{"address": &bson.RegEx{Pattern: address, Options: "i"}})
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
				query = libs.MapMerge(extras...)
		)
		meta.Page = page.Page()
		meta.Size = page.Count()
		meta.Count = page.Count()
		this.postModel.UseSoftDelete()
		listQuery := this.postModel.ListsQuery(query, page)
		// desc createdAt
		err = listQuery.Sort("-createdAt").All(&lists)
		if err == nil {
				meta.Size = len(lists)
				meta.Total, _ = this.postModel.ListsQuery(query, nil).Count()
				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) GetRankingLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
		)
		meta.Page = limit.Page()
		meta.Size = limit.Count()
		meta.Count = limit.Count()
		this.postModel.UseSoftDelete()
		listQuery := this.postModel.ListsQuery(query, limit)
		// 排行版 点赞 ,评论 最多, 发布时间最新
		err = listQuery.Sort("-thumbsUpNum", "-commentNum", "-createdAt").All(&lists)
		if err == nil {
				meta.Size = len(lists)
				meta.Total, _ = this.postModel.ListsQuery(query, nil).Count()
				return lists, meta
		}
		return nil, meta
}

// 获取推荐
func (this *TravelPostServiceImpl) GetRecommendLists(query bson.M, limit models.ListsParams) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
		)
		meta.Page = limit.Page()
		meta.Size = limit.Count()
		meta.Count = limit.Count()
		this.postModel.UseSoftDelete()
		listQuery := this.postModel.ListsQuery(query, limit)
		// 打分最高 ,更新时间最新
		err = listQuery.Sort("-score", "-createdAt").All(&lists)
		if err == nil {
				meta.Size = len(lists)
				meta.Total, _ = this.postModel.ListsQuery(query, nil).Count()
				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) ListsQuery(query bson.M, limit models.ListsParams, sort ...string) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
		)
		meta.Page = limit.Page()
		meta.Size = limit.Count()
		meta.Count = limit.Count()
		this.postModel.UseSoftDelete()
		listQuery := this.postModel.ListsQuery(query, limit)
		if len(sort) != 0 {
				listQuery = listQuery.Sort(sort...)
		}
		// 打分最高 ,更新时间最新
		err = listQuery.All(&lists)
		if err == nil {
				meta.Size = len(lists)
				meta.Total, _ = this.postModel.ListsQuery(query, nil).Count()
				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) Create(notes *models.TravelNotes) error {
		if notes.Type == PostTypeImage && len(notes.Images) <= 0 {
				return common.NewErrors(common.InvalidParametersCode, "图片不能为空")
		}
		// 矫正类型
		if notes.Videos != nil && len(notes.Videos) > 0 {
				notes.Type = PostTypeVideo
		}
		var err = this.postModel.Add(notes)
		if err == nil {
				// 异步更新 附件归属
				go this.attachments(notes)
		}
		return err
}

func (this *TravelPostServiceImpl) GetById(id string) *models.TravelNotes {
		if id == "" {
				return nil
		}
		var (
				err  error
				data = new(models.TravelNotes)
		)
		err = this.postModel.GetById(id, data)
		if err == nil {
				return data
		}
		return nil
}

func (this *TravelPostServiceImpl) attachments(notes *models.TravelNotes) {

		if notes.Videos != nil && len(notes.Videos) > 0 && notes.Type == PostTypeVideo {
				var (
						service = AttachmentServiceOf()
						update  = beego.M{"referName": this.postModel.TableName(), "referId": notes.Id.Hex()}
				)
				for _, id := range notes.Videos {
						_ = service.UpdateById(id, update)
				}
		}
		var (
				service = AttachmentServiceOf()
				update  = beego.M{"referName": this.postModel.TableName(), "referId": notes.Id.Hex()}
		)
		if notes.Images != nil && len(notes.Images) > 0 {
				// 更新图片归属
				for _, id := range notes.Images {
						_ = service.UpdateById(id, update)
				}
				// 设置视频封面
				if notes.Type == PostTypeVideo && len(notes.Videos) == 1 {
						update["coverId"] = bson.ObjectIdHex(notes.Images[0])
				}
		}
		// 设置视频关联
		if notes.Type == PostTypeVideo && len(notes.Videos) == 1 {
				_ = service.UpdateById(notes.Videos[0], update)
		}
}

// 自动化获取封面
func (this *TravelPostServiceImpl) attachVideo(post *models.TravelNotes) bool {
		var (
				count   = 0
				service = AttachmentServiceOf()
				images  []string
		)
		images = images[:0]
		for _, id := range post.Videos {
				attach := service.GetById(id)
				if attach == nil {
						logs.Info("attach id not exists :", id)
						return false
				}
				imageId := service.AutoCoverForVideo(attach, post)
				if imageId != "" {
						logs.Info("AutoCoverForVideo :", imageId)
						count++
				}
		}
		if count == len(post.Images) && count > 0 {
				post.UpdatedAt = time.Now().Local()
				err := this.postModel.Update(bson.M{"_id": post.Id}, post)
				if err == nil {
						return true
				}
				logs.Error(err)
		}
		return false
}

func (this *TravelPostServiceImpl) Search(search beego.M, page models.ListsParams) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				lists = make([]*models.TravelNotes, 2)
				meta  = models.NewMeta()
		)

		lists = lists[:0]
		meta.Page = page.Page()
		meta.Count = page.Count()

		search["status"] = models.StatusAuditPass
		this.postModel.UseSoftDelete()
		meta.Total, err = this.postModel.Lists(search, &lists, page)
		if err == nil {
				meta.Boot()
				meta.Size = len(lists)
				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) IncrComment(id string) error {
		if id == "" {
				return errors.New("comment post id empty")
		}
		return this.postModel.Incr(id, "commentNum")
}

func (this *TravelPostServiceImpl) IncrThumbsUp(id string, incr int) error {
		if id == "" {
				return errors.New("thumbsUp post id empty")
		}
		var post = this.GetById(id)
		if post == nil {
				return errors.New("thumbsUp post id empty")
		}
		defer this.AfterIncr(post.UserId, incr)
		return this.postModel.Incr(id, "thumbsUpNum", incr)
}

func (this *TravelPostServiceImpl) Audit(typ string, ids ...string) bool {
		if len(ids) == 0 {
				return false
		}
		var (
				status int
				arr    []bson.ObjectId
		)
		// 类型判断
		switch typ {
		case "1":
				status = models.StatusAuditPass
		case "2":
				status = models.StatusAuditUnPass
		case "-1":
				status = models.StatusAuditNotPass
		case "0":
				status = models.StatusWaitAudit
		default:
				return false
		}
		for _, v := range ids {
				arr = append(arr, bson.ObjectIdHex(v))
		}
		var err = this.postModel.Update(bson.M{"_id": bson.M{"$in": arr}}, bson.M{"status": status})
		if err == nil {
				return true
		}
		return false
}

func (this *TravelPostServiceImpl) Exists(query bson.M) bool {
		return this.postModel.Exists(query)
}

func (this *TravelPostServiceImpl) AfterIncr(userId string, incr int) {
		var (
				service = UserServiceOf()
				user    = service.GetById(userId)
		)
		if user == nil {
				return
		}
		_ = service.IncrBy(bson.M{"_id": bson.ObjectIdHex(userId)}, "thumbsUpTotal", incr)
}

func (this *TravelPostServiceImpl) UpdateById(id string, data beego.M) error {
		if len(data) == 0 {
				return common.NewErrors(common.EmptyParamCode, "空更新")
		}
		data["updatedAt"] = time.Now().Local()
		return this.postModel.UpdateById(id, data)
}

func (this *TravelPostServiceImpl) AutoVideoCoverImageTask(ids []string) int {
		var count = 0
		for _, id := range ids {
				post := this.GetById(id)
				if post == nil {
						continue
				}
				if post.Type != PostTypeVideo {
						continue
				}
				if len(post.Videos) <= 0 {
						continue
				}
				if this.attachVideo(post) {
						count++
				}
		}
		logs.Info("auto video cover :", count)
		return count
}
