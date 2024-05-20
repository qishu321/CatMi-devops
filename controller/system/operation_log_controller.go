package system

import (
	"CatMi-devops/request"
	"CatMi-devops/service/system_svc"
	"github.com/gin-gonic/gin"
)

type OperationLogController struct{}

// List 记录列表
func (m *OperationLogController) List(c *gin.Context) {
	req := new(request.OperationLogListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.OperationLog.List(c, req)
	})
}

// Delete 删除记录
func (m *OperationLogController) Delete(c *gin.Context) {
	req := new(request.OperationLogDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return system_svc.OperationLog.Delete(c, req)
	})
}
