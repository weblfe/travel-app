package services

import (
		"github.com/globalsign/mgo/bson"
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
)

type UserService interface {
		RemoveByUid(uid string) bool
		Create(model *models.User) common.Errors
		Add(user map[string]interface{}) common.Errors
		Inserts(users []map[string]interface{}, txn ...func([]map[string]interface{}) bool) int
		UpdateByUid(uid string, data map[string]interface{}) bool
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

func (this *UserServiceImpl) UpdateByUid(uid string, data map[string]interface{}) bool {
		if err := this.userModel.Update(bson.M{"_id": bson.ObjectId(uid)}, data); err == nil {
				return true
		}
		return false
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
