package cmdb

import (
	"CatMi-devops/request"
	cmdbsvc "CatMi-devops/service/cmdb_svc"

	"github.com/gin-gonic/gin"
)

type ServerController struct{}

func (m *ServerController) List(c *gin.Context) {
	req := new(request.ListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.List(c, req)
	})
}
func (m *ServerController) Info(c *gin.Context) {
	req := new(request.InfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.Info(c, req)
	})
}
func (m *ServerController) EnableList(c *gin.Context) {
	req := new(request.ListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.EnableList(c, req)
	})
}

func (m *ServerController) Add(c *gin.Context) {
	req := new(request.ServerCmdbReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.Add(c, req)
	})
}

func (m *ServerController) Update(c *gin.Context) {
	req := new(request.UpdateServerCmdbReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.Update(c, req)
	})
}

func (m *ServerController) Delete(c *gin.Context) {
	req := new(request.DeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.Delete(c, req)
	})
}

func (m *ServerController) Enabled(c *gin.Context) {
	req := new(request.EnabledServerCmdbReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.Enabled(c, req)
	})
}

func (m *ServerController) EnabledParams(c *gin.Context) {
	req := new(request.EnabledCmdbReqParams)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.Cmdb.EnabledParams(c, req)
	})
}
