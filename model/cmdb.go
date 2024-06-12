package model

import "gorm.io/gorm"

type ServerGroup struct {
	gorm.Model
	GroupName   string        `json:"groupName" gorm:"column:group_name"`
	ServerCmdbs []*ServerCmdb `gorm:"many2many:group_server_cmdbs" json:"serverCmdbs"`
	Desc        string        `gorm:"type:varchar(50);column:desc;comment:'描述'" json:"desc"`
	Creator     string        `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}
type ServerCmdb struct {
	gorm.Model
	CmdbID    int64  `gorm:"type:bigint;not null;comment:'CMDBID'" json:"cmdbId" validate:"required"`
	CmdbName  string `gorm:"column:cmdb_name;comment:'CMDB名称'" json:"cmdbName"`
	PublicIP  string `gorm:"column:public_ip;comment:'IP地址'" json:"publicIP"`
	Keys      []*Key `gorm:"many2many:server_cmdb_keys;comment:'关联的密钥'" json:"keys"`
	SSHPort   int    `gorm:"column:ssh_port;comment:'SSH端口号'" json:"sshPort"`
	Enabled   int8   `gorm:"type:tinyint(1);default:0;column:enabled;comment:'是否连接成功 1：是 0： 否'" json:"enabled"`
	AuthModel string `gorm:"column:authmodel;comment:'连接服务器所使用的是密钥还是密码'" json:"authmodel" `
	Desc      string `gorm:"type:varchar(50);column:desc;comment:'描述'" json:"desc"`
	Label     string `gorm:"column:label;comment:'标签'" json:"label"`
	Creator   string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}

type Key struct {
	gorm.Model
	KeyName    string `gorm:"type:varchar(128);comment:'key名称'" json:"keyName"`
	ServerName string `gorm:"type:varchar(128);comment:'服务器用户名'" json:"serverName"`
	Password   string `json:"password" db:"password" form:"password"`
	PrivateKey string `json:"private_key" db:"private_key" form:"private_key"`
	Creator    string `gorm:"type:varchar(20);comment:'创建人'" json:"creator"`
}
