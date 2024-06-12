package deployinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type TaskEnv struct{}

// 创建
func (s TaskEnv) Add(server *model.TaskEnv) error {
	return common.DB.Create(server).Error
}
func (s TaskEnv) Check(name string) (code int) {
	var server *model.TaskEnv
	common.DB.Select("id").Where("name = ?", name).First(&server)
	if server.ID > 0 {
		return 1
	}
	return 200
}
func (s TaskEnv) UpdateCheck(name string, id uint) bool {
	var count int64
	common.DB.Model(&model.TaskEnv{}).Where("name = ? AND id != ?", name, id).Count(&count)
	return count > 0
}

// 获取server列表
func (s TaskEnv) List(req *request.SShListReq) (log []*model.TaskEnv, err error) {
	db := common.DB.Model(&model.TaskEnv{}).Order("created_at DESC")
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
func (s TaskEnv) Update(servers *model.TaskEnv) error {
	// 开启事务
	tx := common.DB.Begin()
	if err := tx.Model(&model.TaskEnv{}).Where("id = ?", servers.ID).Updates(servers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Update servers: %w", err)
	}
	// 提交事务
	return tx.Commit().Error

}

// 获取指定server
func (s TaskEnv) Info(name string) (*model.TaskEnv, error) {
	var server *model.TaskEnv
	err := common.DB.Where("name = ?", name).First(&server).Error
	return server, err

}

// 获取server总数
func (s TaskEnv) Count() (count int64, err error) {
	err = common.DB.Model(&model.TaskEnv{}).Count(&count).Error
	return count, err
}

// 删除
func (s TaskEnv) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.TaskEnv{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
