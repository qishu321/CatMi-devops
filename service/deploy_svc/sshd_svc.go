package deploysvc

import (
	"CatMi-devops/config"
	"CatMi-devops/initalize/cmdbinit"
	"CatMi-devops/initalize/deployinit"
	"CatMi-devops/initalize/system"
	"CatMi-devops/model"
	"CatMi-devops/request"
	"CatMi-devops/response"
	"CatMi-devops/utils/tools"
	"encoding/json"
	"fmt"
	"sort"
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
	sem := make(chan struct{}, 60) // 限制并发数为10

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

// SshFileCommand 在远程服务器上并发执行SSH命令，并记录执行日志

func (s SshdSvc) SshFileCommand2(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)
	startTime := time.Now()
	formattedStartTime := startTime.Format("2006-01-02T15:04:05")

	r, ok := req.(*request.EnableTemplateReq)
	if !ok {
		return nil, ReqdeployErr
	}

	var results []response.SShRsp
	var ReqParams []response.SShParams
	// 根据 Sort 对 Stepnames 排序
	sort.SliceStable(r.Stepnames, func(i, j int) bool {
		return r.Stepnames[i].Sort < r.Stepnames[j].Sort
	})
	var wg sync.WaitGroup
	errChan := make(chan error, len(r.Stepnames))
	sem := make(chan struct{}, 50) // 限制并发数为50

	var envList *model.TaskEnv // 声明变量

	// 处理重要性为1的情况，获取环境变量列表
	if r.Important == 1 {
		var err error
		envList, err = deployinit.TaskEnvs.Info(r.Taskenv)
		if err != nil {
			errChan <- fmt.Errorf("获取环境变量失败: %s", err.Error())
		}
	}

	// 处理重要性为1的情况，获取环境变量列表
	for _, param := range r.Stepnames {
		// 获取任务脚本信息
		tasklist, err := deployinit.Tasks.Info(param.Task)
		if err != nil {
			errChan <- fmt.Errorf("获取脚本失败: %s", err.Error())
			continue
		}

		// 替换脚本中的环境变量
		script := tasklist.Script
		if r.Important == 1 && envList != nil {
			script = strings.Replace(script, "{{options}}", r.Options, -1)
		}

		// 获取服务器组信息
		list, err := cmdbinit.ServerGroups.Info(param.GroupName)
		if err != nil {
			errChan <- fmt.Errorf("获取服务器组失败: %s", err.Error())
			continue
		}

		for _, serverGroup := range list {
			for _, serverCmdb := range serverGroup.ServerCmdbs {
				wg.Add(1)
				go func(serverCmdb *model.ServerCmdb) {
					defer wg.Done()
					sem <- struct{}{}        // 获取一个令牌
					defer func() { <-sem }() // 释放令牌

					var encryptedPassword, encryptedPrivateKey string
					var err error

					switch serverCmdb.AuthModel {
					case "PASSWORD":
						encryptedPassword, err = tools.Decrypt(Encryptkey, serverCmdb.Keys[0].Password)
						if err != nil {
							errChan <- fmt.Errorf("Password 解密失败: %s", err.Error())
							return
						}
					case "PrivateKey":
						encryptedPrivateKey, err = tools.Decrypt(Encryptkey, serverCmdb.Keys[0].PrivateKey)
						if err != nil {
							errChan <- fmt.Errorf("PrivateKey 解密失败: %s", err.Error())
							return
						}
					}

					config := &request.SSHClientConfigReq{
						Timeout:    time.Second * time.Duration(5+param.Timeouts),
						PublicIP:   serverCmdb.PublicIP,
						Port:       serverCmdb.SSHPort,
						UserName:   serverCmdb.Keys[0].ServerName,
						AuthModel:  serverCmdb.AuthModel,
						Password:   encryptedPassword,
						PrivateKey: encryptedPrivateKey,
					}

					output, err := tools.CreateFileOnRemoteServer(config, tasklist.Name, tasklist.Type, script)
					if err != nil {
						errChan <- fmt.Errorf("ssh执行失败: %s", err.Error())
						return
					}

					results = append(results, response.SShRsp{
						CmdbName: serverCmdb.CmdbName,
						Data:     output,
					})
					ReqParams = append(ReqParams, response.SShParams{
						CmdbName: serverCmdb.CmdbName,
					})
				}(serverCmdb)
			}
		}
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			results = append(results, response.SShRsp{
				Data: fmt.Sprintf("执行失败: %s", err.Error()),
			})
		}
	}

	sshReqParamsJSON, err := json.Marshal(ReqParams)
	if err != nil {
		fmt.Errorf("sshReqParamsJSON解析失败: %s", err)
	}

	sshRspJSON, err := json.Marshal(results)
	if err != nil {
		fmt.Errorf("sshRspJSON解析失败: %s", err)
	}

	ctxUser, err := system.User.GetCurrentLoginUser(c)
	if err != nil {
		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
	}

	CommandLog := model.CommandLog{
		SSHReqParams: string(sshReqParamsJSON),
		SShRsp:       string(sshRspJSON),
		Ip:           c.ClientIP(),
		Status:       c.Writer.Status(),
		StartTime:    formattedStartTime,
		Creator:      ctxUser.Username,
		// Command:      r.Command,
	}
	err = deployinit.CommandLogs.Add(&CommandLog)
	if err != nil {
		return nil, tools.NewMySqlError(fmt.Errorf("执行日志写入数据库失败: %s", err))
	}

	return results, nil
}
func (s SshdSvc) SshFileCommand(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)
	startTime := time.Now()
	formattedStartTime := startTime.Format("2006-01-02T15:04:05")

	r, ok := req.(*request.EnableTemplateReq)
	if !ok {
		return nil, ReqdeployErr
	}

	var results []response.TemplateSShRsp
	var ReqParams []response.SShParams
	// 根据 Sort 对 Stepnames 排序
	sort.SliceStable(r.Stepnames, func(i, j int) bool {
		return r.Stepnames[i].Sort < r.Stepnames[j].Sort
	})
	var envList *model.TaskEnv // 声明变量

	// 处理重要性为1的情况，获取环境变量列表
	if r.Important == 1 {
		var err error
		envList, err = deployinit.TaskEnvs.Info(r.Taskenv)
		if err != nil {
			fmt.Errorf("获取环境变量失败: %s", err.Error())
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(r.Stepnames))

	for _, param := range r.Stepnames {
		// 在循环内部创建副本
		paramCopy := param

		// 获取任务脚本信息
		tasklist, err := deployinit.Tasks.Info(paramCopy.Task)
		if err != nil {
			fmt.Errorf("获取脚本失败: %s", err.Error())
			continue
		}

		// 替换脚本中的环境变量
		script := tasklist.Script
		if r.Important == 1 && envList != nil {
			script = strings.Replace(script, "{{options}}", r.Options, -1)
		}

		// 获取服务器组信息
		list, err := cmdbinit.ServerGroups.Info(paramCopy.GroupName)
		if err != nil {
			fmt.Errorf("获取服务器组失败: %s", err.Error())
			continue
		}

		for _, serverGroup := range list {
			sem := make(chan struct{}, 50) // 限制并发数为50
			for _, serverCmdb := range serverGroup.ServerCmdbs {
				wg.Add(1)
				go func(serverCmdb *model.ServerCmdb) {
					defer wg.Done()
					sem <- struct{}{}        // 获取一个令牌
					defer func() { <-sem }() // 释放令牌

					var encryptedPassword, encryptedPrivateKey string
					var err error

					switch serverCmdb.AuthModel {
					case "PASSWORD":
						encryptedPassword, err = tools.Decrypt(Encryptkey, serverCmdb.Keys[0].Password)
						if err != nil {
							errChan <- fmt.Errorf("Password 解密失败: %s", err.Error())
							return
						}
					case "PrivateKey":
						encryptedPrivateKey, err = tools.Decrypt(Encryptkey, serverCmdb.Keys[0].PrivateKey)
						if err != nil {
							errChan <- fmt.Errorf("PrivateKey 解密失败: %s", err.Error())
							return
						}
					}

					config := &request.SSHClientConfigReq{
						Timeout:    time.Second * time.Duration(5+paramCopy.Timeouts),
						PublicIP:   serverCmdb.PublicIP,
						Port:       serverCmdb.SSHPort,
						UserName:   serverCmdb.Keys[0].ServerName,
						AuthModel:  serverCmdb.AuthModel,
						Password:   encryptedPassword,
						PrivateKey: encryptedPrivateKey,
					}

					output, err := tools.CreateFileOnRemoteServer(config, tasklist.Name, tasklist.Type, script)
					if err != nil {
						errChan <- fmt.Errorf("ssh执行失败: %s", err.Error())
						return
					}

					// 查找是否已经存在相同 Sort 的 TemplateSShRsp
					found := false
					for i := range results {
						if results[i].Stepname == paramCopy.Stepname {
							results[i].SShRsp = append(results[i].SShRsp, response.SShRsp{
								CmdbName: serverCmdb.CmdbName,
								Data:     output,
							})
							found = true
							break
						}
					}

					// 如果不存在，则创建一个新的 TemplateSShRsp
					if !found {
						results = append(results, response.TemplateSShRsp{
							Stepname: paramCopy.Stepname,
							SShRsp: []response.SShRsp{
								{
									CmdbName: serverCmdb.CmdbName,
									Data:     output,
								},
							},
						})
					}

					ReqParams = append(ReqParams, response.SShParams{
						CmdbName: serverCmdb.CmdbName,
					})
					// 获取当前登录用户信息
					ctxUser, err := system.User.GetCurrentLoginUser(c)
					if err != nil {
						fmt.Errorf("获取当前登陆用户信息失败: %v", err)
					}

					// 只记录当前步骤的结果
					stepResults := []response.TemplateSShRsp{
						{
							Stepname: paramCopy.Stepname,
							SShRsp: []response.SShRsp{
								{
									CmdbName: serverCmdb.CmdbName,
									Data:     output,
								},
							},
						},
					}

					sshReqParamsJSON, err := json.Marshal(ReqParams)
					if err != nil {
						fmt.Errorf("sshReqParamsJSON解析失败: %s", err)
					}

					sshRspJSON, err := json.Marshal(stepResults)
					if err != nil {
						fmt.Errorf("sshRspJSON解析失败: %s", err)
					}

					// 构建命令日志对象
					TemplateLog := model.Template_Log{
						TemplateID:   r.TemplateID,
						Name:         r.Name + "-" + formattedStartTime,
						Stepname:     paramCopy.Stepname,
						Sort:         paramCopy.Sort,
						Command:      script,
						Cmdbnames:    paramCopy.GroupName,
						Timeouts:     paramCopy.Timeouts,
						SshRsp:       string(sshRspJSON),
						SSHReqParams: string(sshReqParamsJSON),
						StartTime:    formattedStartTime,
						Creator:      ctxUser.Username,
					}

					// 将命令日志写入数据库
					err = deployinit.Template_Logs.Add(&TemplateLog)
					if err != nil {
						fmt.Errorf("执行日志写入数据库失败: %s", err)
					}

				}(serverCmdb)
			}
		}

		// 等待当前排序步骤的所有任务完成
		wg.Wait()

		// 在每个排序步骤之间等待 1 秒
		time.Sleep(100)
	}

	close(errChan)

	// 处理错误通道
	for err := range errChan {
		if err != nil {
			results = append(results, response.TemplateSShRsp{
				SShRsp: []response.SShRsp{
					{
						Data: fmt.Sprintf("执行失败: %s", err.Error()),
					},
				},
			})
		}
	}

	return results, nil
}

// func (s SshdSvc) SshFileCommand(c *gin.Context, req interface{}) (data interface{}, rspError interface{}) {
// 	var Encryptkey = []byte(config.Conf.Cmdb.Encryptkey)
// 	startTime := time.Now()
// 	formattedStartTime := startTime.Format("2006-01-02T15:04:05")

// 	r, ok := req.(*request.EnableTemplateReq)
// 	if !ok {
// 		return nil, ReqdeployErr
// 	}

// 	var results []response.TemplateSShRsp
// 	var ReqParams []response.SShParams
// 	// 根据 Sort 对 Stepnames 排序
// 	sort.SliceStable(r.Stepnames, func(i, j int) bool {
// 		return r.Stepnames[i].Sort < r.Stepnames[j].Sort
// 	})
// 	var envList *model.TaskEnv // 声明变量

// 	// 处理重要性为1的情况，获取环境变量列表
// 	if r.Important == 1 {
// 		var err error
// 		envList, err = deployinit.TaskEnvs.Info(r.Taskenv)
// 		if err != nil {
// 			fmt.Errorf("获取环境变量失败: %s", err.Error())
// 		}
// 	}

// 	var wg sync.WaitGroup
// 	errChan := make(chan error, len(r.Stepnames))

// 	for _, param := range r.Stepnames {
// 		// 在循环内部创建副本
// 		paramCopy := param

// 		// 获取任务脚本信息
// 		tasklist, err := deployinit.Tasks.Info(param.Task)
// 		if err != nil {
// 			fmt.Errorf("获取脚本失败: %s", err.Error())
// 			continue
// 		}

// 		// 替换脚本中的环境变量
// 		script := tasklist.Script
// 		if r.Important == 1 && envList != nil {
// 			script = strings.Replace(script, "{{options}}", r.Options, -1)
// 		}

// 		// 获取服务器组信息
// 		list, err := cmdbinit.ServerGroups.Info(param.GroupName)
// 		if err != nil {
// 			fmt.Errorf("获取服务器组失败: %s", err.Error())
// 			continue
// 		}

// 		for _, serverGroup := range list {
// 			sem := make(chan struct{}, 50) // 限制并发数为50
// 			for _, serverCmdb := range serverGroup.ServerCmdbs {
// 				wg.Add(1)
// 				go func(serverCmdb *model.ServerCmdb) {
// 					defer wg.Done()
// 					sem <- struct{}{}        // 获取一个令牌
// 					defer func() { <-sem }() // 释放令牌

// 					var encryptedPassword, encryptedPrivateKey string
// 					var err error

// 					switch serverCmdb.AuthModel {
// 					case "PASSWORD":
// 						encryptedPassword, err = tools.Decrypt(Encryptkey, serverCmdb.Keys[0].Password)
// 						if err != nil {
// 							errChan <- fmt.Errorf("Password 解密失败: %s", err.Error())
// 							return
// 						}
// 					case "PrivateKey":
// 						encryptedPrivateKey, err = tools.Decrypt(Encryptkey, serverCmdb.Keys[0].PrivateKey)
// 						if err != nil {
// 							errChan <- fmt.Errorf("PrivateKey 解密失败: %s", err.Error())
// 							return
// 						}
// 					}

// 					config := &request.SSHClientConfigReq{
// 						Timeout:    time.Second * time.Duration(5+param.Timeouts),
// 						PublicIP:   serverCmdb.PublicIP,
// 						Port:       serverCmdb.SSHPort,
// 						UserName:   serverCmdb.Keys[0].ServerName,
// 						AuthModel:  serverCmdb.AuthModel,
// 						Password:   encryptedPassword,
// 						PrivateKey: encryptedPrivateKey,
// 					}

// 					output, err := tools.CreateFileOnRemoteServer(config, tasklist.Name, tasklist.Type, script)
// 					if err != nil {
// 						errChan <- fmt.Errorf("ssh执行失败: %s", err.Error())
// 						return
// 					}

// 					// 将结果和参数添加到相应的切片中
// 					results = append(results, response.TemplateSShRsp{
// 						Sort: paramCopy.Sort,
// 						SShRsp: response.SShRsp{
// 							CmdbName: serverCmdb.CmdbName,
// 							Data:     output,
// 						},
// 					})
// 					ReqParams = append(ReqParams, response.SShParams{
// 						CmdbName: serverCmdb.CmdbName,
// 					})
// 				}(serverCmdb)
// 			}
// 		}
// 	}

// 	// 等待所有任务完成
// 	wg.Wait()
// 	close(errChan)

// 	// 处理错误通道
// 	for err := range errChan {
// 		if err != nil {
// 			results = append(results, response.TemplateSShRsp{
// 				SShRsp: response.SShRsp{
// 					Data: fmt.Sprintf("执行失败: %s", err.Error()),
// 				},
// 			})
// 		}
// 	}

// 	sshReqParamsJSON, err := json.Marshal(ReqParams)
// 	if err != nil {
// 		fmt.Errorf("sshReqParamsJSON解析失败: %s", err)
// 	}

// 	sshRspJSON, err := json.Marshal(results)
// 	if err != nil {
// 		fmt.Errorf("sshRspJSON解析失败: %s", err)
// 	}

// 	// 获取当前登录用户信息
// 	ctxUser, err := system.User.GetCurrentLoginUser(c)
// 	if err != nil {
// 		fmt.Errorf("获取当前登陆用户信息失败: %v", err)
// 	}

// 	// 构建命令日志对象
// 	CommandLog := model.CommandLog{
// 		SSHReqParams: string(sshReqParamsJSON),
// 		SShRsp:       string(sshRspJSON),
// 		Ip:           c.ClientIP(),
// 		Status:       c.Writer.Status(),
// 		StartTime:    formattedStartTime,
// 		Creator:      ctxUser.Username,
// 		// Command:      r.Command,
// 	}

// 	// 将命令日志写入数据库
// 	err = deployinit.CommandLogs.Add(&CommandLog)
// 	if err != nil {
// 		return nil, tools.NewMySqlError(fmt.Errorf("执行日志写入数据库失败: %s", err))
// 	}

// 	return results, nil
// }
