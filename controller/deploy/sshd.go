package deploy

import (
	"CatMi-devops/request"
	deploysvc "CatMi-devops/service/deploy_svc"

	"github.com/gin-gonic/gin"
)

type SshdController struct{}

func (m *SshdController) Command(c *gin.Context) {
	req := new(request.SSHClientConfigReq)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.SshdSvcs.Command(c, req)
	})
}

func (m *SshdController) SshdCommandParams(c *gin.Context) {
	req := new(request.SSHReqParams)
	Run(c, req, func() (interface{}, interface{}) {
		return deploysvc.SshdSvcs.SshdCommandParams(c, req)
	})
}
