package system_svc

import (
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

type OperationLogSvc struct{}

// List 数据列表
func (l OperationLogSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.OperationLogListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取数据列表
	logs, err := system.OperationLog.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口列表失败: %s", err.Error()))
	}

	rets := make([]model.OperationLog, 0)
	for _, log := range logs {
		rets = append(rets, *log)
	}
	count, err := system.OperationLog.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口总数失败"))
	}

	return response.LogListRsp{
		Total: count,
		Logs:  rets,
	}, nil

	// 获取
	// logs, err := isql.OperationLog.List(&r)
	// if err != nil {
	// 	response.Fail(c, nil, "获取操作日志列表失败: "+err.Error())
	// 	return
	// }
	// return nil, nil
}

// Delete 删除数据
func (l OperationLogSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.OperationLogDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.OperationLogIds {
		filter := tools.H{"id": int(id)}
		if !system.OperationLog.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("该条记录不存在"))
		}
	}
	// 删除接口
	err := system.OperationLog.Delete(r.OperationLogIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除该条记录失败: %s", err.Error()))
	}
	return nil, nil
}
