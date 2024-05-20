package system_svc

import (
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
)

type ApiSvc struct{}

// Add 添加数据
func (l ApiSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ApiAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	api := model.Api{
		Method:   r.Method,
		Path:     r.Path,
		Category: r.Category,
		Remark:   r.Remark,
		Creator:  ctxUser.Username,
	}

	// 创建接口
	err = system.Api.Add(&api)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建接口失败: %s", err.Error()))
	}

	return nil, nil
}

// List 数据列表
func (l ApiSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ApiListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取数据列表
	apis, err := system.Api.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口列表失败: %s", err.Error()))
	}

	rets := make([]model.Api, 0)
	for _, api := range apis {
		rets = append(rets, *api)
	}
	count, err := system.Api.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口总数失败"))
	}

	return response.ApiListRsp{
		Total: count,
		Apis:  rets,
	}, nil
}

// GetTree 数据树
func (l ApiSvc) GetTree(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ApiGetTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	_ = r

	apis, err := system.Api.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取资源列表失败: " + err.Error()))
	}

	// 获取所有的分类
	var categoryList []string
	for _, api := range apis {
		categoryList = append(categoryList, api.Category)
	}
	// 获取去重后的分类
	categoryUniq := funk.UniqString(categoryList)

	apiTree := make([]*response.ApiTreeRsp, len(categoryUniq))

	for i, category := range categoryUniq {
		apiTree[i] = &response.ApiTreeRsp{
			ID:       -i,
			Remark:   category,
			Category: category,
			Children: nil,
		}
		for _, api := range apis {
			if category == api.Category {
				apiTree[i].Children = append(apiTree[i].Children, api)
			}
		}
	}

	return apiTree, nil
}

// Update 更新数据
func (l ApiSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ApiUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !system.Api.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("接口不存在"))
	}

	// 获取当前登陆用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户失败"))
	}

	oldData := new(model.Api)
	err = system.Api.Find(filter, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	api := model.Api{
		Model:    oldData.Model,
		Method:   r.Method,
		Path:     r.Path,
		Category: r.Category,
		Remark:   r.Remark,
		Creator:  ctxUser.Username,
	}
	err = system.Api.Update(&api)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新接口失败: %s", err.Error()))
	}
	return nil, nil
}

// Delete 删除数据
func (l ApiSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ApiDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.ApiIds {
		filter := tools.H{"id": int(id)}
		if !system.Api.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("接口不存在"))
		}
	}
	// 删除接口
	err := system.Api.Delete(r.ApiIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除接口失败: %s", err.Error()))
	}
	return nil, nil
}
