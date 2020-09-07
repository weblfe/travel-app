package controllers

type AuthenticateController struct {
		BaseController
}

func AuthenticateControllerOf() *AuthenticateController  {
  var controller = new(AuthenticateController)
  return controller
}

func (this *AuthenticateController)Authenticate() {

}

func (this *AuthenticateController)Auth2()  {

}

func (this *AuthenticateController)AuthLoginByQQ() {

}

func (this *AuthenticateController)AuthLoginByWeChat()  {

}

func (this *AuthenticateController)AuthLoginByAppleId()  {

}