package services

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
)

type PostService interface {
		Audit(...string) bool
		IncrComment(id string) error
		Exists(query bson.M) bool
		Create(notes *models.TravelNotes) error
		GetById(id string) *models.TravelNotes
		IncrThumbsUp(id string, incr int) error
		ListsQuery(query bson.M, limit models.ListsParams,sort...string) ([]*models.TravelNotes, *models.Meta)
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
		var images []string
		if notes.Images != nil && len(notes.Images) > 0 {
				images = notes.Images[:]
		}
		// 矫正类型
		if notes.Videos != nil && len(notes.Videos) > 0 {
				notes.Type = PostTypeVideo
		}
		var err = this.postModel.Add(notes)
		if err == nil {
				// 异步更新 附件归属
				notes.Images = images
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

func (this *TravelPostServiceImpl) Search(search beego.M, page models.ListsParams) ([]*models.TravelNotes, *models.Meta) {
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
		)
		meta.Page = page.Page()
		meta.Size = len(lists)
		meta.Count = page.Count()
		meta.Total, err = this.postModel.Lists(search, lists, page)
		if err == nil {
				meta.Boot()
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
		return this.postModel.Incr(id, "thumbsUpNum", incr)
}

func (this *TravelPostServiceImpl) Audit(id ...string) bool {
		if len(id) == 0 {
				return false
		}
		var arr []bson.ObjectId
		for _, v := range id {
				arr = append(arr, bson.ObjectIdHex(v))
		}
		var err = this.postModel.Update(bson.M{"_id": bson.M{"$in": arr}}, bson.M{"status": models.StatusAuditPass})
		if err == nil {
				return true
		}
		return false
}

func (this *TravelPostServiceImpl) Exists(query bson.M) bool {
		return this.postModel.Exists(query)
}
