package main

import (
	"CatMi-devops/config"
	"CatMi-devops/initalize/system"
	"CatMi-devops/middleware"
	"CatMi-devops/router"
	"CatMi-devops/utils/common"
	"CatMi-devops/utils/tools"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// 加载配置文件到全局配置结构体
	config.InitConfig()
	// 初始化日志
	log.Print("user:", tools.NewGenPasswd("123456"))

	common.InitLogger()

	// 初始化DB
	common.DBS()

	// 初始化数据库(-init_db)
	initDB := flag.Bool("init_db", false, "initialize database")
	flag.Parse()
	if *initDB { //判断是否需要更新表结构
		common.InitDB()
		common.InitCasbinEnforcer()
		// 初始化mysql数据
		common.InitData()
		return
	}

	// 初始化casbin策略管理器
	common.InitCasbinEnforcer()

	// 初始化Validator数据校验
	common.InitValidate()

	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中
	// 这里开启3个goroutine处理channel将日志记录到数据库
	for i := 0; i < 3; i++ {
		go system.OperationLog.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	// 注册所有路由
	r := router.InitRoutes()

	host := "0.0.0.0"
	port := config.Conf.System.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("listen: %s\n", err)
		}
	}()

	common.Log.Info(fmt.Sprintf("Server is running at %s:%d/%s", host, port, config.Conf.System.UrlPathPrefix))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(quit, os.Interrupt)
	<-quit
	common.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Fatal("Server forced to shutdown:", err)
	}

	common.Log.Info("Server exiting!")

}
