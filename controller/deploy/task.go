package deploy

import (
	"CatMi-devops/request"
	deploysvc "CatMi-devops/service/deploy_svc"

	"github.com/gin-gonic/gin"
)

type TaskController struct{}

func (m *TaskController) List(c *gin.Context) {
	req := new(request.SShListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskSvcs.List(c, req)
	})
}

func (m *TaskController) Info(c *gin.Context) {
	req := new(request.SShInfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskSvcs.Info(c, req)
	})
}

func (m *TaskController) Delete(c *gin.Context) {
	req := new(request.SShDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskSvcs.Delete(c, req)
	})
}

func (m *TaskController) Add(c *gin.Context) {
	req := new(request.TaskReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskSvcs.Add(c, req)
	})
}

func (m *TaskController) Update(c *gin.Context) {
	req := new(request.UpdateTaskReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TaskSvcs.Update(c, req)
	})
}
