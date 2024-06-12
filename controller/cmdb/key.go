package cmdb

import (
	"CatMi-devops/request"
	cmdbsvc "CatMi-devops/service/cmdb_svc"

	"github.com/gin-gonic/gin"
)

type KeyController struct{}

func (m *KeyController) List(c *gin.Context) {
	req := new(request.ListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Keys.List(c, req)
	})
}
func (m *KeyController) Info(c *gin.Context) {
	req := new(request.InfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Keys.Info(c, req)
	})
}

func (m *KeyController) Add(c *gin.Context) {
	req := new(request.KeyReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Keys.Add(c, req)
	})
}

func (m *KeyController) Update(c *gin.Context) {
	req := new(request.KeyReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Keys.Update(c, req)
	})
}

func (m *KeyController) Delete(c *gin.Context) {
	req := new(request.DeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Keys.Delete(c, req)
	})
}
