package cmdb

import (
	"CatMi-devops/request"
	cmdbsvc "CatMi-devops/service/cmdb_svc"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

func (m *GroupController) List(c *gin.Context) {
	req := new(request.ListReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.CmdbGroup.List(c, req)
	})
}
func (m *GroupController) Info(c *gin.Context) {
	req := new(request.InfoReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.CmdbGroup.Info(c, req)
	})
}

func (m *GroupController) Add(c *gin.Context) {
	req := new(request.ServerGroupReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.CmdbGroup.Add(c, req)
	})
}

func (m *GroupController) Update(c *gin.Context) {
	req := new(request.UpdateServerGroupReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.CmdbGroup.Update(c, req)
	})
}

func (m *GroupController) Delete(c *gin.Context) {
	req := new(request.DeleteReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.CmdbGroup.Delete(c, req)
	})
}

func (m *GroupController) DelGroupId(c *gin.Context) {
	req := new(request.DelGroupIdServerGroupReq)
	Run(c, req, func() (interface{}, interface{}) {
		return cmdbsvc.CmdbGroup.DelGroupId(c, req)
	})
}
