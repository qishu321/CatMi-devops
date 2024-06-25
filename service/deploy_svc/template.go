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

type TemplateSvc struct{}

func (s TemplateSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.TemplateReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var err error

	code := deployinit.Templates.Check(r.Name)
	if code != 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("服务器名称重复，请换新服务器名称"))
	}
	// Create key object
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}

	StepnamesJSON, err := json.Marshal(r.Stepnames)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("CmdbNames 解析失败"))
	}
	// TaskenvJSON, err := json.Marshal(r.Taskenv)
	// if err != nil {
	// 	return nil, tools.NewMySqlError(fmt.Errorf("TaskenvJSON 解析失败"))
	// }

	Server := model.Template{
		TemplateID: tools.GenerateRandomNumber(),
		Taskenv:    r.Taskenv,
		Important:  r.Important,
		Name:       r.Name,
		Stepnames:  string(StepnamesJSON),
		Creator:    ctxUser.Username,
	}
	err = deployinit.Templates.Add(&Server)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建服务器失败: %s", err.Error()))
	}

	return Server, nil
}

func (s TemplateSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.UpdateTemplateReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var err error

	if deployinit.Templates.UpdateCheck(r.Name, r.ID) {
		return nil, tools.NewMySqlError(fmt.Errorf("模板名称重复，请换新模板名称"))
	}
	// Create key object
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}

	StepnamesJSON, err := json.Marshal(r.Stepnames)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("CmdbNames 解析失败"))
	}
	// TaskenvJSON, err := json.Marshal(r.Taskenv)
	// if err != nil {
	// 	return nil, tools.NewMySqlError(fmt.Errorf("TaskenvJSON 解析失败"))
	// }

	Server := model.Template{
		Model:     gorm.Model{ID: r.ID},
		Important: r.Important,
		Taskenv:   r.Taskenv,
		Name:      r.Name,
		Stepnames: string(StepnamesJSON),
		Creator:   ctxUser.Username,
	}
	err = deployinit.Templates.Update(&Server)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新模板失败: %s", err.Error()))
	}

	return Server, nil
}

func (s TemplateSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShListReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.Templates.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取模板列表失败: %s", err.Error()))
	}

	lists := make([]model.Template, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}
	count, err := deployinit.Templates.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取模板总数失败"))
	}
	return response.TemplateRsp{
		Total:     count,
		Templates: lists,
	}, nil

}

func (s TemplateSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShDeleteReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	err := deployinit.Templates.Delete(r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除模板失败: %s", err.Error()))
	}

	return nil, nil

}
func (s TemplateSvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShInfoReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.Templates.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取模板失败: %s", err.Error()))
	}
	return list, nil

}
