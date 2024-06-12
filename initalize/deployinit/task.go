package deployinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type Task struct{}

// 创建
func (s Task) Add(server *model.Task) error {
	return common.DB.Create(server).Error
}
func (s Task) Check(name string) (code int) {
	var server *model.Task
	common.DB.Select("id").Where("name = ?", name).First(&server)
	if server.ID > 0 {
		return 1
	}
	return 200
}
func (s Task) UpdateCheck(name string, id uint) bool {
	var count int64
	common.DB.Model(&model.Task{}).Where("name = ? AND id != ?", name, id).Count(&count)
	return count > 0
}

// 获取server列表
func (s Task) List(req *request.SShListReq) (log []*model.Task, err error) {
	db := common.DB.Model(&model.Task{}).Order("created_at DESC")
	name := strings.TrimSpace(req.Searchname)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&log).Error
	return log, err
}
func (s Task) Update(servers *model.Task) error {
	// 开启事务
	tx := common.DB.Begin()
	if err := tx.Model(&model.Task{}).Where("id = ?", servers.ID).Updates(servers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Update servers: %w", err)
	}
	// 提交事务
	return tx.Commit().Error

}

// 获取指定server
func (s Task) Info(name string) (*model.Task, error) {
	var server *model.Task
	err := common.DB.Where("name = ?", name).First(&server).Error
	return server, err

}

// 获取server总数
func (s Task) Count() (count int64, err error) {
	err = common.DB.Model(&model.Task{}).Count(&count).Error
	return count, err
}

// 删除
func (s Task) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.Task{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
