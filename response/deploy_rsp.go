package response

import "CatMi-devops/model"

type SShRsp struct {
	CmdbName string      `json:"cmdbName"` // 服务器名称
	Data     interface{} `json:"data"`
}
type TemplateSShRsp struct {
	Stepname string   `json:"stepname"`
	SShRsp   []SShRsp `json:"sshrsp"`
}

type SShParams struct {
	CmdbName string `json:"CmdbName"`
}

type TemplateRsp struct {
	Total     int64            `json:"total"`
	Templates []model.Template `json:"Templates"`
}

type Template_logRsp struct {
	Total         int64                `json:"total"`
	Templates_log []model.Template_Log `json:"Templates_log"`
}

type TaskRsp struct {
	Total int64        `json:"total"`
	Tasks []model.Task `json:"Task"`
}

type TaskEnvRsp struct {
	Total    int64           `json:"total"`
	TaskEnvs []model.TaskEnv `json:"TaskEnv"`
}
