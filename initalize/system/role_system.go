package system

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type RoleSystem struct{}

// Exist 判断资源是否存在
func (s RoleSystem) Exist(filter map[string]interface{}) bool {
	var dataObj model.Role
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// List 获取数据列表
func (s RoleSystem) List(req *request.RoleListReq) ([]*model.Role, error) {
	var list []*model.Role
	db := common.DB.Model(&model.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// Count 获取资源总数
func (s RoleSystem) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Role{}).Count(&count).Error
	return count, err
}

// Add 创建资源
func (s RoleSystem) Add(role *model.Role) error {
	return common.DB.Create(role).Error
}

// Update 更新资源
func (s RoleSystem) Update(role *model.Role) error {
	return common.DB.Model(&model.Role{}).Where("id = ?", role.ID).Updates(role).Error
}

// Find 获取单个资源
func (s RoleSystem) Find(filter map[string]interface{}, data *model.Role) error {
	return common.DB.Where(filter).First(&data).Error
}

// Delete 删除资源
func (s RoleSystem) Delete(roleIds []uint) error {
	var roles []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies, _ := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("删除角色成功, 删除角色关联权限接口失败")
				}
			}
		}

	}
	return err
}

// Delete 根据角色ID获取角色
func (s RoleSystem) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// GetRoleMenusById 获取角色的权限菜单
func (s RoleSystem) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := common.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// UpdateRoleMenus 更新角色的权限菜单
func (s RoleSystem) UpdateRoleMenus(role *model.Role) error {
	return common.DB.Model(role).Association("Menus").Replace(role.Menus)
}

// UpdateRoleApis 更新角色的权限接口（先全部删除再新增）
func (s RoleSystem) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies, _ := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("更新角色的权限接口失败")
		}
	}
	isAdded, _ := common.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("更新角色的权限接口失败")
	}
	err = common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
	} else {
		return err
	}
}
