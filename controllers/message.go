package controllers

type MessageController struct {
		BaseController
}

// MessageControllerOf 消息模块 controller
func MessageControllerOf() *MessageController {
		return new(MessageController)
}

// GetApplyFriendsMsgList 获取申请添加用户消息列表接口
// @router /message/apply/friends/list [get]
func (this *MessageController)GetApplyFriendsMsgList()  {

}

// GetMessageList 获取用户消息列表接口
// @router /message/list  [get]
func (this *MessageController)GetMessageList()  {

}

// RemoveMessageById
// @router /message/:id  [delete]
func (this *MessageController)RemoveMessageById()  {

}

// ApplyAddFriend 申请添加好友 接口
// @router /apply/friends [post]
func (this *MessageController)ApplyAddFriend()  {

}

// DealApplyById 通过|决绝 好友申请 接口
// @router /apply/:id   [post]
func (this *MessageController)DealApplyById()  {

}