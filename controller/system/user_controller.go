package system

import (
	"CatMi-devops/request"
	"CatMi-devops/service/system_svc"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// Add 添加记录
func (m *UserController) Add(c *gin.Context) {
	req := new(request.UserAddReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.Add(c, req)
	})
}

// Update 更新记录
func (m *UserController) Update(c *gin.Context) {
	req := new(request.UserUpdateReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.Update(c, req)
	})
}

// List 记录列表
func (m *UserController) List(c *gin.Context) {
	req := new(request.UserListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.List(c, req)
	})
}

// Delete 删除记录
func (m UserController) Delete(c *gin.Context) {
	req := new(request.UserDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.Delete(c, req)
	})
}

// ChangePwd 更新密码
func (m UserController) ChangePwd(c *gin.Context) {
	req := new(request.UserChangePwdReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.ChangePwd(c, req)
	})
}

// ChangeUserStatus 更改用户状态
func (m UserController) ChangeUserStatus(c *gin.Context) {
	req := new(request.UserChangeUserStatusReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.ChangeUserStatus(c, req)
	})
}

// GetUserInfo 获取当前登录用户信息
func (uc UserController) GetUserInfo(c *gin.Context) {
	req := new(request.UserGetUserInfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.User.GetUserInfo(c, req)
	})
}
