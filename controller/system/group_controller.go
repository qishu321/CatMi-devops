package system

import (
	"CatMi-devops/request"
	"CatMi-devops/service/system_svc"
	"github.com/gin-gonic/gin"
)

type GroupController struct{}

// List 记录列表
func (m *GroupController) List(c *gin.Context) {
	req := new(request.GroupListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.List(c, req)
	})
}

// UserInGroup 在分组内的用户
func (m *GroupController) UserInGroup(c *gin.Context) {
	req := new(request.UserInGroupReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.UserInGroup(c, req)
	})
}

// UserNoInGroup 不在分组的用户
func (m *GroupController) UserNoInGroup(c *gin.Context) {
	req := new(request.UserNoInGroupReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.UserNoInGroup(c, req)
	})
}

// GetTree 接口树
func (m *GroupController) GetTree(c *gin.Context) {
	req := new(request.GroupListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.GetTree(c, req)
	})
}

// Add 新建记录
func (m *GroupController) Add(c *gin.Context) {
	req := new(request.GroupAddReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.Add(c, req)
	})
}

// Update 更新记录
func (m *GroupController) Update(c *gin.Context) {
	req := new(request.GroupUpdateReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.Update(c, req)
	})
}

// Delete 删除记录
func (m *GroupController) Delete(c *gin.Context) {
	req := new(request.GroupDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.Delete(c, req)
	})
}

// AddUser 添加用户
func (m *GroupController) AddUser(c *gin.Context) {
	req := new(request.GroupAddUserReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.AddUser(c, req)
	})
}

// RemoveUser 移除用户
func (m *GroupController) RemoveUser(c *gin.Context) {
	req := new(request.GroupRemoveUserReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Group.RemoveUser(c, req)
	})
}
