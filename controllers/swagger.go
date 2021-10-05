package controllers

import (
	"github.com/astaxie/beego/config/env"
	"github.com/astaxie/beego/context"
	swag "github.com/weblfe/beego-swagger"
	"github.com/weblfe/travel-app/docs"
	"github.com/weblfe/travel-app/libs"
	"sync"
)

func SwaggerHandlerOf() func(*context.Context) {
	var once = sync.Once{}
	return func(c *context.Context) {
		once.Do(func() {
			docs.SwaggerInfo.Host = libs.VariableParse(env.Get("APP_URL", docs.SwaggerInfo.Host))
		})
		swag.Handler(c)
	}
}
