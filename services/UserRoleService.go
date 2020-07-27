package services

import "github.com/weblfe/travel-app/models"

type UserRoleService interface {
		GetRoleById(string) *models.UserRoleType
		GetRoleDesc(int) string
}

type userRoleService struct {
		BaseService
		model *models.UserRolesConfigModel
}

var (
		_UserRoleServiceIns UserRoleService
)

func UserRoleServiceOf() UserRoleService {
		if _UserRoleServiceIns == nil {
				_UserRoleServiceIns = newUserRoleService()
		}
		return _UserRoleServiceIns
}

func newUserRoleService() *userRoleService {
		var service = new(userRoleService)
		service.Init()
		return service
}

func (this *userRoleService) Init() {
		this.init()
		this.model = models.UserRolesConfigModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return UserRoleServiceOf()
		}
}

func (this *userRoleService) GetRoleById(id string) *models.UserRoleType {
		var (
				data = models.NewUserRole()
				err  = this.model.GetById(id, data)
		)
		if err == nil {
				return data
		}
		return nil
}

func (this *userRoleService) GetRoleDesc(role int) string {
		var (
				data = models.NewUserRole()
				err  = this.model.GetByKey("role", role, data)
		)
		if err == nil {
				return data.Name
		}
		return ""
}
