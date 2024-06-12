package request

type Page struct {
	PageNum  int `json:"pageNum" form:"pageNum" comment:"分页页码"`
	PageSize int `json:"pageSize" form:"pageSize" comment:"分页大小"`
}
type ServerGroupReq struct {
	GroupName string  `json:"groupName" form:"groupName" comment:"服务器组名称" validate:"required"`
	Desc      string  `json:"desc" form:"desc" comment:"描述"`
	Ids       []int64 `json:"ids"`
}
type UpdateServerGroupReq struct {
	ID        uint    `json:"id" db:"id" form:"id"  validate:"required"`
	GroupName string  `json:"groupName" form:"groupName" comment:"服务器组名称" validate:"required"`
	Desc      string  `json:"desc" form:"desc" comment:"描述"`
	Ids       []int64 `json:"ids"`
}

type DelGroupIdServerGroupReq struct {
	ID  uint    `json:"id" db:"id" form:"id"  validate:"required"`
	Ids []int64 `json:"ids" validate:"required"`
}

type EnabledCmdbReqParams struct {
	AuthModel string `form:"authmodel" json:"authmodel" binding:"required"`
	Data      []struct {
		KeyNames string `json:"keyName"`
		PublicIP string `json:"publicIP"`
		SSHPort  int    `json:"sshPort"`
		ID       uint   `json:"ID"`
		Cmdbid   int64  `json:"cmdbId"`
		Keys     []struct {
			ServerName string `json:"serverName" form:"serverName" comment:"服务器用户名"`
			Password   string `json:"password" form:"password" comment:"密码"`
			PrivateKey string `json:"privateKey" form:"privateKey" comment:"私钥"`
		}
	} `json:"data"`
}

type EnabledServerCmdbReq struct {
	ID        uint   `json:"id" db:"id" form:"id"`
	Cmdbid    int64  `json:"cmdbid"  comment:"CMDB ID"`
	PublicIP  string `json:"public_ip" db:"public_ip" form:"public_ip" comment:"IP地址" validate:"required"`
	KeyNames  string `json:"keyName" form:"keyName" comment:"key名称" validate:"required"`
	SSHPort   int    `json:"ssh_port" db:"ssh_port" form:"ssh_port" comment:"SSH端口号" validate:"required"`
	Enabled   int8   `json:"enabled" db:"enabled" form:"enabled" validate:"oneof=0 1"`
	AuthModel string `form:"authmodel" json:"authmodel" binding:"required"`
}

type ServerCmdbReq struct {
	Cmdbname string `json:"cmdbname" db:"cmdbname" form:"cmdbname" comment:"CMDB名称" validate:"required"`
	PublicIP string `json:"public_ip" db:"public_ip" form:"public_ip" comment:"IP地址" validate:"required"`
	KeyNames string `json:"keyName" form:"keyName" comment:"key名称" validate:"required"`
	SSHPort  int    `json:"ssh_port" db:"ssh_port" form:"ssh_port" comment:"SSH端口号" validate:"required"`
	Enabled  int8   `json:"enabled" db:"enabled" form:"enabled" validate:"oneof=0 1"`
	Desc     string `json:"desc" form:"desc" comment:"描述"`
	Label    string `json:"label" db:"label" form:"label" comment:"标签"`
	Creator  string `json:"creator" form:"creator" comment:"创建人"`
}
type UpdateServerCmdbReq struct {
	ID       uint   `json:"id" db:"id" form:"id"  validate:"required"`
	Cmdbid   int64  `json:"cmdbid" validate:"required" comment:"CMDB ID"`
	Cmdbname string `json:"cmdbname" db:"cmdbname" form:"cmdbname" comment:"CMDB名称" validate:"required"`
	PublicIP string `json:"public_ip" db:"public_ip" form:"public_ip" comment:"IP地址"`
	Keyid    int64  `json:"keyid" validate:"required"`
	SSHPort  int    `json:"ssh_port" db:"ssh_port" form:"ssh_port" comment:"SSH端口号"`
	Enabled  int8   `json:"enabled" db:"enabled" form:"enabled" validate:"oneof=0 1"`
	Desc     string `json:"desc" form:"desc" comment:"描述"`
	Label    string `json:"label" db:"label" form:"label" comment:"标签"`
	Creator  string `json:"creator" form:"creator" comment:"创建人"`
}

type KeyReq struct {
	KeyName    string `json:"keyName" form:"keyName" comment:"key名称" validate:"required"`
	ServerName string `json:"serverName" form:"serverName" comment:"服务器用户名"`
	Password   string `json:"password" form:"password" comment:"密码"`
	PrivateKey string `json:"privateKey" form:"privateKey" comment:"私钥"`
	Creator    string `json:"creator" form:"creator" comment:"创建人"`
}
type DeleteReq struct {
	Ids []uint `json:"Ids" validate:"required"`
}

type ListReq struct {
	Name    string `json:"name" form:"name" comment:"名称"`
	Creator string `json:"creator" form:"creator" comment:"创建人"`
	Page
}
type InfoReq struct {
	Name string `json:"Name" form:"Name" comment:"名称" validate:"required"`
}
