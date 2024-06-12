package deploysvc

import (
	"CatMi-devops/config"
	"CatMi-devops/initalize/deployinit"
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SshdSvc struct{}

func (s SshdSvc) Command(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)

	r, ok := req.(*request.SSHClientConfigReq)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c
	var encryptedPassword, encryptedPrivateKey string
	var err error

	// Encrypt password if provided
	if r.Password != "" {
		encryptedPassword, err = tools.Decrypt(Encryptkey, r.Password)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("Password 解密失败: %s", err.Error()))

		}
	}
	// Encrypt private key if provided
	if r.PrivateKey != "" {
		encryptedPrivateKey, err = tools.Decrypt(Encryptkey, r.PrivateKey)
		if err != nil {
			return nil, tools.NewMySqlError(fmt.Errorf("PrivateKey 解密失败: %s", err.Error()))
		}
	}
	config := &request.SSHClientConfigReq{
		Timeout:    time.Second * time.Duration(5+r.TimeoutS),
		PublicIP:   r.PublicIP,
		Port:       r.Port,
		UserName:   r.UserName,
		AuthModel:  r.AuthModel,
		Password:   encryptedPassword,
		PrivateKey: encryptedPrivateKey,
	}
	output, err := tools.SshCommand(config, r.Command)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("ssh执行失败: %s", err.Error()))
	}
	formattedOutput := strings.ReplaceAll(output, "\n", "<br>")

	return formattedOutput, nil
}

func (s SshdSvc) SshdCommandParams(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)
	startTime := time.Now()
	formattedStartTime := startTime.Format("2006-01-02T15:04:05")

	r, ok := req.(*request.SSHReqParams)
	if !ok {
		return nil, ReqdeployErr
	}
	_ = c

	var results []response.SShRsp
	var ReqParams []response.SShParams

	var wg sync.WaitGroup
	errChan := make(chan error, len(r.Data))
	sem := make(chan struct{}, 20) // 限制并发数为10

	for _, param := range r.Data {
		wg.Add(1)
		go func(param request.ServerInfo) {
			defer wg.Done()
			sem <- struct{}{}        // 获取一个令牌
			defer func() { <-sem }() // 释放令牌
			var encryptedPassword, encryptedPrivateKey string
			var err error
			switch param.AuthModel {
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
				Timeout:    time.Second * time.Duration(5+r.TimeoutS),
				PublicIP:   param.PublicIP,
				Port:       param.SSHPort,
				UserName:   param.Keys[0].ServerName,
				AuthModel:  param.AuthModel,
				Password:   encryptedPassword,
				PrivateKey: encryptedPrivateKey,
			}
			output, err := tools.SshCommand(config, r.Command)
			if err != nil {
				errChan <- fmt.Errorf("ssh执行失败: %s", err.Error())
				return
			}
			// 执行成功，添加执行结果到数组中
			results = append(results, response.SShRsp{
				CmdbName: param.Cmdbname,
				Data:     output,
			})
			ReqParams = append(ReqParams, response.SShParams{
				CmdbName: param.Cmdbname,
			})

		}(param)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		// 如果存在错误，则将错误信息添加到执行结果数组中
		if err != nil {
			results = append(results, response.SShRsp{
				Data: fmt.Sprintf("执行失败: %s", err.Error()),
			})
		}
	}
	// 将 r 和 results 转换为 JSON 字符串
	sshReqParamsJSON, err := json.Marshal(ReqParams)
	if err != nil {
		fmt.Errorf("sshReqParamsJSON解析失败: %s", err)
	}

	sshRspJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Errorf("sshRspJSON解析失败: %v", err)
	}
	// 获取当前用户
	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}

	CommandLog := model.CommandLog{
		SSHReqParams: string(sshReqParamsJSON),
		SShRsp:       string(sshRspJSON),
		Ip:           c.ClientIP(),
		Status:       c.Writer.Status(),
		StartTime:    fmt.Sprintf("%v", formattedStartTime),
		Creator:      ctxUser.Username,
		Command:      r.Command,
	}
	err = deployinit.CommandLogs.Add(&CommandLog)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("执行日志写入数据库失败: %s", err.Error()))
	}

	return results, nil
}
