package deployinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
)

type CommandLog struct{}

// 创建
func (s CommandLog) Add(server *model.CommandLog) error {
	return common.DB.Create(server).Error
}

// 获取server列表
func (s CommandLog) List(req *request.SShListReq) (log []*model.CommandLog, err error) {
	db := common.DB.Model(&model.CommandLog{}).Order("created_at DESC")
	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&log).Error
	return log, err
}

// 获取server总数
func (s CommandLog) Count() (count int64, err error) {
	err = common.DB.Model(&model.CommandLog{}).Count(&count).Error
	return count, err
}

// 删除
func (s CommandLog) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()

	// // 查找所有要删除的服务器记录，并预加载其关联的 Keys
	// if err := tx.Preload("Keys").Where("id IN (?)", ids).Find(&servers).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// 删除服务器记录
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.CommandLog{}).Error; err != nil {
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
