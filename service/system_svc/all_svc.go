package system_svc

import (
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/utils/tools"
	"fmt"
)

var (
	ReqAssertErr = tools.NewRspError(tools.SystemErr, fmt.Errorf("请求异常"))

	Api          = &ApiSvc{}
	User         = &UserSvc{}

	Group        = &GroupSvc{}
	Role         = &RoleSvc{}
	Menu         = &MenuSvc{}
	OperationLog = &OperationLogSvc{}
	Base         = &BaseSvc{}
)

// CommonAddGroup 标准创建分组
func CommonAddGroup(group *model.Group) error {
	// 在数据库中创建组
	err := system.Group.Add(group)
	if err != nil {
		return err
	}

	// 默认创建分组之后，需要将admin添加到分组中
	adminInfo := new(model.User)
	err = system.User.Find(tools.H{"id": 1}, adminInfo)
	if err != nil {
		return err
	}

	err = system.Group.AddUserToGroup(group, []model.User{*adminInfo})
	if err != nil {
		return err
	}

	return nil
}

// CommonUpdateGroup 标准更新分组
func CommonUpdateGroup(oldGroup, newGroup *model.Group) error {
	err := system.Group.Update(newGroup)
	if err != nil {
		return err
	}
	return nil
}

// CommonAddUser 标准创建用户
func CommonAddUser(user *model.User, groups []*model.Group) error {
	// 用户信息的预置处理
	if user.Nickname == "" {
		user.Nickname = "佚名"
	}
	if user.GivenName == "" {
		user.GivenName = user.Nickname
	}
	if user.Introduction == "" {
		user.Introduction = user.Nickname
	}
	if user.Mail == "" {
		user.Mail = "该用户邮箱为空"
	}
	if user.JobNumber == "" {
		user.JobNumber = "该用户工号为空"
	}
	if user.Departments == "" {
		user.Departments = "默认:研发中心"
	}
	if user.Position == "" {
		user.Position = "默认:技术"
	}
	if user.PostalAddress == "" {
		user.PostalAddress = "默认:地球"
	}
	if user.Mobile == "" {
		user.Mobile = "emptyMobile"
	}

	// 先将用户添加到MySQL
	err := system.User.Add(user)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("向MySQL创建用户失败：" + err.Error()))
	}

	// 处理用户归属的组
	for _, group := range groups {
		// 先将用户和部门信息维护到MySQL
		err := system.Group.AddUserToGroup(group, []model.User{*user})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("向MySQL添加用户到分组关系失败：" + err.Error()))
		}
	}
	return nil
}

// CommonUpdateUser 标准更新用户
func CommonUpdateUser(oldUser, newUser *model.User, groupId []uint) error {
	err := system.User.Update(newUser)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("在MySQL更新用户失败：" + err.Error()))
	}

	//判断部门信息是否有变化有变化则更新相应的数据库
	oldDeptIds := tools.StringToSlice(oldUser.DepartmentId, ",")
	addDeptIds, removeDeptIds := tools.ArrUintCmp(oldDeptIds, groupId)

	// 先处理添加的部门
	addgroups, err := system.Group.GetGroupByIds(addDeptIds)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("根据部门ID获取部门信息失败" + err.Error()))
	}
	for _, group := range addgroups {
		// 先将用户和部门信息维护到MySQL
		err := system.Group.AddUserToGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("向MySQL添加用户到分组关系失败：" + err.Error()))
		}
	}

	// 再处理删除的部门
	removegroups, err := system.Group.GetGroupByIds(removeDeptIds)
	if err != nil {
		return tools.NewMySqlError(fmt.Errorf("根据部门ID获取部门信息失败" + err.Error()))
	}
	for _, group := range removegroups {
		err := system.Group.RemoveUserFromGroup(group, []model.User{*newUser})
		if err != nil {
			return tools.NewMySqlError(fmt.Errorf("在MySQL将用户从分组移除失败：" + err.Error()))
		}
	}
	return nil
}