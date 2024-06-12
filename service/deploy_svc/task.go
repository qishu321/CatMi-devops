package deploysvc

import (
	"CatMi-devops/initalize/deployinit"
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskSvc struct{}

func (s TaskSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.TaskReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var err error

	code := deployinit.Tasks.Check(r.Name)
	if code != 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("脚本名称重复，请换新脚本名称"))
	}
	// Create key object
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}
	script := r.Script
	for key, value := range r.Variables {
		script = strings.Replace(script, "{{"+key+"}}", value, -1)
	}
	Server := model.Task{
		Name:        r.Name,
		Type:        r.Type,
		Description: r.Description,
		Script:      script,
		Creator:     ctxUser.Username,
	}
	err = deployinit.Tasks.Add(&Server)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建脚本失败: %s", err.Error()))
	}

	return Server, nil
}

func (s TaskSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.UpdateTaskReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var err error

	if deployinit.Tasks.UpdateCheck(r.Name, r.ID) {
		return nil, tools.NewMySqlError(fmt.Errorf("脚本名称重复，请换新脚本名称"))
	}
	// Create key object
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}

	Server := model.Task{
		Model:       gorm.Model{ID: r.ID},
		Name:        r.Name,
		Type:        r.Type,
		Description: r.Description,
		Script:      r.Script,
		Creator:     ctxUser.Username,
	}
	err = deployinit.Tasks.Update(&Server)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新脚本失败: %s", err.Error()))
	}

	return Server, nil
}

func (s TaskSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShListReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.Tasks.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本列表失败: %s", err.Error()))
	}

	lists := make([]model.Task, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}
	count, err := deployinit.Tasks.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本总数失败"))
	}
	return response.TaskRsp{
		Total: count,
		Tasks: lists,
	}, nil

}

func (s TaskSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShDeleteReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	err := deployinit.Tasks.Delete(r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除脚本失败: %s", err.Error()))
	}

	return nil, nil

}
func (s TaskSvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShInfoReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.Tasks.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本失败: %s", err.Error()))
	}
	return list, nil
}
