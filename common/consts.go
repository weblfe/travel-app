package common

const (
		WarnCode                  = 3000
		SuccessCode               = 0
		ErrorCode                 = 5000
		DevelopCode               = 3111
		UnLoginCode               = 4001
		CreateFailCode            = 1001
		PermissionCode            = 4003
		LimitCode                 = 4100
		EmptyParamCode            = 1002
		InvalidParametersCode     = 4101
		ServiceFailed             = 4041
		ParamVerifyFailed         = 4020
		Error                     = "error"
		Success                   = "success"
		UnLoginError              = "please login!"
		NotFoundError             = "not found!"
		PermissionError           = "Permission denied!"
		LimitError                = "access api too frequently!"
		InvalidParametersError    = "invalid parameters"
		CreateFail                = "create record fail!"
		NotFound                  = 4004
		RecordNotFound            = 4040
		AccessForbid              = 4003
		InvalidTokenCode          = 4100
		InvalidTokenError         = "invalid token"
		VerifyNotMatch            = 1003
		PasswordOrAccountNotMatch = "password or account not match"
		UserAccountForbid         = "user account forbid"
		RegisterFail              = 1004
		RegisterFailTip           = "register failed!"
		AppTokenCookie            = "authorization"
		DevelopCodeError          = "developing"
		ServiceFailedError        = "server failed!"
		ParamVerifyFailedError    = "param verify error!"
		RecordNotFoundError       = "empty"
		Page                      = 1  // 默认页
		Count                     = 20 // 默认分页量
)
