package deploy

import (
	"CatMi-devops/request"
	deploysvc "CatMi-devops/service/deploy_svc"

	"github.com/gin-gonic/gin"
)

type TaskEnvController struct{}

func (m *TaskEnvController) List(c *gin.Context) {
	req := new(request.SShListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskEnvSvcs.List(c, req)
	})
}

func (m *TaskEnvController) Info(c *gin.Context) {
	req := new(request.SShInfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskEnvSvcs.Info(c, req)
	})
}

func (m *TaskEnvController) Delete(c *gin.Context) {
	req := new(request.SShDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskEnvSvcs.Delete(c, req)
	})
}

func (m *TaskEnvController) Add(c *gin.Context) {
	req := new(request.TaskEnvReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskEnvSvcs.Add(c, req)
	})
}

func (m *TaskEnvController) Update(c *gin.Context) {
	req := new(request.UpdateTaskEnvReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskEnvSvcs.Update(c, req)
	})
}
