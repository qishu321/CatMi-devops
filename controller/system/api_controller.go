package system

import (
	"CatMi-devops/request"
	"CatMi-devops/service/system_svc"
	"github.com/gin-gonic/gin"
)

type ApiController struct{}

// List 记录列表
func (m *ApiController) List(c *gin.Context) {
	req := new(request.ApiListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Api.List(c, req)
	})
}

// GetTree 接口树
func (m *ApiController) GetTree(c *gin.Context) {
	req := new(request.ApiGetTreeReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Api.GetTree(c, req)
	})
}

// Add 新建记录
func (m *ApiController) Add(c *gin.Context) {
	req := new(request.ApiAddReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Api.Add(c, req)
	})
}

// Update 更新记录
func (m *ApiController) Update(c *gin.Context) {
	req := new(request.ApiUpdateReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Api.Update(c, req)
	})
}

// Delete 删除记录
func (m *ApiController) Delete(c *gin.Context) {
	req := new(request.ApiDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Api.Delete(c, req)
	})
}