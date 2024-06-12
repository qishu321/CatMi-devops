package cmdbsvc

import (
	"CatMi-devops/utils/tools"
	"fmt"
)

var (
	ReqCmdbErr = tools.NewRspError(tools.SystemErr, fmt.Errorf("请求异常"))
	Keys       = &KeySvc{}
	CmdbGroup  = &CmdbGroupSvc{}
	Cmdb       = &CmdbSvc{}
)
