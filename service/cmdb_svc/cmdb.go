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
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CmdbSvc struct{}

func (s CmdbSvc) Delete(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.DeleteReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	err := cmdbinit.ServerCmdbs.Delete(r.Ids)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("删除服务器失败: %s", err.Error()))
	}

	return nil, nil

}
func (s CmdbSvc) EnabledParams(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)

	r, ok := req.(*request.EnabledCmdbReqParams)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c

	var wg sync.WaitGroup
	errChan := make(chan error, len(r.Data))
	sem := make(chan struct{}, 10) // 限制并发数为10

	for _, param := range r.Data {
		wg.Add(1)
		go func(param struct {
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
		}) {
			defer wg.Done()
			sem <- struct{}{}        // 获取一个令牌
			defer func() { <-sem }() // 释放令牌

			// list, err := cmdbinit.Keys.Infos(param.KeyNames)
			// if err != nil {
			// 	errChan <- fmt.Errorf("获取key失败: %s", err.Error())
			// 	return
			// }

			var encryptedPassword, encryptedPrivateKey string
			var err error
			switch r.AuthModel {
			case "PASSWORD":
				encryptedPassword, err = tools.Decrypt(Encryptkey, param.Keys[0].Password)
				if err != nil {
					errChan <- fmt.Errorf("Password 解密失败: %s", err.Error())
					return
				}
			case "PrivateKey":
				encryptedPrivateKey, err = tools.Decrypt(Encryptkey, param.Keys[0].PrivateKey)
				if err != nil {
					errChan <- fmt.Errorf("PrivateKey 解密失败: %s", err.Error())
					return
				}
			}

			config := &request.SSHClientConfigReq{
				Timeout:    time.Second * 5,
				PublicIP:   param.PublicIP,
				Port:       param.SSHPort,
				UserName:   param.Keys[0].ServerName,
				AuthModel:  r.AuthModel,
				Password:   encryptedPassword,
				PrivateKey: encryptedPrivateKey,
			}
			_, err = tools.SshCommand(config, "hostname")
			if err != nil {
				errChan <- fmt.Errorf("ssh执行失败: %s", err.Error())
				return
			}

			ServerCmdb := model.ServerCmdb{
				Model:     gorm.Model{ID: param.ID},
				CmdbID:    param.Cmdbid,
				PublicIP:  param.PublicIP,
				SSHPort:   param.SSHPort,
				Enabled:   1,
				AuthModel: r.AuthModel,
			}
			err = cmdbinit.ServerCmdbs.EnabledUpdate(&ServerCmdb)
			if err != nil {
				errChan <- fmt.Errorf("更新服务器失败: %s", err.Error())
				return
			}
		}(param)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return nil, tools.NewMySqlError(err)
		}
	}

	return ok, nil
}

func (s CmdbSvc) Enabled(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var encryptedPassword, encryptedPrivateKey string
	var err error
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)

	r, ok := req.(*request.EnabledServerCmdbReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.Keys.Infos(r.KeyNames)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取key失败: %s", err.Error()))
	}

	switch r.AuthModel {
	case "PASSWORD":
		encryptedPassword, err = tools.Decrypt(Encryptkey, list.Password)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("Password 解密失败: %s", err.Error()))
		}
	case "PrivateKey":
		encryptedPrivateKey, err = tools.Decrypt(Encryptkey, list.PrivateKey)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("PrivateKey 解密失败: %s", err.Error()))
		}
	}

	config := &request.SSHClientConfigReq{
		Timeout:    time.Second * 5,
		PublicIP:   r.PublicIP,
		Port:       r.SSHPort,
		UserName:   list.ServerName,
		AuthModel:  r.AuthModel,
		Password:   encryptedPassword,
		PrivateKey: encryptedPrivateKey,
	}
	_, err = tools.SshCommand(config, "hostname")
	if err != nil {
		ServerCmdb := model.ServerCmdb{
			Model:     gorm.Model{ID: r.ID},
			CmdbID:    r.Cmdbid,
			PublicIP:  r.PublicIP,
			SSHPort:   r.SSHPort,
			Enabled:   0,
			AuthModel: r.AuthModel,
		}
		err = cmdbinit.ServerCmdbs.EnabledUpdate(&ServerCmdb)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("更新服务器失败: %s", err.Error()))
		}
		return nil, tools.NewMySqlError(fmt.Errorf("ssh执行失败: %s", err.Error()))
	}

	ServerCmdb := model.ServerCmdb{
		Model:     gorm.Model{ID: r.ID},
		CmdbID:    r.Cmdbid,
		PublicIP:  r.PublicIP,
		SSHPort:   r.SSHPort,
		Enabled:   1,
		AuthModel: r.AuthModel,
	}
	err = cmdbinit.ServerCmdbs.EnabledUpdate(&ServerCmdb)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新服务器失败: %s", err.Error()))
	}

	return ok, nil
}

func (s CmdbSvc) Add(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.ServerCmdbReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var err error

	code := cmdbinit.ServerCmdbs.Check(r.Cmdbname)
	if code != 200 {
		return nil, tools.NewMySqlError(fmt.Errorf("服务器名称重复，请换新服务器名称"))
	}
	// Create key object
	list, err := cmdbinit.Keys.Info(r.KeyNames)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取key失败: %s", err.Error()))
	}
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取当前登陆用户信息失败"))
	}

	ServerCmdb := model.ServerCmdb{
		CmdbID:   tools.GenerateRandomNumber(),
		CmdbName: r.Cmdbname,
		PublicIP: r.PublicIP,
		SSHPort:  r.SSHPort,
		Enabled:  r.Enabled,
		Keys:     list,
		Creator:  ctxUser.Username,
		Desc:     r.Desc,
		Label:    r.Label,
	}
	err = cmdbinit.ServerCmdbs.Add(&ServerCmdb)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("创建服务器失败: %s", err.Error()))
	}

	return ServerCmdb, nil
}

func (s CmdbSvc) Update(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {

	r, ok := req.(*request.UpdateServerCmdbReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	var err error

	ServerCmdb := model.ServerCmdb{
		Model:    gorm.Model{ID: r.ID},
		CmdbID:   r.Cmdbid,
		PublicIP: r.PublicIP,
		SSHPort:  r.SSHPort,
		Enabled:  r.Enabled,
		Desc:     r.Desc,
		Label:    r.Label,
	}
	err = cmdbinit.ServerCmdbs.Update(&ServerCmdb, r.Keyid)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("更新服务器失败: %s", err.Error()))
	}

	return nil, nil
}

func (s CmdbSvc) List(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ListReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.ServerCmdbs.List(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器列表失败: %s", err.Error()))
	}

	lists := make([]model.ServerCmdb, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}
	count, err := cmdbinit.ServerCmdbs.Count()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器总数失败"))
	}
	return response.ServerCmdbRsp{
		Total:       count,
		ServerCmdbs: lists,
	}, nil

}
func (s CmdbSvc) EnableList(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.ListReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c

	list, err := cmdbinit.ServerCmdbs.EnableList(r)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器失败: %s", err.Error()))
	}
	lists := make([]model.ServerCmdb, 0)
	for _, server := range list {
		lists = append(lists, *server)
	}

	count, err := cmdbinit.ServerCmdbs.EnableCount()
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器总数失败"))
	}

	return response.ServerCmdbRsp{
		Total:       count,
		ServerCmdbs: lists,
	}, nil

}

func (s CmdbSvc) Info(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	r, ok := req.(*request.InfoReq)
	if !ok {
		return nil, ReqCmdbErr
	}
	_ = c
	list, err := cmdbinit.ServerCmdbs.Info(r.Name)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("获取服务器失败: %s", err.Error()))
	}
	return list, nil

}
