package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

// 用户关注
type UserFocus struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                            // ID
		Status        int           `json:"status" bson:"status"`                     // 状态
		UserId        bson.ObjectId `json:"userId" bson:"userId"`                     // 用户ID
		PostId        bson.ObjectId `json:"postId,omitempty" json:"postId,omitempty"` // 文章ID
		FocusUserId   bson.ObjectId `json:"focusUserId" bson:"focusUserId"`           // 被关注的用户ID
		Extras        beego.M       `json:"extras" bson:"extras"`                     // 扩展数据
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`               // 创建时间
		UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`               // 更新时间
		dataClassImpl `json:",omitempty" bson:",omitempty"`
}

// 用户关注数据模型
type UserFocusModel struct {
		BaseModel
}

func UserFocusModelOf() *UserFocusModel {
		var model = new(UserFocusModel)
		model._Self = model
		model.Init()
		return model
}

func NewUserFocus() *UserFocus {
		var focus = new(UserFocus)
		focus.Init()
		return focus
}

const (
		UserFocusTable = "user_focus"
)

func (this *UserFocus) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

func (this *UserFocus) data() beego.M {
		return beego.M{
				"id":          this.Id.Hex(),
				"status":      this.Status,
				"userId":      this.UserId.Hex(),
				"postId":      this.PostId.Hex(),
				"extras":      this.Extras,
				"focusUserId": this.FocusUserId.Hex(),
				"createdAt":   this.CreatedAt.Unix(),
				"updatedAt":   this.UpdatedAt.Unix(),
		}
}

// 保存
func (this *UserFocus) save() error {
		var (
				model = UserFocusModelOf()
				data  = model.GetByUnique(this.data())
		)
		if data == nil {
				this.InitDefault()
				return model.Add(this)
		}
		return model.Update(bson.M{"_id": data.Id}, this.M(func(m beego.M) beego.M {
				delete(m, "id")
				delete(m, "createdAt")
				m["updatedAt"] = time.Now().Local()
				return m
		}))
}

func (this *UserFocus) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.UpdatedAt.IsZero() {
				this.UpdatedAt = time.Now().Local()
		}
		if this.Extras == nil {
				this.Extras = beego.M{}
		}
		if this.Status == 0 {
				this.Status = 1
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
}

// 设置数值
func (this *UserFocus) setAttributes(data map[string]interface{}, safe ...bool) {
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

func (this *UserFocus) Set(key string, v interface{}) *UserFocus {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "status":
				this.SetNumInt(&this.Status, v)
		case "userId":
				this.SetObjectId(&this.UserId, v)
		case "postId":
				this.SetObjectId(&this.PostId, v)
		case "extras":
				this.SetMapper(&this.Extras, v)
		case "focusUserId":
				this.SetObjectId(&this.FocusUserId, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "updatedAt":
				this.SetTime(&this.UpdatedAt, v)
		}
		return this
}

// 表
func (this *UserFocusModel) TableName() string {
		return UserFocusTable
}

// 创建索引
func (this *UserFocusModel) CreateIndex() {
		_ = this.Collection().EnsureIndex(mgo.Index{
				Key:    []string{"userId", "focusUserId"},
				Unique: true,
				Sparse: false,
		})
		_ = this.Collection().EnsureIndexKey("status")
		_ = this.Collection().EnsureIndexKey("postId")
}

func (this *UserFocusModel) GetByUnique(m beego.M) *UserFocus {
		var (
				err   error
				query = bson.M{
						"userId":      this.id(m["userId"]),
						"focusUserId": this.id(m["focusUserId"]),
				}
				data = NewUserFocus()
		)
		err = this.NewQuery(query).One(data)
		if err == nil {
				return data
		}
		return nil
}

// 是否用户互关注
func (this *UserFocusModel) GetFocusTwo(userId, userId2 string) bool {
		var query = bson.M{
				"$or": []beego.M{
						{"userId": this.id(userId), "focusUserId": this.id(userId2), "status": 1},
						{"userId": this.id(userId2), "focusUserId": this.id(userId), "status": 1},
				},
		}
		var n, err = this.NewQuery(query).Count()
		if err != nil {
				return false
		}
		return n == 2
}

// 获取 用户关注列表
func (this *UserFocusModel) GetUserFocusLists(userId string, params ...ListsParams) ([]*UserFocus, ListsParams) {
		if len(params) == 0 {
				params = append(params, NewListParam(1, 10))
		}
		var (
				err   error
				page  = params[0]
				items = make([]*UserFocus, page.Count())
				query = bson.M{"userId": bson.ObjectIdHex(userId), "status": 1}
		)
		items = items[:0]
		err = this.NewQuery(query).Limit(page.Count()).Skip(page.Skip()).All(&items)
		if err == nil {
				page.SetTotal(this.GetUserFocusCount(userId))
				return items, page
		}
		page.SetTotal(0)
		return nil, page
}

// 获取用户关注数
func (this *UserFocusModel) GetUserFocusCount(userId string) int {
		var (
				query  = bson.M{"userId": bson.ObjectIdHex(userId), "status": 1}
				n, err = this.NewQuery(query).Count()
		)
		if err != nil {
				return 0
		}
		return n
}

// 获取用户被关注数
func (this *UserFocusModel) GetFocusCount(userId string) int {
		var (
				query  = bson.M{"focusUserId": bson.ObjectIdHex(userId), "status": 1}
				n, err = this.NewQuery(query).Count()
		)
		if err != nil {
				return 0
		}
		return n
}

// ID
func (this *UserFocusModel) id(v interface{}) bson.ObjectId {
		if v == nil || v == "" {
				return ""
		}
		if str, ok := v.(string); ok {
				return bson.ObjectIdHex(str)
		}
		if id, ok := v.(bson.ObjectId); ok {
				return id
		}
		return ""
}