package deploy

import (
	"CatMi-devops/request"
	deploysvc "CatMi-devops/service/deploy_svc"

	"github.com/gin-gonic/gin"
)

type TemplateController struct{}

func (m *TemplateController) List(c *gin.Context) {
	req := new(request.SShListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateSvcs.List(c, req)
	})
}

func (m *TemplateController) Info(c *gin.Context) {
	req := new(request.SShInfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateSvcs.Info(c, req)
	})
}

func (m *TemplateController) Delete(c *gin.Context) {
	req := new(request.SShDeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateSvcs.Delete(c, req)
	})
}

func (m *TemplateController) Add(c *gin.Context) {
	req := new(request.TemplateReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateSvcs.Add(c, req)
	})
}

func (m *TemplateController) Update(c *gin.Context) {
	req := new(request.UpdateTemplateReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.TemplateSvcs.Update(c, req)
	})
}
