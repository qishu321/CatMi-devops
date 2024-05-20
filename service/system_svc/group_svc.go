package system_svc

import (
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type GroupSvc struct{}

// Add 添加数据
func (l GroupSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupAddReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	if system.Group.Exist(tools.H{"group_name": r.GroupName}) {
		return nil, tools.NewValidatorError(fmt.Errorf("该分组对应DN已存在"))
	}

	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	group := model.Group{
		ParentId:  r.ParentId,
		GroupName: r.GroupName,
		Remark:    r.Remark,
		Creator:   ctxUser.Username,
		Source:    "platform", //默认是平台添加
	}

	// 然后在数据库中创建组
	err = system.Group.Add(&group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("向MySQL创建分组失败"))
	}

	return nil, nil
}

// List 数据列表
func (l GroupSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	// 获取数据列表
	groups, err := system.Group.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组列表失败: %s", err.Error()))
	}

	rets := make([]model.Group, 0)
	for _, group := range groups {
		rets = append(rets, *group)
	}
	count, err := system.Group.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组总数失败"))
	}

	return response.GroupListRsp{
		Total:  count,
		Groups: rets,
	}, nil
}

// GetTree 数据树
func (l GroupSvc) GetTree(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupListReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	var groups []*model.Group
	groups, err := system.Group.ListTree(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取资源列表失败: " + err.Error()))
	}

	tree := system.GenGroupTree(0, groups)

	return tree, nil
}

// Update 更新数据
func (l GroupSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupUpdateReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": int(r.ID)}
	if !system.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	// 获取当前登陆用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户失败"))
	}

	oldGroup := new(model.Group)
	err = system.Group.Find(filter, oldGroup)
	if err != nil {
		return nil, tools.NewMySqlError(err)
	}

	newGroup := model.Group{
		Model:     oldGroup.Model,
		GroupName: r.GroupName,
		Remark:    r.Remark,
		Creator:   ctxUser.Username,
	}

	err = system.Group.Update(&newGroup)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("向MySQL更新分组失败"))
	}
	return nil, nil
}

// Delete 删除数据
func (l GroupSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupDeleteReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	for _, id := range r.GroupIds {
		filter := tools.H{"id": int(id)}
		if !system.Group.Exist(filter) {
			return nil, tools.NewMySqlError(fmt.Errorf("有分组不存在"))
		}
	}

	groups, err := system.Group.GetGroupByIds(r.GroupIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组列表失败: %s", err.Error()))
	}

	// 从MySQL中删除
	err = system.Group.Delete(groups)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除接口失败: %s", err.Error()))
	}

	return nil, nil
}

// AddUser 添加用户到分组
func (l GroupSvc) AddUser(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupAddUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !system.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	users, err := system.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取用户列表失败: %s", err.Error()))
	}

	group := new(model.Group)
	err = system.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	// 先添加到MySQL
	err = system.Group.AddUserToGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("添加用户到分组失败: %s", err.Error()))
	}

	return nil, nil
}

func (l GroupSvc) updataUser(newUser *model.User) error {
	err := system.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("在MySQL更新用户失败：" + err.Error()))
	}
	return nil
}

// RemoveUser 移除用户
func (l GroupSvc) RemoveUser(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.GroupRemoveUserReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !system.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	users, err := system.User.GetUserByIds(r.UserIds)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取用户列表失败: %s", err.Error()))
	}

	group := new(model.Group)
	err = system.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	// 再操作MySQL
	err = system.Group.RemoveUserFromGroup(group, users)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("将用户从MySQL移除失败: %s", err.Error()))
	}

	for _, user := range users {
		oldData := new(model.User)
		err = system.User.Find(tools.H{"id": user.ID}, oldData)
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
		newData := oldData

		var newDepts []string
		var newDeptIds []string
		// 删掉移除的分组名字
		for _, v := range strings.Split(oldData.Departments, ",") {
			if v != group.GroupName {
				newDepts = append(newDepts, v)
			}
		}
		// 删掉移除的分组id
		for _, v := range strings.Split(oldData.DepartmentId, ",") {
			if v != strconv.Itoa(int(r.GroupID)) {
				newDeptIds = append(newDeptIds, v)
			}
		}

		newData.Departments = strings.Join(newDepts, ",")
		newData.DepartmentId = strings.Join(newDeptIds, ",")
		err = l.updataUser(newData)
		if err != nil {
			return nil, tools.NewOperationError(fmt.Errorf("处理用户的部门数据失败:" + err.Error()))
		}
	}

	return nil, nil
}

// UserInGroup 在分组内的用户
func (l GroupSvc) UserInGroup(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.UserInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !system.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	group := new(model.Group)
	err := system.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	rets := make([]response.Guser, 0)

	for _, user := range group.Users {
		if r.Nickname != "" && !strings.Contains(user.Nickname, r.Nickname) {
			continue
		}
		rets = append(rets, response.Guser{
			UserId:       int64(user.ID),
			UserName:     user.Username,
			NickName:     user.Nickname,
			Mail:         user.Mail,
			JobNumber:    user.JobNumber,
			Mobile:       user.Mobile,
			Introduction: user.Introduction,
		})
	}

	return response.GroupUsers{
		GroupId:     int64(group.ID),
		GroupName:   group.GroupName,
		GroupRemark: group.Remark,
		UserList:    rets,
	}, nil
}

// UserNoInGroup 不在分组内的用户
func (l GroupSvc) UserNoInGroup(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.UserNoInGroupReq)
	if !ok {
		return nil, ReqAssertErr
	}
	_ = c

	filter := tools.H{"id": r.GroupID}

	if !system.Group.Exist(filter) {
		return nil, tools.NewMySqlError(fmt.Errorf("分组不存在"))
	}

	group := new(model.Group)
	err := system.Group.Find(filter, group)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取分组失败: %s", err.Error()))
	}

	var userList []*model.User
	userList, err = system.User.ListAll()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取资源列表失败: " + err.Error()))
	}

	rets := make([]response.Guser, 0)
	for _, user := range userList {
		in := true
		for _, groupUser := range group.Users {
			if user.Username == groupUser.Username {
				in = false
				break
			}
		}
		if in {
			if r.Nickname != "" && !strings.Contains(user.Nickname, r.Nickname) {
				continue
			}
			rets = append(rets, response.Guser{
				UserId:       int64(user.ID),
				UserName:     user.Username,
				NickName:     user.Nickname,
				Mail:         user.Mail,
				JobNumber:    user.JobNumber,
				Mobile:       user.Mobile,
				Introduction: user.Introduction,
			})
		}
	}

	return response.GroupUsers{
		GroupId:     int64(group.ID),
		GroupName:   group.GroupName,
		GroupRemark: group.Remark,
		UserList:    rets,
	}, nil
}