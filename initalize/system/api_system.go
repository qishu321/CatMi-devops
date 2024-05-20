package system

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type ApiSystem struct{}

// List 获取数据列表
func (s ApiSystem) List(req *request.ApiListReq) ([]*model.Api, error) {
	var list []*model.Api
	db := common.DB.Model(&model.Api{}).Order("created_at DESC")

	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err := db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&list).Error
	return list, err
}

// List 获取数据列表
func (s ApiSystem) ListAll() (list []*model.Api, err error) {
	err = common.DB.Model(&model.Api{}).Order("created_at DESC").Find(&list).Error

	return list, err
}

// Count 获取数据总数
func (s ApiSystem) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&model.Api{}).Count(&count).Error
	return count, err
}

// Add 添加资源
func (s ApiSystem) Add(api *model.Api) error {
	return common.DB.Create(api).Error
}

// Update 更新资源
func (s ApiSystem) Update(api *model.Api) error {
	// 根据id获取接口信息
	var oldApi model.Api
	err := common.DB.First(&oldApi, api.ID).Error
	if err != nil {
		return errors.New("根据接口ID获取接口信息失败")
	}
	err = common.DB.Model(api).Where("id = ?", api.ID).Updates(api).Error
	if err != nil {
		return err
	}
	// 更新了method和path就更新casbin中policy
	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies, _ := common.CasbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		// 接口在casbin的policy中存在才进行操作
		if len(policies) > 0 {
			// 先删除
			isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("更新权限接口失败")
			}
			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}
			// 新增
			isAdded, _ := common.CasbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("更新权限接口失败")
			}
			// 加载policy
			err := common.CasbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("更新权限接口成功，权限接口策略加载失败")
			} else {
				return err
			}
		}
	}
	return err
}

// Find 获取单个资源
func (s ApiSystem) Find(filter map[string]interface{}, data *model.Api) error {
	return common.DB.Where(filter).First(&data).Error
}

// Exist 判断资源是否存在
func (s ApiSystem) Exist(filter map[string]interface{}) bool {
	var dataObj model.Api
	err := common.DB.Debug().Order("created_at DESC").Where(filter).First(&dataObj).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Delete 批量删除
func (s ApiSystem) Delete(ids []uint) error {
	var apis []model.Api
	for _, id := range ids {
		// 根据ID获取用户
		api := new(model.Api)
		err := s.Find(tools.H{"id": id}, api)
		if err != nil {
			return fmt.Errorf("根据ID获取接口信息失败: %v", err)
		}
		apis = append(apis, *api)
	}

	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.Api{}).Error
	// 如果删除成功，删除casbin中policy
	if err == nil {
		for _, api := range apis {
			policies, _ := common.CasbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("删除权限接口失败")
				}
			}
		}
		// 重新加载策略
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("删除权限接口成功，权限接口策略加载失败")
		} else {
			return err
		}
	}
	return err
}


// GetApisById 根据接口ID获取接口列表
func (s ApiSystem) GetApisById(apiIds []uint) ([]*model.Api, error) {
	var apis []*model.Api
	err := common.DB.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}
