package cmdbsvc

import (
	"CatMi-devops/initalize/cmdbinit"
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CmdbGroupSvc struct{}

func (s CmdbGroupSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.DeleteReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	err := cmdbinit.ServerGroups.Delete(r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除服务器组失败: %s", err.Error()))
	}

	return nil, nil

}
func (s CmdbGroupSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.ServerGroupReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var err error

	code := cmdbinit.ServerGroups.Check(r.GroupName)
	if code != 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("服务器组名称重复，请换新服务器组名称"))
	}
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	ServerGroup := model.ServerGroup{
		GroupName: r.GroupName,
		Desc:      r.Desc,
		Creator:   ctxUser.Username,
	}
	err = cmdbinit.ServerGroups.Add(&ServerGroup, r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建服务器组失败: %s", err.Error()))
	}

	return ServerGroup, nil
}

func (s CmdbGroupSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.UpdateServerGroupReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var err error

	ServerGroup := model.ServerGroup{
		Model:     gorm.Model{ID: r.ID},
		GroupName: r.GroupName,
		Desc:      r.Desc,
	}
	err = cmdbinit.ServerGroups.Update(&ServerGroup, r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新服务器组失败: %s", err.Error()))
	}

	return nil, nil
}

func (s CmdbGroupSvc) DelGroupId(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.DelGroupIdServerGroupReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var err error

	ServerGroup := model.ServerGroup{
		Model: gorm.Model{ID: r.ID},
	}
	err = cmdbinit.ServerGroups.DelGroupId(&ServerGroup, r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新服务器组失败: %s", err.Error()))
	}

	return nil, nil
}

func (s CmdbGroupSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ListReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.ServerGroups.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器组列表失败: %s", err.Error()))
	}

	lists := make([]model.ServerGroup, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}
	count, err := cmdbinit.ServerGroups.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器组总数失败"))
	}
	return response.ServerGroupReq{
		Total:  count,
		Groups: lists,
	}, nil

}
func (s CmdbGroupSvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.InfoReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.ServerGroups.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器组失败: %s", err.Error()))
	}
	return list, nil

}
