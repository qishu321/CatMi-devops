package deployinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type Template_Log struct{}

// 创建 enable接口，然后批量写入进来
func (s Template_Log) Add(server *model.Template_Log) error {
	return common.DB.Create(server).Error
}

// 获取server列表
func (s Template_Log) List(req *request.SShListReq) (log []*model.Template_Log, err error) {
	subQuery := common.DB.Model(&model.Template_Log{}).Select("MIN(id) as id").Group("name")
	db := common.DB.Model(&model.Template_Log{}).Where("id IN (?)", subQuery).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&log).Error
	return log, err
}

// 获取指定server
func (s Template_Log) Info(name string) ([]*model.Template_Log, error) {
	var server []*model.Template_Log
	err := common.DB.Where("name = ?", name).Find(&server).Error
	return server, err

}

// 获取server总数
func (s Template_Log) CountDistinctNames() (count int64, err error) {
	err = common.DB.Model(&model.Template_Log{}).Distinct("name").Count(&count).Error
	return count, err
}

// 删除
func (s Template_Log) Delete(ids []uint) error {
	// 开启事务
	tx := common.DB.Begin()

	// // 查找所有要删除的服务器记录，并预加载其关联的 Keys
	// if err := tx.Preload("Keys").Where("id IN (?)", ids).Find(&servers).Error; err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// 删除服务器记录
	if err := tx.Where("id IN (?)", ids).Unscoped().Delete(&model.Template_Log{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 提交事务
	return tx.Commit().Error
}
