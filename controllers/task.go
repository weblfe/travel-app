package controllers

import "github.com/weblfe/travel-app/repositories"

type TaskController struct {
		BaseController
}

// TaskControllerOf 任务控制器
func TaskControllerOf() *TaskController {
		return new(TaskController)
}

// Create 创建任务
// @router /task/create [post]
func (this *TaskController) Create() {
		this.Send(nil)
}

// Hook 手动触发任务
// @router /task/hook [post]
func (this *TaskController) Hook() {
		this.Send(nil)
}

// Add 添加任务
// @router /task/add [post]
func (this *TaskController) Add() {
		this.Send(nil)
}

// Stop 停止任务
// @router /task/stop [patch]
func (this *TaskController) Stop() {
		this.Send(nil)
}

// Update 更新任务
// @router /task/update [patch]
func (this *TaskController) Update() {
		this.Send(nil)
}

// Remove 删除任务
// @router /task/del [delete]
func (this *TaskController) Remove() {
		this.Send(nil)
}

// Lists 任务列表
// @router /task/lists [get]
func (this *TaskController) Lists() {
		this.Send(nil)
}

// SyncAssetsToOss 同步数据到oss
// @router /task/sync/assets [post]
func (this *TaskController) SyncAssetsToOss() {
		this.Send(repositories.NewTaskRepository(this).SyncAssetsTask())
}
