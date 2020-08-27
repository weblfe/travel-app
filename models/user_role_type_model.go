package models

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/transforms"
		"time"
)

type UserRolesConfigModel struct {
		BaseModel
}

// 用户身份类型
type UserRoleType struct {
		Id            bson.ObjectId `json:"id" bson:"_id"`                              // ID
		Name          string        `json:"name" bson:"name"`                           // 身份名
		Role          int           `json:"role" json:"role"`                           // 身份ID
		State         int           `json:"state" json:"state"`                         // 状态
		Level         int           `json:"level" json:"level"`                         // 身份等级 最高 99999,(1~99999)
		Extras        bson.M        `json:"extras,omitempty" json:"extras,omitempty"`   // 扩展信息
		Comment       string        `json:"comment,omitempty" json:"comment,omitempty"` // 备注
		CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`                 // 创建时间
		dataClassImpl `json:",omitempty" bson:",omitempty"`
}

const (
		UserRootRole          = 7
		UserRoleTypeTableName = "user_roles_config"
)

func NewUserRole() *UserRoleType {
		var data = new(UserRoleType)
		data.Init()
		return data
}

func UserRolesConfigModelOf() *UserRolesConfigModel {
		var model = new(UserRolesConfigModel)
		model.Bind(model)
		model.Init()
		return model
}

func (this *UserRoleType) Init() {
		this.AddFilters(transforms.FilterEmpty)
		this.SetProvider(DataProvider, this.data)
		this.SetProvider(SaverProvider, this.save)
		this.SetProvider(DefaultProvider, this.setDefaults)
		this.SetProvider(AttributesProvider, this.setAttributes)
}

// 数据输出
func (this *UserRoleType) data() bson.M {
		return bson.M{
				"id":        this.Id.Hex(),
				"name":      this.Name,
				"role":      this.Role,
				"state":     this.State,
				"extras":    this.Extras,
				"comment":   this.Comment,
				"level":     this.Level,
				"createdAt": this.CreatedAt,
		}
}

// 属性设置器
func (this *UserRoleType) setAttributes(data map[string]interface{}, safe ...bool) {
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

// 默认处理
func (this *UserRoleType) setDefaults() {
		if this.CreatedAt.IsZero() {
				this.CreatedAt = time.Now().Local()
		}
		if this.State == 0 {
				this.State = 1
		}
		if this.Extras == nil {
				this.Extras = map[string]interface{}{}
		}
		if this.Id == "" {
				this.Id = bson.NewObjectId()
		}
}

// 保存逻辑
func (this *UserRoleType) save() error {
		var (
				userRole = NewUserRole()
				model    = UserRolesConfigModelOf()
		)
		// 检查是更新还是创建
		if this.Id != "" {
				err := model.GetByObjectId(this.Id, userRole)
				if err == nil {
						userRole.SetAttributes(this.M())
						return model.Update(beego.M{"_id": this.Id}, userRole)
				}
		}
		userRole = model.GetByUnique(this.M())
		if userRole != nil {
				userRole.SetAttributes(this.M())
				return model.Update(beego.M{"_id": this.Id}, userRole)
		}
		this.InitDefault()
		return model.Add(this)
}

// setter
func (this *UserRoleType) Set(key string, v interface{}) *UserRoleType {
		switch key {
		case "id":
				this.SetObjectId(&this.Id, v)
		case "role":
				this.SetNumInt(&this.Role, v)
		case "createdAt":
				this.SetTime(&this.CreatedAt, v)
		case "state":
				this.SetNumInt(&this.State, v)
		case "name":
				this.SetString(&this.Name, v)
		case "comment":
				this.SetString(&this.Comment, v)
		case "level":
				this.SetNumInt(&this.Level, v)
				if this.Level > 10000 || this.Level < 1 {
						this.Level = 1
				}
		case "extra":
				this.SetBsonMapper(&this.Extras, v)
		}
		return this
}

func (this *UserRolesConfigModel) TableName() string {
		return UserRoleTypeTableName
}

func (this *UserRolesConfigModel) GetByUnique(data map[string]interface{}) *UserRoleType {
		if len(data) == 0 {
				return nil
		}
		var (
				role   = NewUserRole()
				roleId = data["role"]
		)
		if roleId == nil || roleId == "" {
				return nil
		}
		if err := this.FindOne(bson.M{"role": roleId}, role); err == nil {
				return role
		}
		return nil
}

func (this *UserRolesConfigModel) CreateIndex(force ...bool) {
		this.createIndex(this.getCreateIndexHandler(), force...)
}

func (this *UserRolesConfigModel) getCreateIndexHandler() func(*mgo.Collection) {
		return func(doc *mgo.Collection) {
				// null unique username
				this.logs(doc.EnsureIndex(mgo.Index{
						Key:    []string{"role"},
						Unique: true,
						Sparse: false,
				}))
				this.logs(doc.EnsureIndexKey("name", "state", "level"))
		}
}

func (this *UserRolesConfigModel) GetRoleName(role int) string {
		var roleInfo = NewUserRole()
		if err := this.GetByKey("role", role, roleInfo); err == nil {
				return roleInfo.Name
		}
		return ""
}

// 批量导入
func (this *UserRolesConfigModel) Adds(items []map[string]interface{}) error {
		if len(items) == 0 {
				return ErrEmptyData
		}
		var result []interface{}
		for _, it := range items {
				role := this.GetByUnique(it)
				if role != nil {
						_ = this.Update(bson.M{"_id": role.Id}, it)
				} else {
						var role = NewUserRole()
						role.SetAttributes(it, false)
						role.InitDefault()
						result = append(result, role)
				}
		}
		if len(result) == 0 {
				return nil
		}
		if err := this.Inserts(result); err != nil {
				return err
		}
		return nil
}
