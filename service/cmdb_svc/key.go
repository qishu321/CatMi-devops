package cmdbsvc

import (
	"CatMi-devops/config"
	"CatMi-devops/initalize/cmdbinit"
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"fmt"

	"github.com/gin-gonic/gin"
)

type KeySvc struct{}

func (s KeySvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.DeleteReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	err := cmdbinit.Keys.Delete(r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除key失败: %s", err.Error()))
	}

	return nil, nil

}
func (s KeySvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)

	r, ok := req.(*request.KeyReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var encryptedPassword, encryptedPrivateKey string
	var err error

	// Encrypt password if provided
	if r.Password != "" {
		encryptedPassword, err = tools.Encrypt(Encryptkey, r.Password)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("Password 加密失败: %s", err.Error()))

		}
	}
	// Encrypt private key if provided
	if r.PrivateKey != "" {
		encryptedPrivateKey, err = tools.Encrypt(Encryptkey, r.PrivateKey)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("PrivateKey 加密失败: %s", err.Error()))
		}
	}
	code := cmdbinit.Keys.Checkcmdb(r.KeyName)
	if code != 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("key名重复，请换新key名"))
	}
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	// Create key object
	key := model.Key{
		KeyName:    r.KeyName,
		ServerName: r.ServerName,
		Password:   encryptedPassword,
		PrivateKey: encryptedPrivateKey,
		Creator:    ctxUser.Username,
	}
	err = cmdbinit.Keys.Add(&key)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建key失败: %s", err.Error()))
	}

	return response.Data{Data: key}, nil
}

func (s KeySvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)

	r, ok := req.(*request.KeyReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var encryptedPassword, encryptedPrivateKey string
	var err error

	// Encrypt password if provided
	if r.Password != "" {
		encryptedPassword, err = tools.Encrypt(Encryptkey, r.Password)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("Password 加密失败: %s", err.Error()))

		}
	}
	// Encrypt private key if provided
	if r.PrivateKey != "" {
		encryptedPrivateKey, err = tools.Encrypt(Encryptkey, r.PrivateKey)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("PrivateKey 加密失败: %s", err.Error()))
		}
	}
	code := cmdbinit.Keys.Checkcmdb(r.KeyName)
	if code == 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("key名不存在，请确认key名"))
	}
	// Create key object
	key := model.Key{
		KeyName:    r.KeyName,
		ServerName: r.ServerName,
		Password:   encryptedPassword,
		PrivateKey: encryptedPrivateKey,
		Creator:    r.Creator,
	}
	err = cmdbinit.Keys.Update(&key)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新key失败: %s", err.Error()))
	}

	return response.Data{Data: key}, nil
}

func (s KeySvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ListReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.Keys.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取key列表失败: %s", err.Error()))
	}
	lists := make([]model.Key, 0)
	for _, key := range list {
		lists = append(lists, *key)
	}
	count, err := cmdbinit.Keys.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取接口总数失败"))
	}
	return response.KeyRsp{
		Total: count,
		Keys:  lists,
	}, nil

}
func (s KeySvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.InfoReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.Keys.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取key失败: %s", err.Error()))
	}
	return list, nil

}
