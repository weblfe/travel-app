package controllers

import (
	"github.com/astaxie/beego/config/env"
	"github.com/astaxie/beego/context"
	swag "github.com/weblfe/beego-swagger"
	"github.com/weblfe/travel-app/docs"
)

func SwaggerHandlerOf() func(*context.Context) {
	docs.SwaggerInfo.Host = env.Get("APP_URL", docs.SwaggerInfo.Host)
	return swag.Handler
}
