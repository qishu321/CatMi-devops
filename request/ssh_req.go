package request

import "time"

type SSHClientConfigReq struct {
	UserName   string        `form:"username" json:"username" binding:"required"`
	Password   string        `form:"password" json:"password"`
	PublicIP   string        `form:"public_ip" json:"public_ip" binding:"required"`
	Port       int           `form:"port" json:"port" binding:"required"`
	Command    string        `form:"command" json:"command"`
	AuthModel  string        `form:"authmodel" json:"authmodel" binding:"required"`
	PrivateKey string        `form:"private_key" json:"private_key"`
	Timeout    time.Duration `form:"timeout" json:"timeout"`   //超时时间
	TimeoutS   int           `form:"timeouts" json:"timeouts"` //超时时间

}
type TemplateReq struct {
	TemplateID int64      `gorm:"type:bigint;not null;comment:'模板id'" json:"templateid" `
	Name       string     `gorm:"type:varchar(50);comment:'模板名称'" json:"name" validate:"required"`
	Stepnames  []Stepname `gorm:"-" json:"steps" validate:"required,dive"`
}

type Stepname struct {
	Stepname  string   `gorm:"type:varchar(100);comment:'步骤名称'" json:"stepname" validate:"required"`
	Sort      int      `gorm:"type:int;default:999;comment:'步骤顺序(1-999)'" json:"sort" validate:"required"`
	Command   string   `gorm:"type:text;comment:'执行命令'" json:"command" `
	Cmdbnames []string `gorm:"type:text;comment:'步骤绑定的服务器列表'" json:"cmdbnames" `
	Timeouts  int      `gorm:"type:int;comment:'超时时间'" json:"timeouts"`
}

type UpdateTemplateReq struct {
	ID        uint       `json:"id" db:"id" form:"id"  validate:"required"`
	Name      string     `gorm:"type:varchar(50);comment:'模板名称'" json:"name"`
	Stepnames []Stepname `gorm:"-" json:"steps" validate:"required,dive"`
}

type SSHReqParams struct {
	Command  string       `form:"command" json:"command"  binding:"required"`
	TimeoutS int          `form:"timeouts" json:"timeouts"` //超时时间
	Data     []ServerInfo `json:"datalist"`
}

type ServerInfo struct {
	Cmdbname  string `json:"cmdbName" db:"cmdbName" form:"cmdbName" comment:"CMDB名称" validate:"required"`
	PublicIP  string `json:"publicIP"`
	SSHPort   int    `json:"sshPort"`
	ID        uint   `json:"ID"`
	AuthModel string `form:"authmodel" json:"authmodel" binding:"required"`
	Keys      []struct {
		ServerName string `json:"serverName" form:"serverName" comment:"服务器用户名"`
		Password   string `json:"password" form:"password" comment:"密码"`
		PrivateKey string `json:"private_key" form:"private_key" comment:"私钥"`
	}
}
type TaskReq struct {
	Type        string            ` json:"type"  form:"type" validate:"required"`
	Name        string            ` json:"name" form:"name" validate:"required"`
	Script      string            `json:"script" form:"script" validate:"required"`
	Description string            `json:"description"  form:"description"`
	Variables   map[string]string `json:"variables" form:"variables"`
}
type UpdateTaskReq struct {
	ID          uint              `json:"id" db:"id" form:"id"  validate:"required"`
	Type        string            ` json:"type"  form:"type" validate:"required"`
	Name        string            ` json:"name" form:"name" validate:"required"`
	Script      string            `json:"script" form:"script" validate:"required"`
	Description string            `json:"description"  form:"description"`
	Variables   map[string]string `json:"variables" form:"variables"`
}

type SShListReq struct {
	Page
	Name       string `json:"name"`
	Searchname string `json:"searchname"`
	Creator    string `json:"creator" form:"creator" comment:"创建人"`
}

type SShInfoReq struct {
	Name string `json:"name"`
}

type SShDeleteReq struct {
	Ids []uint `json:"Ids" validate:"required"`
}

type Options struct {
	Environment map[string]string `json:"environment"`
}

type TaskEnvReq struct {
	Name        string              ` json:"name" form:"name" validate:"required"`
	Options     []map[string]string `json:"options" form:"options" validate:"required"`
	Description string              `json:"description"  form:"description"`
	Important   int8                ` json:"important" form:"important"`
}

type UpdateTaskEnvReq struct {
	ID          uint                `json:"id" db:"id" form:"id"  validate:"required"`
	Name        string              ` json:"name" form:"name" validate:"required"`
	Options     []map[string]string `json:"options" form:"options" validate:"required"`
	Description string              `json:"description"  form:"description"`
	Important   int8                ` json:"important" form:"important"`
}
