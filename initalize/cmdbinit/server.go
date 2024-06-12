package cmdbinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type ServerCmdb struct{}

// 创建
func (s ServerCmdb) Add(server *model.ServerCmdb) error {
	return common.DB.Preload("Keys").Create(server).Error
}

// AddServerKey 添加key到server
func (s ServerCmdb) AddServerKey(server *model.ServerCmdb, key []model.Key) error {
	return common.DB.Model(&server).Association("Keys").Append(key)
}

// 是否重复server_name
func (s ServerCmdb) Check(name string) (code int) {
	var server *model.ServerCmdb
	common.DB.Select("id").Where("cmdb_name = ?", name).First(&server)
	if server.ID > 0 {
		return 1
	}
	return 200
}

// 获取server总数
func (s ServerCmdb) Count() (count int64, err error) {
	err = common.DB.Model(&model.ServerCmdb{}).Count(&count).Error
	return count, err
}

// 删除
func (s ServerCmdb) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()

	// // 查找所有要删除的服务器记录，并预加载其关联的 Keys
	// if err := tx.Preload("Keys").Where("id IN (?)", ids).Find(&servers).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// 删除服务器记录
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.ServerCmdb{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 清空每个服务器记录的关联的 Keys
	err := tx.Exec("DELETE FROM server_cmdb_keys WHERE server_cmdb_id IN (?)", ids).Error
	if err != nil {
		return err
	}
	// 提交事务
	return tx.Commit().Error
}

// 获取server列表
func (s ServerCmdb) List(req *request.ListReq) (servers []*model.ServerCmdb, err error) {
	db := common.DB.Model(&model.ServerCmdb{}).Order("created_at DESC")
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("cmdb_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("Keys").Find(&servers).Error
	return servers, err
}

// 获取已启用server列表
func (s ServerCmdb) EnableList(req *request.ListReq) (servers []*model.ServerCmdb, err error) {
	db := common.DB.Model(&model.ServerCmdb{}).Order("created_at DESC")
	db = db.Where("enabled = ?", 1)
	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Preload("Keys").Find(&servers).Error
	return servers, err
}

func (s ServerCmdb) EnableCount() (count int64, err error) {
	err = common.DB.Model(&model.ServerCmdb{}).Where("enabled = ?", 1).Count(&count).Error
	return count, err
}

// 获取指定server
func (s ServerCmdb) Info(name string) ([]*model.ServerCmdb, error) {
	var server []*model.ServerCmdb
	err := common.DB.Preload("Keys").Where("cmdb_name = ?", name).First(&server).Error
	return server, err

}

func (s ServerCmdb) EnabledUpdate(servers *model.ServerCmdb) error {
	// 开启事务
	tx := common.DB.Begin()
	if err := tx.Model(&model.ServerCmdb{}).Where("id = ? and cmdb_id = ?", servers.ID, servers.CmdbID).Updates(servers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Update servers: %w", err)
	}
	// 提交事务
	return tx.Commit().Error

}

// 更新指定server
func (s ServerCmdb) Update(servers *model.ServerCmdb, keyid int64) error {
	// 开启事务
	tx := common.DB.Begin()
	err := tx.Exec("DELETE FROM server_cmdb_keys WHERE server_cmdb_id IN (?)", servers.ID).Error
	if err != nil {
		return err
	}
	if err := tx.Model(&model.ServerCmdb{}).Where("id = ? and cmdb_id = ?", servers.ID, servers.CmdbID).Updates(servers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Update servers: %w", err)
	}
	err = tx.Exec("INSERT INTO server_cmdb_keys (server_cmdb_id, key_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE server_cmdb_id = server_cmdb_id", servers.ID, keyid).Error
	if err != nil {
		return err
	}
	// 提交事务
	return tx.Commit().Error

}
