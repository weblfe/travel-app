package test

import (
		"github.com/weblfe/travel-app/libs"
		"github.com/weblfe/travel-app/services"
		"testing"
)

func TestInitDataService(t *testing.T) {
		var (
				path    = libs.VariableParse("${GOPATH}/src/github.com/weblfe/travel-app/static/database")
				service = services.GetInitDataServiceInstance()
		)
		service.SetInit(path)
		service.Init()
}
