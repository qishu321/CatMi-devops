package cmdbinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type ServerGroup struct{}

// 创建
func (s ServerGroup) Add(server *model.ServerGroup, ids []int64) error {
	// 在数据库中创建服务器组记录
	if err := common.DB.Create(server).Error; err != nil {
		return err
	}
	// 遍历ids数组，将每个ID与服务器组建立关联
	for _, id := range ids {
		// 执行插入操作，如果存在主键冲突则执行更新操作
		err := common.DB.Exec("INSERT INTO group_server_cmdbs (server_group_id, server_cmdb_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE server_group_id = server_group_id", server.ID, id).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// 是否重复group_name
func (s ServerGroup) Check(name string) (code int) {
	var server *model.ServerGroup
	common.DB.Select("id").Where("group_name = ?", name).First(&server)
	if server.ID > 0 {
		return 1
	}
	return 200
}

// 获取server组总数
func (s ServerGroup) Count() (count int64, err error) {
	err = common.DB.Model(&model.ServerGroup{}).Count(&count).Error
	return count, err
}

// 删除
func (s ServerGroup) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()
	// 删除服务器组记录
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.ServerGroup{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 清空每个服务器组记录的关联的 Keys
	err := tx.Exec("DELETE FROM group_server_cmdbs WHERE server_group_id IN (?)", ids).Error
	if err != nil {
		return err
	}
	// 提交事务
	return tx.Commit().Error
}

// 获取server组列表
func (s ServerGroup) List(req *request.ListReq) (servers []*model.ServerGroup, err error) {
	db := common.DB.Model(&model.ServerGroup{}).Order("created_at DESC")
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("group_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("ServerCmdbs").Find(&servers).Error
	return servers, err
}

// 获取指定server组
func (s ServerGroup) Info(name string) ([]*model.ServerGroup, error) {
	var server []*model.ServerGroup
	err := common.DB.Preload("ServerCmdbs.Keys").Where("group_name = ?", name).First(&server).Error
	return server, err
}

// 更新指定server组
func (s ServerGroup) Update(group *model.ServerGroup, ids []int64) error {
	// 开启事务
	tx := common.DB.Begin()
	err := tx.Exec("DELETE FROM group_server_cmdbs WHERE server_group_id IN (?)", group.ID).Error
	if err != nil {
		return err
	}
	if err := tx.Model(&model.ServerGroup{}).Where("id = ?", group.ID).Updates(group).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Update servers: %w", err)
	}
	// 遍历ids数组，将每个ID与服务器组建立关联
	for _, id := range ids {
		// 执行插入操作，如果存在主键冲突则执行更新操作
		err := tx.Exec("INSERT INTO group_server_cmdbs (server_group_id, server_cmdb_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE server_group_id = server_group_id", group.ID, id).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 提交事务
	return tx.Commit().Error

}

// 删除指定server组
func (s ServerGroup) DelGroupId(group *model.ServerGroup, ids []int64) error {
	// 开启事务
	tx := common.DB.Begin()
	err := tx.Exec("DELETE FROM group_server_cmdbs WHERE server_group_id = ? AND server_cmdb_id IN (?)", group.ID, ids).Error

	if err != nil {
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
