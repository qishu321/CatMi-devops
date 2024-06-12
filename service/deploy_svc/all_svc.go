package deploysvc

import (
	"CatMi-devops/utils/tools"
	"fmt"
)

var (
	ReqdeployErr   = tools.NewRspError(tools.SystemErr, fmt.Errorf("请求异常"))
	SshdSvcs       = &SshdSvc{}
	CommandLogSvcs = &CommandLogSvc{}
	TemplateSvcs   = &TemplateSvc{}
	TaskSvcs       = &TaskSvc{}
	TaskEnvSvcs    = &TaskEnvSvc{}
)
