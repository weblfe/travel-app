package middlewares

// 注册中间键
func init()  {
		GetTokenMiddleware()
		GetAuthMiddleware()
		GetAttachTicketMiddleware()
		GetRoleMiddleware()
}
