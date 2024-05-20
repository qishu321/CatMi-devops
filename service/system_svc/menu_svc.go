package system_svc

import (
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

type MenuSvc struct{}

// Add 添加数据
func (l MenuSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.MenuAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if system.Menu.Exist(tools.H{"name": r.Name}) {
		return nil, tools.NewMySqlError(fmt.Errorf("菜单名称已存在"))

	}

	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	menu := model.Menu{
		Name:       r.Name,
		Title:      r.Title,
		Icon:       r.Icon,
		Path:       r.Path,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Sort:       r.Sort,
		Status:     r.Status,
		Hidden:     r.Hidden,
		NoCache:    r.NoCache,
		AlwaysShow: r.AlwaysShow,
		Breadcrumb: r.Breadcrumb,
		ActiveMenu: r.ActiveMenu,
		ParentId:   r.ParentId,
		Creator:    ctxUser.Username,
	}

	err = system.Menu.Add(&menu)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建记录失败: %s", err.Error()))
	}

	return nil, nil
}

// // List 数据列表
// func (l MenuSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
// 	_, ok := req.(*request.MenuListReq)
// 	if !ok {
// 		return nil, ReqAssertErr
// 	}
// 	_ = c

// 	menus, err := system.Menu.List()
// 	if err != nil {
// 		return nil, tools.NewMySqlError(fmt.Errorf("获取资源列表失败: %s", err.Error()))
// 	}

// 	rets := make([]model.Menu, 0)
// 	for _, menu := range menus {
// 		rets = append(rets, *menu)
// 	}
// 	count, err := system.Menu.Count()
// 	if err != nil {
// 		return nil, tools.NewMySqlError(fmt.Errorf("获取资源总数失败"))
// 	}

// 	return response.MenuListRsp{
// 		Total: count,
// 		Menus: rets,
// 	}, nil
// }

// Update 更新数据
func (l MenuSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.MenuUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !system.Menu.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("该ID对应的记录不存在"))
	}

	// 获取当前登陆用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户失败"))
	}

	oldData := new(model.Menu)
	err = system.Menu.Find(filter, oldData)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取记录失败: %s", err.Error()))
	}

	menu := model.Menu{
		Model:      oldData.Model,
		Name:       r.Name,
		Title:      r.Title,
		Icon:       r.Icon,
		Path:       r.Path,
		Redirect:   r.Redirect,
		Component:  r.Component,
		Sort:       r.Sort,
		Status:     r.Status,
		Hidden:     r.Hidden,
		NoCache:    r.NoCache,
		AlwaysShow: r.AlwaysShow,
		Breadcrumb: r.Breadcrumb,
		ActiveMenu: r.ActiveMenu,
		ParentId:   r.ParentId,
		Creator:    ctxUser.Username,
	}

	err = system.Menu.Update(&menu)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新记录失败: %s", err.Error()))
	}

	return nil, nil
}

// Delete 删除数据
func (l MenuSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.MenuDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.MenuIds {
		filter := tools.H{"id": int(id)}
		if !system.Menu.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("该ID对应的记录不存在"))
		}
	}

	// 删除接口
	err := system.Menu.Delete(r.MenuIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除接口失败: %s", err.Error()))
	}
	return nil, nil
}

// GetTree 数据树
func (l MenuSvc) GetTree(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	_, ok := req.(*request.MenuGetTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	menus, err := system.Menu.List()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取资源列表失败: " + err.Error()))
	}

	tree := system.GenMenuTree(0, menus)

	return tree, nil
}

// GetAccessTree 获取用户菜单树
func (l MenuSvc) GetAccessTree(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.MenuGetAccessTreeReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c
	// 校验
	filter := tools.H{"id": r.ID}
	if !system.User.Exist(filter) {
		return nil, tools.NewValidatorError(fmt.Errorf("该用户不存在"))
	}
	user := new(model.User)
	err := system.User.Find(filter, user)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("在MySQL查询用户失败: " + err.Error()))
	}
	var roleIds []uint
	for _, role := range user.Roles {
		roleIds = append(roleIds, role.ID)
	}
	menus, err := system.Menu.ListUserMenus(roleIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取资源列表失败: " + err.Error()))
	}

	tree := system.GenMenuTree(0, menus)

	return tree, nil
}
