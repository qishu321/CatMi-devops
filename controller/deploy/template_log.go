package deploy

import (
	"CatMi-devops/request"
	deploysvc "CatMi-devops/service/deploy_svc"

	"github.com/gin-gonic/gin"
)

type Template_Log_Controller struct{}

func (m *Template_Log_Controller) List(c *gin.Context) {
	req := new(request.SShListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateLogSvcs.List(c, req)
	})
}

func (m *Template_Log_Controller) Info(c *gin.Context) {
	req := new(request.SShInfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateLogSvcs.Info(c, req)
	})
}
