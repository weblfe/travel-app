package services

import (
		"context"
		"errors"
		"fmt"
		"github.com/astaxie/beego"
		"github.com/astaxie/beego/config/env"
		"github.com/astaxie/beego/logs"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/qiniu/api.v7/v7/storage"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/models"
		"github.com/weblfe/travel-app/plugins"
		"io"
		"math"
		"os"
		"path/filepath"
		"strconv"
		"strings"
		"time"
)

type AttachmentService interface {
		SyncOssTask() int
		Remove(query beego.M) bool
		Get(mediaId string) *models.Attachment
		UpdateById(string, beego.M) error
		GetUrl(string) string
		GetById(id string) *models.Attachment
		GetAccessUrl(string) string
		Lists(page, count int) ([]*models.Attachment, *models.Meta)
		AutoCoverForVideo(attachment *models.Attachment, posts ...*models.TravelNotes) string
		Save(reader io.ReadCloser, extras ...beego.M) *models.Attachment
}

const (
		AttachTypeDoc         = models.AttachTypeDoc
		AttachTypeText        = models.AttachTypeText
		AttachTypeVideo       = models.AttachTypeVideo
		AttachTypeImage       = models.AttachTypeImage
		AttachTypeImageAvatar = models.AttachTypeImageAvatar
)

type AttachmentServiceImpl struct {
		BaseService
		model *models.AttachmentModel
}

func AttachmentServiceOf() AttachmentService {
		var service = new(AttachmentServiceImpl)
		service.Init()
		return service
}

func (this *AttachmentServiceImpl) Init() {
		this.model = models.AttachmentModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return AttachmentServiceOf()
		}
		this.init()
}

func (this *AttachmentServiceImpl) Remove(query beego.M) bool {
		if err := this.model.Remove(query, true); err == nil {
				return true
		}
		return false
}

func (this *AttachmentServiceImpl) Get(mediaId string) *models.Attachment {
		var attach = UrlTicketServiceOf().GetTicketInfoToSimple(mediaId)
		if attach == nil {
				var model, err = this.model.GetByMediaId(mediaId)
				if err == nil {
						return model
				}
				return nil
		}
		if model, err := this.model.GetByMediaId(attach.MediaId); model != nil && err == nil {
				return model
		}
		return nil
}

func (this *AttachmentServiceImpl) GetById(id string) *models.Attachment {
		var model, err = this.model.GetByMediaId(id)
		if model != nil && err == nil {
				return model
		}
		return nil
}

func (this *AttachmentServiceImpl) Save(reader io.ReadCloser, extras ...beego.M) *models.Attachment {
		var m = this.defaultsExtras(libs.MapMerge(extras...))
		if reader == nil {
				return nil
		}
		model := this.save(reader, m)
		if model == nil {
				return nil
		}
		if ty, ok := m["type"]; ok && ty != nil && ty != "" {
				tyName := ty.(string)
				if tyName != "file" && tyName != "files" {
						model.FileType = tyName
				}
		}
		if this.Create(model) {
				return model
		}
		return nil
}

func (this *AttachmentServiceImpl) Create(attach *models.Attachment) bool {
		if attach == nil {
				return false
		}
		attach = this.onlySaveOne(attach)
		if err := this.model.Add(attach.Defaults()); err == nil {
				go this.after(attach)
				return true
		}
		return false
}

func (this *AttachmentServiceImpl) after(attachment *models.Attachment) {
		switch attachment.FileType {
		case "video":
				this.video(attachment)
		}
		var fs, err = os.Open(attachment.GetLocal())
		if err != nil {
				logs.Error(err)
				return
		}
		defer this.closer(fs)
		extras := attachment.M()
		extras["key"] = attachment.GetBase()
		this.Uploader(fs, extras)
}

func (this *AttachmentServiceImpl) video(attachment *models.Attachment) bool {
		var err error
		if attachment.Duration == 0 && attachment.Path != "" {
				attachment.Duration, err = libs.GetMp4FileDuration(attachment.GetLocal())
				if err != nil {
						return false
				}
				err = this.model.Update(bson.M{"_id": attachment.Id}, beego.M{"duration": attachment.Duration})
				if err != nil {
						logs.Error(err)
				}
		}
		if "" == this.AutoCoverForVideo(attachment) {
				return false
		}
		return true
}

func (this *AttachmentServiceImpl) delete(fs string) {
		if fs != "" {
				err := os.Remove(fs)
				if err != nil {
						logs.Error(err)
				}
		}
}

// 文件仅保存一份
func (this *AttachmentServiceImpl) onlySaveOne(attach *models.Attachment) *models.Attachment {
		// 开关文件保存一份 通过hash
		if env.Get("ATTACHMENT_ONLY_ONE_SAVE", "on") == "off" {
				return attach
		}
		if attach.Hash != "" {
				oldAttach := this.GetByHash(attach.Hash)
				if oldAttach != nil && attach.Path != "" && attach.FileName != "" {
						existsFs := filepath.Join(oldAttach.Path, oldAttach.FileName)
						if !libs.IsExits(existsFs) {
								return attach
						}
						fs := filepath.Join(attach.Path, attach.FileName)
						defer this.delete(fs)
						attach.ExtrasInfo["originSavePath"] = attach.Path
						attach.ExtrasInfo["originSaveFileName"] = attach.FileName
						attach.Path = oldAttach.Path
						attach.FileName = oldAttach.FileName
						attach.Duration = oldAttach.Duration
				}
		}
		return attach
}

func (this *AttachmentServiceImpl) defaultsExtras(m beego.M) beego.M {
		_, ok := m["path"]
		if len(m) == 0 || !ok {
				m["path"] = this.getAttachmentPath()
		}
		return m
}

func (this *AttachmentServiceImpl) getAttachmentPath() string {
		var (
				year, month, day = time.Now().Date()
				date             = fmt.Sprintf("%d-%d-%d", year, month, day)
		)
		return PathsServiceOf().StoragePath("/" + date)
}

func (this *AttachmentServiceImpl) save(reader io.ReadCloser, extras beego.M) *models.Attachment {
		var (
				path           = extras["path"]
				oss, ossBucket = extras["oss"], extras["oss_bucket"]
		)
		if reader == nil {
				return nil
		}
		if path != "" {
				extras["path"] = path
				_ = os.MkdirAll(path.(string), os.ModePerm)
				res, ok := GetFileSystem().SaveByReader(reader, extras)
				if ok && len(res) > 0 {
						return models.NewAttachment().Load(libs.MapMerge(res, extras)).Defaults()
				}
				return nil
		}
		if oss != "" && ossBucket != "" {
				return this.Uploader(reader, extras)
		}
		return nil
}

// oss 上传器
func (this *AttachmentServiceImpl) Uploader(reader io.ReadCloser, extras beego.M) *models.Attachment {
		var (
				id  = getId(extras)
				key = extras["key"]
		)
		// 空流， 未知附件
		if reader == nil || id == "" {
				return nil
		}

		var params = plugins.OssParams{
				TypeName:  getType(extras),
				Storage:   storage.Config{},
				Reader:    reader,
				Size:      getSize(extras),
				Key:       key.(string),
				Extras:    nil,
				PutPolicy: nil,
		}
		params.Result = result()
		params.Extras = putExtras(extras)
		params.PutPolicy = putPolicy(extras)
		var uploader = plugins.GetOSS().CreateUploader(&params, func(cfg *storage.Config) {
				cfg.Zone = &storage.ZoneHuanan
				cfg.UseHTTPS = libs.Boolean(plugins.GetQinNiuProperty("USE_HTTPS", "false"))
		})

		var res, err = uploader(context.Background())

		if err != nil {
				logs.Error(err)
				return nil
		}
		var (
				oss    string
				attr   = models.NewAttachment()
				bucket = params.PutPolicy.Scope
		)
		if oss == "" {
				oss = params.Provider
		}
		if strings.Contains(bucket, ":") {
				var keys = strings.Split(bucket, ":")
				bucket = keys[0]
		}
		body, ok := res.(*ReturnBody)
		if !ok || body == nil {
				logs.Error(errors.New("empty oss return body"))
				return nil
		}
		err = this.model.FindOne(beego.M{"_id": id}, attr)
		if err != nil {
				logs.Error(err)
				return nil
		}
		if body.Id != id.Hex() {
				logs.Error(errors.New("id not matched"))
				return nil
		}
		// 更新记录
		attr.Oss = oss
		attr.Cdn = oss
		attr.OssBucket = bucket
		attr.Width, _ = strconv.Atoi(body.Width)
		attr.Height, _ = strconv.Atoi(body.Height)
		attr.CdnUrl = plugins.GetOssAccessUrl(body.Path, oss, bucket)
		if attr.Duration == 0 && body.Duration != "" {
				t, _ := strconv.ParseFloat(body.Duration, 10)
				if t != 0 {
						d := int64(math.Floor(t))
						attr.Duration = time.Duration(d) * time.Second
				}
		}
		if attr.Size == 0 && body.Size != "" {
				n, err1 := strconv.Atoi(body.Size)
				if err1 != nil {
						logs.Error(err1)
				}
				attr.Size = int64(n)
		}
		attr.ExtrasInfo["oss"] = body
		attr.UpdatedAt = time.Now().Local()
		err = this.model.Update(beego.M{"_id": attr.Id}, attr)
		if err != nil {
				logs.Error(err)
		} else {
				logs.Info("cdn put success id: " + attr.Id.Hex())
				// @todo 自预热
		}
		return attr
}

func (this *AttachmentServiceImpl) SaveToOssById(id string) error {
		var data = this.GetById(id)
		if data == nil {
				return errors.New("not exists id")
		}
		var fs = data.GetLocal()
		if fs == "" {
				return errors.New("file not exists")
		}
		var fd, err = os.Open(fs)
		if err != nil {
				return err
		}
		if nil != this.Uploader(fd, data.M()) {
				return nil
		}
		return errors.New("save failed ,id: " + id)
}

func (this *AttachmentServiceImpl) GetByHash(hash string) *models.Attachment {
		var attach = models.NewAttachment()
		if err := this.model.GetByKey("hash", hash, attach); err == nil {
				return attach
		}
		return nil
}

func (this *AttachmentServiceImpl) UpdateById(id string, update beego.M) error {
		if len(update) == 0 || id == "" {
				return common.NewErrors(common.InvalidParametersCode, "更新参数不能为空")
		}
		return this.model.UpdateById(id, update)
}

func (this *AttachmentServiceImpl) GetUrl(mediaId string) string {
		var data = this.Get(mediaId)
		if data == nil {
				return ""
		}
		if data.CdnUrl != "" {
				return data.CdnUrl
		}
		return data.Url
}

func (this *AttachmentServiceImpl) Lists(page, count int) ([]*models.Attachment, *models.Meta) {
		var (
				meta  = models.NewMeta()
				query = bson.M{"status": models.StatusOk}
				items = make([]*models.Attachment, count)
		)
		items = items[:0]
		meta.Count = count
		meta.Page = page

		var err = this.model.NewQuery(query).Limit(count).Skip((page - 1) * count).All(&items)
		if err == nil {
				meta.Total, err = this.model.NewQuery(query).Count()
				if err != nil {
						logs.Error(err)
				}
		}
		meta.Size = len(items)
		meta.Boot()
		return items, meta
}

func (this *AttachmentServiceImpl) GetAccessUrl(mediaId string) string {
		var data = this.Get(mediaId)
		if data == nil {
				return ""
		}
		return UrlTicketServiceOf().GetTicketUrlByAttach(data)
}

func (this *AttachmentServiceImpl) AutoCoverForVideo(attachment *models.Attachment, posts ...*models.TravelNotes) string {
		if attachment == nil || attachment.FileType != AttachTypeVideo {
				return ""
		}
		if attachment.CoverId != "" {
				return attachment.CoverId.Hex()
		}
		var fs = attachment.GetLocal()
		if !libs.IsExits(fs) {
				return ""
		}
		var (
				ext         = filepath.Ext(fs)
				name        = fmt.Sprintf("%d.%s", time.Now().Unix(), "jpg")
				storageName = strings.Replace(fs, ext, name, 1)
		)
		if plugins.ScreenShot(fs, storageName) {
				fd, _ := os.Open(storageName)
				defer this.closer(fd)
				data := beego.M{
						"userId":  attachment.UserId,
						"referId": attachment.Id.Hex(), "fileType": AttachTypeImage,
						"referName": models.AttachmentTable,
						"filename":  filepath.Base(fd.Name()),
				}
				image := this.Save(fd, data)
				if image == nil {
						return ""
				}
				err := this.UpdateById(attachment.Id.Hex(), beego.M{"coverId": image.Id})
				if err != nil {
						logs.Error(err)
				}
				if len(posts) > 0 {
						posts[0].Images = append(posts[0].Images, image.Id.Hex())
						return image.Id.Hex()
				}
				if attachment.ReferId != "" && attachment.ReferName == models.TravelNotesTable {
						_ = PostServiceOf().UpdateById(attachment.ReferId, beego.M{"images": []string{image.Id.Hex()}})
				}
				return image.Id.Hex()
		}
		logs.Info("ScreenShot failed", fs)
		return ""
}

func (this *AttachmentServiceImpl) closer(closer io.Closer) {
		var err = closer.Close()
		if err != nil {
				logs.Error(err)
		}
}

func (this *AttachmentServiceImpl) SyncOssTask() int {
		var (
				count int
				query = bson.M{"cdnUrl": beego.M{"$in": []interface{}{nil, ""}, "$exists": false}, "status": models.StatusOk}
		)
		var (
				iter       = this.model.NewQuery(query)
				total, err = this.model.NewQuery(query).Count()
		)
		if err != nil {
				logs.Error(err)
				return 0
		}
		if iter == nil {
				logs.Warn("iter for SyncOssTask nil ")
				return total
		}
		// 异步循环
		go this.ossAsyncTask(iter, &count)
		return total
}

// 异步同步任务
func (this *AttachmentServiceImpl) ossAsyncTask(iter *mgo.Query, count *int) {
		var (
				page = 1
				size  = 100
				items = make([]*models.Attachment, 10)
		)
		// iter.Iter()
		for {
				items = items[:0]
				// iter.Done()
				err := iter.Skip((page-1)*size).Limit(size).All(&items)
				if err != nil {
						logs.Error(err)
						break
				}
				if len(items) == 0 {
						logs.Warn("empty attache for async task")
						break
				}
				logs.Info(fmt.Sprintf("count: %d", len(items)))
				for _, it := range items {
						logs.Info("start..." + it.Id.Hex(),it.CdnUrl)
						if it.CdnUrl != ""  {
								logs.Info("cdnUrl: " + it.CdnUrl)
								continue
						}
						logs.Info("filename..." + it.GetLocal())
						if it.Size >= 251658240 {
								logs.Info("filename size to lager" , it.Size)
							continue
						}
						fs, err1 := os.Open(it.GetLocal())
						if err1 != nil {
								logs.Error(err1)
								continue
						}
						data := it.M()
						data["key"] = it.GetBase()
						if this.Uploader(fs, data) != nil {
								*count++
								logs.Info("success..." + it.Id.Hex())
						}else{
								logs.Info("failed..." + it.Id.Hex())
						}
						this.closer(fs)
				}
		}
}

func getSize(extras beego.M) int64 {
		if v, ok := extras["size"]; ok {
				return v.(int64)
		}
		return 0
}

func getMediaId(v interface{}) string {
		if v == nil {
				return bson.NewObjectId().Hex()
		}
		if id, ok := v.(string); ok {
				return id
		}
		if id, ok := v.(bson.ObjectId); ok {
				return id.Hex()
		}
		return bson.NewObjectId().Hex()
}

// 上传类型
func getType(extras beego.M) string {
		if v, ok := extras["type"]; ok {
				var t = v.(string)
				if t == AttachTypeImage {
						return plugins.QinNiuBucketImg
				}
				if t == AttachTypeVideo {
						return plugins.QinNiuBucketVideo
				}
		}
		if v, ok := extras["fileType"]; ok {
				var t = v.(string)
				if t == AttachTypeImage {
						return plugins.QinNiuBucketImg
				}
				if t == AttachTypeVideo {
						return plugins.QinNiuBucketVideo
				}
		}
		return ""
}

// 返回数据结构
func result() *ReturnBody {
		return new(ReturnBody)
}

// 扩展数据
func putExtras(extras beego.M) *storage.PutExtra {
		var (
				params = &storage.PutExtra{
						Params: map[string]string{
								"x:app":       os.Getenv("APP_NAME"),
								"x:uuid":      getMediaId(extras["id"]),
								"x:filename":  getFileName(extras["filename"]),
								"x:timestamp": fmt.Sprintf("%v", time.Now().Unix()),
						},
				}
		)
		return params
}

// 上传策略
func putPolicy(extras beego.M) *storage.PutPolicy {
		var (
				ty     = getType(extras)
				params = &storage.PutPolicy{ReturnBody: ""}
				body   = map[string]string{
						"name":        "$(fname)",
						"size":        "$(fsize)",
						"key":         "$(key)",
						"type":        "$(mimeType)",
						"hash":        "$(etag)",
						"x:uuid":      "$(x:uuid)",
						"x:app":       "$(x:app)",
						"x:type":      "$(x:type)",
						"x:filename":  "$(x:filename)",
						"x:timestamp": "$(x:timestamp)",
				}
		)
		if ty == plugins.QinNiuBucketImg {
				body["w"] = "$(imageInfo.width)"
				body["h"] = "$(imageInfo.height)"
				body["color"] = "$(exif.ColorSpace.val)"
		}
		if ty == plugins.QinNiuBucketVideo {
				body["w"] = "$(avinfo.video.width)"
				body["h"] = "$(avinfo.video.height)"
				body["duration"] = "$(avinfo.video.duration)"
		}
		var info, err = libs.Json().Marshal(body)
		if err != nil {
				logs.Error(err)
		}
		if len(info) > 0 {
				params.ReturnBody = string(info)
		}
		params.Expires = getExpires(extras)
		return params
}

// token过期时长
func getExpires(extras beego.M) uint64 {
		if e, ok := extras["expires"]; ok {
				if n, ok := e.(int64); ok {
						return uint64(n)
				}
		}
		return 1800
}

func getFileName(v interface{}) string {
		if v == nil {
				return ""
		}
		if name, ok := v.(string); ok {
				return name
		}
		return ""
}

func getId(extras beego.M) bson.ObjectId {
		if id, ok := extras["id"]; ok {
				if id == "" || id == nil {
						return ""
				}
				if str, ok := id.(string); ok {
						return bson.ObjectIdHex(str)
				}
				if obj, ok := id.(bson.ObjectId); ok {
						return obj
				}
		}
		return ""
}

// 返回数据
type ReturnBody struct {
		Hash      string `json:"hash"`
		Size      string `json:"size"`
		Ty        string `json:"x:type"`
		FileType  string `json:"type"`
		Timestamp string `json:"x:timestamp"`
		App       string `json:"x:app"`
		Name      string `json:"name"`
		Id        string `json:"x:uuid"`
		Path      string `json:"key"`
		Color     string `json:"color,omitempty"`
		FileName  string `json:"x:filename,omitempty"`
		Width     string `json:"w,omitempty"`
		Height    string `json:"h,omitempty"`
		Duration  string `json:"duration,omitempty"`
}

// 声频 audio
type ReturnAudio struct {
		BitRate    int64       `json:"bit_rate"`
		Channels   string      `json:"channels"`
		CodeName   string      `json:"code_name"`
		CodecType  string      `json:"codec_type"`
		Duration   string      `json:"duration"`
		Height     string      `json:"height"`
		Width      string      `json:"width"`
		Index      string      `json:"index"`
		NbFrames   string      `json:"nb_frames"`
		SampleFmt  string      `json:"sample_fmt"`
		RFrameRate string      `json:"r_frame_rate"`
		SampleRate string      `json:"sample_rate"`
		StartTime  string      `json:"start_time"`
		Tags       *ReturnTags `json:"tags"`
}

// 视频 video
type ReturnVideo struct {
		BitRate            int64       `json:"bit_rate"`
		CodeName           string      `json:"code_name"`
		CodecType          string      `json:"codec_type"`
		DisplayAspectRatio string      `json:"display_aspect_ratio"`
		Duration           string      `json:"duration"`
		Height             string      `json:"height"`
		Width              string      `json:"width"`
		Index              string      `json:"index"`
		NbFrames           string      `json:"nb_frames"`
		PixFmt             string      `json:"pix_fmt"`
		RFrameRate         string      `json:"r_frame_rate"`
		SampleAspectRatio  string      `json:"sample_aspect_ratio"`
		StartTime          string      `json:"start_time"`
		Tags               *ReturnTags `json:"tags"`
}

// 格式
type ReturnFormat struct {
		BitRate        int64       `json:"bit_rate"`
		Duration       float64     `json:"duration"`
		FormatLongName string      `json:"format_long_name"`
		FormatName     string      `json:"format_name"`
		NbFrames       int         `json:"nb_frames"`
		Size           int64       `json:"size"`
		StartTime      string      `json:"start_time"`
		Tags           *ReturnTags `json:"tags"`
}

// tags
type ReturnTags struct {
		CreationTime string `json:"creation_time"`
}

// 	七牛云
//  音频
//  "audio" : {
//        "bit_rate":"64028",
//        "channels":1,
//        "codec_name":"mp3",
//        "codec_type":"audio",
//        "duration":"30.105556",
//        "index":1,
//        "nb_frames":"1153",
//        "r_frame_rate":"0/0",
//        "sample_fmt":"s16p",
//        "sample_rate":"44100",
//        "start_time":"0.000000",
//        "tags":{
//            "creation_time":"2012-10-21 01:13:54"
//        }
//    }
//
//   "format" : {
//        "bit_rate":"918325",
//        "duration":"30.106000",
//        "format_long_name":"QuickTime / MOV",
//        "format_name":"mov,mp4,m4a,3gp,3g2,mj2",
//        "nb_streams":2,
//        "size":"3455888",
//        "start_time":"0.000000",
//        "tags":{
//            "creation_time":"2012-10-21 01:13:54"
//        }
//    }
//
//    视频
//    "video": {
//        "bit_rate":"856559",
//        "codec_name":"h264",
//        "codec_type":"video",
//        "display_aspect_ratio":"4:3",
//        "duration":"29.791667",
//        "height":480,
//        "index":0,
//        "nb_frames":"715",
//        "pix_fmt":"yuv420p",
//        "r_frame_rate":"24/1",
//        "sample_aspect_ratio":"1:1",
//        "start_time":"0.000000",
//        "tags":{
//            "creation_time":"2012-10-21 01:13:54"
//        },
//        "width":640
//    }
//
