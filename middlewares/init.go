package middlewares

// 注册 中间键命名空间
func init()  {
		GetTokenMiddleware()
		GetAuthMiddleware()
		GetAttachTicketMiddleware()
		GetRoleMiddleware()
		GetCorsMiddleware()
		GetHeaderWares()
}
