package deploysvc

import (
	"CatMi-devops/initalize/deployinit"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"

	"github.com/gin-gonic/gin"
)

type TemplateLogSvc struct{}

func (s TemplateLogSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShListReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c

	list, err := deployinit.Template_Logs.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器失败: %s", err.Error()))
	}
	lists := make([]model.Template_Log, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}
	count, err := deployinit.Templates.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取模板总数失败"))
	}
	return response.Template_logRsp{
		Total:         count,
		Templates_log: lists,
	}, nil

}

func (s TemplateLogSvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShInfoReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	list, err := deployinit.Template_Logs.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取脚本失败: %s", err.Error()))
	}
	return list, nil
}
