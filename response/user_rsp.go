package response

import "CatMi-devops/model"

type UserListRsp struct {
	Total int          `json:"total"`
	Users []model.User `json:"users"`
}
