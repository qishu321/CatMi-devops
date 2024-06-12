package deployinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type Template struct{}

// 创建
func (s Template) Add(server *model.Template) error {
	return common.DB.Create(server).Error
}
func (s Template) Check(name string) (code int) {
	var server *model.Template
	common.DB.Select("id").Where("name = ?", name).First(&server)
	if server.ID > 0 {
		return 1
	}
	return 200
}
func (s Template) UpdateCheck(name string, id uint) bool {
	var count int64
	common.DB.Model(&model.Template{}).Where("name = ? AND id != ?", name, id).Count(&count)
	return count > 0
}

// 获取server列表
func (s Template) List(req *request.SShListReq) (log []*model.Template, err error) {
	db := common.DB.Model(&model.Template{}).Order("created_at DESC")
	name := strings.TrimSpace(req.Searchname)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&log).Error
	return log, err
}
func (s Template) Update(servers *model.Template) error {
	// 开启事务
	tx := common.DB.Begin()
	if err := tx.Model(&model.Template{}).Where("id = ?", servers.ID).Updates(servers).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to Update servers: %w", err)
	}
	// 提交事务
	return tx.Commit().Error

}

// 获取指定server
func (s Template) Info(name string) ([]*model.Template, error) {
	var server []*model.Template
	err := common.DB.Where("name = ?", name).First(&server).Error
	return server, err

}

// 获取server总数
func (s Template) Count() (count int64, err error) {
	err = common.DB.Model(&model.Template{}).Count(&count).Error
	return count, err
}

// 删除
func (s Template) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()

	// // 查找所有要删除的服务器记录，并预加载其关联的 Keys
	// if err := tx.Preload("Keys").Where("id IN (?)", ids).Find(&servers).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// 删除服务器记录
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.Template{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
