package services

import (
		"github.com/weblfe/travel-app/common"
		"github.com/weblfe/travel-app/models"
)

type UserService interface {
		RemoveByUid(uid string) bool
		Create(model *models.UserModel) common.Errors
		Add(user map[string]interface{}) common.Errors
		Inserts(users []map[string]interface{}, txn ...func([]map[string]interface{}) bool) int
		UpdateByUid(uid string, data map[string]interface{}) bool
		Lists(page int, size int, args ...interface{}) (items []*models.UserModel, total int, more bool)
		GetByMobile(mobile string) *models.UserModel
		GetByEmail(email string) *models.UserModel
		GetById(id string) *models.UserModel
		GetByUserName(name string) *models.UserModel
}

type UserServiceImpl struct {
		BaseService
	 userModel	*models.UserModel
}

func UserServiceOf() UserService {
		var service = new(UserServiceImpl)
		service.Init()
		return service
}

func (this *UserServiceImpl) RemoveByUid(uid string) bool {
		panic("implement me")
}

func (this *UserServiceImpl) Create(model *models.UserModel) common.Errors {
		panic("implement me")
}

func (this *UserServiceImpl) Add(user map[string]interface{}) common.Errors {
		panic("implement me")
}

func (this *UserServiceImpl) Inserts(users []map[string]interface{}, txn ...func([]map[string]interface{}) bool) int {
		panic("implement me")
}

func (this *UserServiceImpl) UpdateByUid(uid string, data map[string]interface{}) bool {
		panic("implement me")
}

func (this *UserServiceImpl) Lists(page int, size int, args ...interface{}) (items []*models.UserModel, total int, more bool) {
		panic("implement me")
}

func (this *UserServiceImpl) GetByMobile(mobile string) *models.UserModel {
		panic("implement me")
}

func (this *UserServiceImpl) GetByEmail(email string) *models.UserModel {
		panic("implement me")
}

func (this *UserServiceImpl) GetById(id string) *models.UserModel {
		panic("implement me")
}

func (this *UserServiceImpl) GetByUserName(name string) *models.UserModel {
		panic("implement me")
}

func (this *UserServiceImpl) Init() {
	 this.userModel = models.UserModelOf()
	 this.Constructor = func(args ...interface{}) interface{} {
			 return UserServiceOf()
	 }
}


