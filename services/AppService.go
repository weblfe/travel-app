package services

import "github.com/weblfe/travel-app/models"

type AppService interface {
		GetAppVersion() string
		GetAboutUs() string
		GetPrivacy() string
		GetUserAgreement() string
		GetAppCustomers() []string
		GetAppInfos() map[string]string
}

type AppServiceImpl struct {
		BaseService
		appModel *models.AppModel
}

func AppServiceOf() AppService {
		var service = new(AppServiceImpl)
		service.Init()
		return service
}

func (this *AppServiceImpl) GetAppVersion() string {
		panic("implement me")
}

func (this *AppServiceImpl) GetAboutUs() string {
		panic("implement me")
}

func (this *AppServiceImpl) GetPrivacy() string {
		panic("implement me")
}

func (this *AppServiceImpl) GetUserAgreement() string {
		panic("implement me")
}

func (this *AppServiceImpl) GetAppCustomers() []string {
		panic("implement me")
}

func (this *AppServiceImpl) GetAppInfos() map[string]string {
		panic("implement me")
}

func (this *AppServiceImpl) Init() {
		this.appModel = models.AppModelOf()
		this.Constructor = func(args ...interface{}) interface{} {
				return AppServiceOf()
		}

}
