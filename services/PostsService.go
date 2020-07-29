package services

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
)

type PostService interface {
		Create(notes *models.TravelNotes) error
		GetById(id string) *models.TravelNotes
		IncrComment(id string) error
		IncrThumbsUp(id string, incr int) error
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
		meta.Size = len(lists)
		meta.Count = page.Count()
		defer this.postModel.Release()
		listQuery := this.postModel.ListsQuery(query, page)
		err = listQuery.Sort("-createdAt").All(&lists)
		if err == nil {
				meta.Total, _ = listQuery.Count()
				meta.Boot()

				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) ListByTags(tags []string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta) {
		extras = append(extras, beego.M{"$in": beego.M{"tags": tags}})
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
				query = libs.MapMerge(extras...)
		)
		meta.Page = page.Page()
		meta.Size = len(lists)
		meta.Count = page.Count()
		meta.Total, err = this.postModel.Lists(query, &lists, page)
		if err == nil {
				meta.Boot()
				return lists, meta
		}
		return nil, meta
}

func (this *TravelPostServiceImpl) ListByAddress(address string, page models.ListsParams, extras ...beego.M) ([]*models.TravelNotes, *models.Meta) {
		extras = append(extras, beego.M{"$regexp": beego.M{"address": address}})
		var (
				err   error
				lists []*models.TravelNotes
				meta  = models.NewMeta()
				query = libs.MapMerge(extras...)
		)
		meta.Page = page.Page()
		meta.Size = len(lists)
		meta.Count = page.Count()
		meta.Total, err = this.postModel.Lists(query, &lists, page)
		if err == nil {
				meta.Boot()
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
