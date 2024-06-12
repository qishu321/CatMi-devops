package model

import (
	"gorm.io/gorm"
)

type CommandLog struct {
	gorm.Model
	Creator      string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
	Command      string `form:"command" json:"command"`
	SSHReqParams string `json:"SSHReqParams"`
	SShRsp       string `json:"SShRsp"`
	Ip           string `gorm:"type:varchar(20);comment:'Ip地址'" json:"ip"`
	Status       int    `gorm:"type:int(4);comment:'响应状态码'" json:"status"`
	StartTime    string `gorm:"type:varchar(2048);comment:'发起时间'" json:"startTime"`
}

// 创建执行模板
type Template struct {
	gorm.Model
	TemplateID int64  `gorm:"type:bigint;not null;comment:'模板id'" json:"templateid" `
	Name       string `gorm:"type:varchar(50);comment:'模板名称'" json:"name"`
	Stepnames  string `gorm:"type:text;comment:'步骤详情'" json:"steps"` // 存储为 JSON 格式
	Creator    string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}

// 执行日志
type Template_Log struct {
	gorm.Model
	TemplateID   int64  `gorm:"type:bigint;not null;comment:'模板id'" json:"templateid"`
	Name         string `gorm:"type:varchar(50);comment:'模板名称'" json:"name"`
	Stepname     string `gorm:"type:varchar(100);comment:'步骤名称'" json:"stepname"`
	Sort         int    `gorm:"type:int;default:999;comment:'步骤顺序(1-999)'" json:"sort"`
	Command      string `gorm:"type:text;comment:'执行命令'" json:"command"`
	Status       uint   `gorm:"type:tinyint(1);default:1;comment:'步骤状态(正常/失败, 默认正常)'" json:"status"`
	Cmdbnames    string `gorm:"type:text;comment:'步骤绑定的服务器列表'" json:"cmdbnames"`
	Timeouts     int    `gorm:"type:int;comment:'超时时间'" json:"timeouts"`
	SSHReqParams string `json:"ssh_req_params"`
	SshRsp       string `json:"ssh_rsp"`
	StartTime    string `gorm:"type:varchar(2048);comment:'发起时间'" json:"startTime"`
	Creator      string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}

// 执行脚本
type Task struct {
	gorm.Model
	Type        string `gorm:"type:varchar(50);comment:'脚本类型,bash,python'" json:"type"  form:"type"`
	Name        string `gorm:"type:varchar(50);comment:'脚本名称'" json:"name"`
	Description string `gorm:"type:text;comment:'脚本备注'" json:"description"  form:"description"`
	Script      string `json:"script" form:"script" gorm:"type:text;comment:'脚本内容'"`
	Creator     string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}

// 脚本参数
type TaskEnv struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);comment:'参数名称'" json:"name"`
	Description string `gorm:"type:text;comment:'参数备注'" json:"description"  form:"description"`
	Options     string `json:"options" form:"options" gorm:"type:text;comment:'参数内容'"`
	Important   int8   `gorm:"default:0;type:tinyint(1);column:important;comment:'是否标记为必填 1:是 0: 否'" json:"important"`
	Creator     string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}
