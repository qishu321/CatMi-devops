package response

import "CatMi-devops/model"

type RoleListRsp struct {
	Total int64        `json:"total"`
	Roles []model.Role `json:"roles"`
}