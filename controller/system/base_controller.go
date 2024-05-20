package system

import (
	"CatMi-devops/request"
	"CatMi-devops/service/system_svc"
	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// Dashboard 系统首页展示数据
func (m *BaseController) Dashboard(c *gin.Context) {
	req := new(request.BaseDashboardReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Base.Dashboard(c, req)
	})
}

// GetPasswd 生成加密密码
func (m *BaseController) GetPasswd(c *gin.Context) {
	req := new(request.GetPasswdReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.Base.GetPasswd(c, req)
	})
}
