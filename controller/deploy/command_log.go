package deploy

import (
	"CatMi-devops/request"
	deploysvc "CatMi-devops/service/deploy_svc"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

func (m *LogController) List(c *gin.Context) {
	req := new(request.SShListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.CommandLogSvcs.List(c, req)
	})
}
