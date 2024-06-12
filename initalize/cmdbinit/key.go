package cmdbinit

import (
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"fmt"
	"strings"
)

type KeyCmdb struct{}

// 创建
func (s KeyCmdb) Add(key *model.Key) error {
	return common.DB.Create(key).Error
}

// 是否重复key_name
func (s KeyCmdb) Checkcmdb(name string) (code int) {
	var Key *model.Key
	common.DB.Select("id").Where("key_name = ?", name).First(&Key)
	if Key.ID > 0 {
		return 1
	}
	return 200
}

// 获取key总数
func (s KeyCmdb) Count() (count int64, err error) {
	err = common.DB.Model(&model.Key{}).Count(&count).Error
	return count, err
}

// 删除
func (s KeyCmdb) Delete(ids []uint) error {
	return common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.Key{}).Error
}

// 获取key列表
func (s KeyCmdb) List(req *request.ListReq) (keys []*model.Key, err error) {
	db := common.DB.Model(&model.Key{}).Order("created_at DESC")
	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("key_name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	pageReq := tools.NewPageOption(req.PageNum, req.PageSize)
	err = db.Offset(pageReq.PageNum).Limit(pageReq.PageSize).Find(&keys).Error
	return keys, err
}

// 获取指定key
func (s KeyCmdb) Info(name string) ([]*model.Key, error) {
	var Key []*model.Key
	err := common.DB.Where("key_name = ?", name).First(&Key).Error
	return Key, err

}

func (s KeyCmdb) Infos(name string) (*model.Key, error) {
	var Key *model.Key
	err := common.DB.Where("key_name = ?", name).First(&Key).Error
	return Key, err

}

// 更新指定key
func (s KeyCmdb) Update(keys *model.Key) error {
	return common.DB.Model(&model.Key{}).Where("key_name= ?", keys.KeyName).Updates(keys).Error
}
