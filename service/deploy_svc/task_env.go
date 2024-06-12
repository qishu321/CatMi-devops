package deploysvc

import (
	"CatMi-devops/initalize/deployinit"
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskEnvSvc struct{}

func (s TaskEnvSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.TaskEnvReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var err error

	code := deployinit.TaskEnvs.Check(r.Name)
	if code != 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("脚本参数名称重复，请换新脚本参数名称"))
	}
	// Create key object
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}
	OptionsJSON, err := json.Marshal(r.Options)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("OptionsJSON 解析失败"))
	}

	Server := model.TaskEnv{
		Name:        r.Name,
		Options:     string(OptionsJSON),
		Important:   r.Important,
		Description: r.Description,
		Creator:     ctxUser.Username,
	}
	err = deployinit.TaskEnvs.Add(&Server)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建脚本参数失败: %s", err.Error()))
	}

	return Server, nil
}

func (s TaskEnvSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.UpdateTaskEnvReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var err error

	if deployinit.TaskEnvs.UpdateCheck(r.Name, r.ID) {
		return nil, tools.NewMySqlError(fmt.Errorf("脚本参数名称重复，请换新脚本参数名称"))
	}
	// Create key object
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}
	OptionsJSON, err := json.Marshal(r.Options)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("OptionsJSON 解析失败"))
	}

	Server := model.TaskEnv{
		Model:       gorm.Model{ID: r.ID},
		Name:        r.Name,
		Options:     string(OptionsJSON),
		Important:   r.Important,
		Description: r.Description,
		Creator:     ctxUser.Username,
	}
	err = deployinit.TaskEnvs.Update(&Server)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新脚本参数失败: %s", err.Error()))
	}

	return Server, nil
}

func (s TaskEnvSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShListReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.TaskEnvs.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本参数列表失败: %s", err.Error()))
	}

	lists := make([]model.TaskEnv, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}
	count, err := deployinit.TaskEnvs.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本参数总数失败"))
	}
	return response.TaskEnvRsp{
		Total:    count,
		TaskEnvs: lists,
	}, nil

}

func (s TaskEnvSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShDeleteReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	err := deployinit.TaskEnvs.Delete(r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除脚本参数失败: %s", err.Error()))
	}

	return nil, nil

}
func (s TaskEnvSvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShInfoReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.TaskEnvs.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本参数失败: %s", err.Error()))
	}
	return list, nil
}
