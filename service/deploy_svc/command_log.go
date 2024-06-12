package deploysvc

import (
	"CatMi-devops/initalize/deployinit"
	"CatMi-devops/request"
	"CatMi-devops/utils/tools"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CommandLogSvc struct{}

func (s CommandLogSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.SShListReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c

	list, err := deployinit.CommandLogs.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器失败: %s", err.Error()))
	}
	return list, nil
}
