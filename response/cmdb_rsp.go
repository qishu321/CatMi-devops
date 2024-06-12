package response

import "CatMi-devops/model"

type Data struct {
	Data interface{} `json:"data"`
}

type ServerGroupReq struct {
	Total  int64               `json:"total"`
	Groups []model.ServerGroup `json:"serverGroups"`
}

type ServerCmdbRsp struct {
	Total       int64              `json:"total"`
	ServerCmdbs []model.ServerCmdb `json:"serverCmdbs"`
}

type KeyRsp struct {
	Total int64       `json:"total"`
	Keys  []model.Key `json:"keys"`
}
