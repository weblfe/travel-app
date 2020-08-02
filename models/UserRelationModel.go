package models

import (
		"errors"
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

// 用户关系记录
type UserRelation struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                    // ID
		UserId        string        `json:"userId" bson:"userId"`             // 用户ID
		TargetUserId  string        `json:"targetUserId" bson:"targetUserId"` // 关系目标用户ID
		TargetType    string        `json:"targetType" bson:"targetType"`     // 关系类型
		Status        int           `json:"status" bson:"status"`             // 状态
		Extras        beego.M       `json:"extras" bson:"extras"`             // 扩展信息
		Tags          []string      `json:"tags" bson:"tags"`                 // 分类tags
		Versions      []string      `json:"versions" bson:"versions"`         // 状态版本变化历史 ["1-2020-08-01","2-2020-10-01"]
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`       // 创建时间
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`       // 更新时间
		DeletedAt     int64         `json:"deletedAt" bson:"deletedAt"`       // 删除时间
		dataClassImpl `bson:",omitempty"  json:",omitempty"`                  // 工具类
}

type UserRelationModel struct {
		BaseModel
}

const (
		UserRelationTable = "user_relations" // 表名
		TargetTypeFriend  = "friend"         // 朋友关系
		StatusOk          = 1                // 状态 正常
		StatusUnKnown     = 0                // 状态 初始化
		StatusCancel      = 2                // 状态 取消
)

func NewUserRelation() *UserRelation {
		var data = new(UserRelation)
		data.Init()
		return data
}

func UserRelationModelOf() *UserRelationModel {
		var model = new(UserRelationModel)
		return model.init()
}

func (this *UserRelation) data() beego.M {
		return beego.M{
				"id":           this.Id.Hex(),
				"userId":       this.UserId,
				"targetUserId": this.TargetUserId,
				"targetType":   this.TargetType,
				"status":       this.Status,
				"extras":       this.Extras,
				"tags":         this.Tags,
				"tagsText":     this.getTagsText(),
				"versions":     this.getVersions(),
				"createdAt":    this.CreatedAt.Unix(),
				"updatedAt":    this.UpdatedAt.Unix(),
				"deletedAt":    this.DeletedAt,
		}
}

func (this *UserRelation) getTags() []string {
		if this.Tags == nil {
				return []string{}
		}
		return this.Tags
}

func (this *UserRelation) getTagsText() []string {
		if this.Tags == nil {
				return []string{}
		}
		// @todo
		return []string{}
}

func (this *UserRelation) getVersions() []string {
		if this.Versions == nil {
				return []string{}
		}
		return this.Versions
}

func (this *UserRelation) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.defaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *UserRelation) defaults() {
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.Status == 0 {
				this.Status = 1
		}
		if this.TargetType == "" {
				this.TargetType = TargetTypeFriend
		}
}

func (this *UserRelation) save() error {
		var (
				model = UserRelationModelOf()
				data  = model.GetByUnique(this.data())
		)
		if data == nil {
				this.InitDefault()
				return model.Add(this)
		}
		return model.Update(bson.M{"_id": data.Id}, this.M(func(m beego.M) beego.M {
				delete(m, "id")
				delete(m, "createdAt")
				delete(m, "tagsText")
				m["updatedAt"] = time.Now().Local()
				return m
		}))
}

func (this *UserRelation) Set(key string, v interface{}) *UserRelation {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "userId":
				this.SetString(&this.UserId, v)
		case "targetUserId":
				this.SetString(&this.TargetUserId, v)
		case "targetType":
				this.SetString(&this.TargetType, v)
		case "status":
				this.SetNumInt(&this.Status, v)
		case "extras":
				this.SetMapper(&this.Extras, v)
		case "tags":
				this.SetStringArr(&this.Tags, v)
		case "versions":
				this.SetStringArr(&this.Versions, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		case "deletedAt":
				this.SetNumIntN(&this.DeletedAt, v)
		}
		return this
}

func (this *UserRelation) setAttributes(data map[string]interface{}, safe ...bool) {
		for key, v := range data {
				if !safe[0] {
						if this.Excludes(key) {
								continue
						}
						if this.IsEmpty(v) {
								continue
						}
				}
				this.Set(key, v)
		}
}

func (this *UserRelationModel) TableName() string {
		return CollectTable
}

func (this *UserRelationModel) CreateIndex() {
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"userId", "targetUserId", "targetType"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndexKey("status")
}

func (this *UserRelationModel) init() *UserRelationModel {
		this._Self = this
		this.Init()
		return this
}

// 唯一
func (this *UserRelationModel) GetByUnique(m beego.M) *UserRelation {
		var (
				err   error
				data  = NewUserRelation()
				query = bson.M{"userId": "", "targetUserId": "", "targetType": ""}
		)
		for key := range query {
				v, ok := m[key]
				if !ok {
						return nil
				}
				str, ok := v.(string)
				if ok && str != "" {
						query[key] = str
						continue
				}
				return nil
		}
		err = this.NewQuery(query).One(data)
		if err == nil {
				return data
		}
		return nil
}

// 保存关系记录
func (this *UserRelationModel) SaveInfo(userId string, targetUserId string, extras ...beego.M) error {
		if userId == "" || targetUserId == "" {
				return errors.New("userId or targetUserId empty")
		}
		if userId == targetUserId {
				return errors.New("userId must not eq targetUserId")
		}
		if len(extras) == 0 {
				extras = append(extras, beego.M{"status": 1})
		}
		var (
				user   = NewUserRelation()
				status = getStatus(extras[0])
				typ    = getTargetType(extras[0])
				query  = bson.M{"userId": userId, "targetUserId": targetUserId, "targetType": typ}
				friend = this.GetByUnique(beego.M(query))
		)
		delete(extras[0],"status")
		delete(extras[0],"targetType")
		if friend == nil {
				user.UserId = userId
				user.TargetUserId = targetUserId
				user.Status = status
				user.Extras = beego.M{}
				user.TargetType = typ
				if len(extras) > 0 {
						user.Extras = Merger(user.Extras, extras[0])
				}
				return user.Save()
		}
		if status != StatusOk {
				return errors.New("error status")
		}
		friend.Status = 1
		friend.Extras = beego.M{}
		friend.TargetType = typ
		if len(extras) > 0 {
				friend.Extras = Merger(friend.Extras, extras[0])
		}
		friend.UpdatedAt = time.Now().Local()
		return this.Update(bson.M{"_id": friend.Id}, friend)
}

// 统计数量
func (this *UserFocusModel) Count(m beego.M) int64 {
		var count, err = this.NewQuery(bson.M(m)).Count()
		if err == nil {
				return int64(count)
		}
		return 0
}

// 状态
func getStatus(m beego.M, defaults ...int) int {
		if len(defaults) == 0 {
				defaults = append(defaults, StatusOk)
		}
		var value, ok = m["status"]
		if !ok {
				return defaults[0]
		}
		if v, ok := value.(int); ok {
				return v
		}
		return defaults[0]
}

// 获取类型
func getTargetType(m beego.M, defaults ...string) string {
		if len(defaults) == 0 {
				defaults = append(defaults, TargetTypeFriend)
		}
		var value, ok = m["targetType"]
		if !ok {
				return defaults[0]
		}
		if v, ok := value.(string); ok {
				return v
		}
		return defaults[0]
}
