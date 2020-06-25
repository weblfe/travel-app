package services

type AppService interface {
		GetAppVersion() string
		GetAboutUs() string
		GetPrivacy() string
		GetUserAgreement() string
		GetAppCustomers() []string
		GetAppInfos() map[string]string
}

type AppServiceImpl struct {

}



func AppServiceOf() AppService  {

}