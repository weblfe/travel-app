package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/attachment/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"],
        beego.ControllerComments{
            Method: "Ticket",
            Router: `/attachment/ticket`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"],
        beego.ControllerComments{
            Method: "Upload",
            Router: `/attachment/upload`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"],
        beego.ControllerComments{
            Method: "Uploads",
            Router: `/attachment/uploads`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"],
        beego.ControllerComments{
            Method: "GetByMediaId",
            Router: `/attachments/:mediaId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:AttachmentController"],
        beego.ControllerComments{
            Method: "DownloadByMediaId",
            Router: `/attachments/download/:mediaId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:CaptchaController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:CaptchaController"],
        beego.ControllerComments{
            Method: "SendEmailCaptcha",
            Router: `/captcha/email`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:CaptchaController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:CaptchaController"],
        beego.ControllerComments{
            Method: "SendMobileCaptcha",
            Router: `/captcha/mobile`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:CaptchaController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:CaptchaController"],
        beego.ControllerComments{
            Method: "SendWeChatCaptcha",
            Router: `/captcha/wechat`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Index",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Index",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Index",
            Router: `/`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:IndexController"],
        beego.ControllerComments{
            Method: "Index",
            Router: `/`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"],
        beego.ControllerComments{
            Method: "DealApplyById",
            Router: `/apply/:id`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"],
        beego.ControllerComments{
            Method: "ApplyAddFriend",
            Router: `/apply/friends`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"],
        beego.ControllerComments{
            Method: "RemoveMessageById",
            Router: `/message/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"],
        beego.ControllerComments{
            Method: "GetApplyFriendsMsgList",
            Router: `/message/apply/friends/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:MessageController"],
        beego.ControllerComments{
            Method: "GetMessageList",
            Router: `/message/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"],
        beego.ControllerComments{
            Method: "DetailById",
            Router: `/posts/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"],
        beego.ControllerComments{
            Method: "RemoveById",
            Router: `/posts/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"],
        beego.ControllerComments{
            Method: "Create",
            Router: `/posts/create`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:PostsController"],
        beego.ControllerComments{
            Method: "Update",
            Router: `/posts/update`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "FocusOff",
            Router: `/focus/off/:userId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "FocusOn",
            Router: `/focus/on/:userId`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "Register",
            Router: `/register`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "ResetPassword",
            Router: `/reset/password`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetUserFriends",
            Router: `/user/friends`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetUserInfo",
            Router: `/user/info`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/weblfe/travel-app/controllers:UserController"],
        beego.ControllerComments{
            Method: "UpdateUserInfo",
            Router: `/user/info`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}