package services

import (
		"github.com/astaxie/beego"
		"github.com/globalsign/mgo"
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
		"regexp"
		"strings"
)

type UserService interface {
		RemoveByUid(uid string) bool
		Create(model *models.User) common.Errors
		Add(user map[string]interface{}) common.Errors
		Inserts(users []map[string]interface{}, txn ...func([]map[string]interface{}) bool) int
		UpdateByUid(uid string, data map[string]interface{}) error
		Lists(page int, size int, args ...interface{}) (items []*models.User, total int, more bool)
		GetByMobile(mobile string) *models.User
		GetByEmail(email string) *models.User
		GetById(id string) *models.User
		GetByUserName(name string) *models.User
}

type UserServiceImpl struct {
		BaseService
		userModel *models.UserModel
}

func UserServiceOf() UserService {
		var service = new(UserServiceImpl)
		service.Init()
		return service
}

func (this *UserServiceImpl) RemoveByUid(uid string) bool {
		if err := this.userModel.Remove(bson.M{"_id": uid}); err == nil {
				return true
		}
		return false
}

func (this *UserServiceImpl) Create(user *models.User) common.Errors {
		if user.AvatarId == "" {
				avatar := AvatarServerOf().GetDefaultAvatar(user.Gender)
				if avatar != nil {
						user.AvatarId = avatar.Id.Hex()
				}
		}
		err := this.userModel.Add(user)
		if err == nil {
				return nil
		}
		return common.NewErrors(err.Error(), common.CreateFailCode)
}

func (this *UserServiceImpl) Add(user map[string]interface{}) common.Errors {
		var userData = new(models.User)
		userData.Load(user)
		if userData.UserName == "" && userData.Mobile == "" && userData.Email == "" {
				return common.NewErrors("miss params", common.EmptyParamCode)
		}
		return this.Create(userData.Defaults())
}

func (this *UserServiceImpl) Inserts(users []map[string]interface{}, txn ...func([]map[string]interface{}) bool) int {
		var count = 0
		if len(txn) == 0 {
				for _, user := range users {
						if err := this.Add(user); err == nil {
								count++
						}
				}
				return count
		}
		if txn[0](users) {
				return len(users)
		}
		return 0
}

func (this *UserServiceImpl) UpdateByUid(uid string, data map[string]interface{}) error {
		var (
				modifies []string
				arr, ok  = data["modifies"]
		)
		if ok && arr != nil {
				if strArr, ok := arr.([]string); ok {
						modifies = strArr
				}
				delete(data, "modifies")
		}
		if len(data) == 0 {
				return common.NewErrors("无更新字段")
		}
		if len(modifies) == 0 {
				user := this.GetById(uid)
				if user == nil {
						return common.NewErrors("用户不存在")
				}
				data = user.M(getDiffFilter(data))
		}
		if len(data) == 0 {
				return nil
		}
		// 无需ID
		err := this.userModel.Update(bson.M{"_id": bson.ObjectIdHex(uid)}, data)
		if info, ok := err.(*mgo.LastError); ok {
				if  this.userModel.IsDuplicate(err) {
						return common.NewErrors(info.Code, strings.Join(getDupKeys(info), " ")+" 已存在!")
				}
				return common.NewErrors(info.Code, info.Err)
		}
		if err == nil {
				return nil
		}
		return common.NewErrors(err)
}

func (this *UserServiceImpl) UpdateUserAddressById(id string, addr *models.UserAddress) error {
		var (
				userId    = bson.ObjectIdHex(id)
				addrModel = models.UserAddressModelOf()
				userAddr  = addrModel.GetAddressByUserId(userId, addr.Type)
		)
		if userAddr == nil {
				addr.UserId = userId
				addr.Type = models.AddressTypeRegister
				addr.InitDefault()
				err := addr.Save()
				if err == nil {
						return this.UpdateByUid(id, beego.M{"addressId": addr.Id})
				}
				return err
		}
		userAddr.SetAttributes(addr.M(), true)
		return userAddr.Update()
}

func (this *UserServiceImpl) Lists(page int, size int, query ...interface{}) (items []*models.User, total int, more bool) {
		if len(query) == 0 {
				query = append(query, nil)
		}
		if len(query) < 2 {
				query = append(query, nil)
		}
		limit := models.NewListParam(page, size)
		total, err := this.userModel.Lists(query[0], items, limit, query[1])
		if err != nil {
				return nil, 0, false
		}
		limit.SetTotal(total)
		return items, total, limit.More()
}

func (this *UserServiceImpl) GetByMobile(mobile string) *models.User {
		var user = new(models.User)
		if err := this.userModel.GetByKey("mobile", mobile, user); err != nil {
				return nil
		}
		return user
}

func (this *UserServiceImpl) GetByEmail(email string) *models.User {
		var user = new(models.User)
		if err := this.userModel.GetByKey("email", email, user); err != nil {
				return nil
		}
		return user
}

func (this *UserServiceImpl) GetById(id string) *models.User {
		var user = new(models.User)
		if err := this.userModel.GetById(id, user); err != nil {
				return nil
		}
		return user
}

func (this *UserServiceImpl) GetByUserName(name string) *models.User {
		var user = new(models.User)
		if err := this.userModel.GetByKey("username", name, user); err != nil {
				return nil
		}
		return user
}

func (this *UserServiceImpl) Init() {
		this.init()
		this.userModel = models.UserModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return UserServiceOf()
		}
}

// 字段对比过滤器
func getDiffFilter(b beego.M) func(m beego.M) beego.M {
		return func(m beego.M) beego.M {
				var result = make(beego.M)
				for key, v := range b {
						diff := m[key]
						if v != nil && diff != v {
								result[key] = v
						}
				}
				return result
		}
}

// 字段更新过滤器
func getUpdateFilter(fields []string, data beego.M) func(m beego.M) beego.M {
		return func(m beego.M) beego.M {
				var result = make(beego.M)
				for _, key := range fields {
						v, ok := data[key]
						if ok && v != nil {
								result[key] = v
						}
				}
				return result
		}
}

// 获取重复键
func getDupKeys(err *mgo.LastError) []string {
		var (
				keys    []string
				reg     = regexp.MustCompile(`.+ dup key: (.+)`)
				regs    = regexp.MustCompile(`.+ dup keys: (.+)`)
				keysReg = regexp.MustCompile(`.+ (\w+:).+`)
		)
		arr := reg.FindAllStringSubmatch(err.Err, -1)
		if len(arr) == 0 {
				arr = regs.FindAllStringSubmatch(err.Err, -1)
		}
		str := arr[0][1]
		if str != "" {
				arr := keysReg.FindAllStringSubmatch(str, -1)
				if len(arr) < 0 {
						return keys
				}
				for _, k := range arr[0][1:] {
						keys = append(keys, k[0:len(k)-1])
				}
		}
		return keys
}
