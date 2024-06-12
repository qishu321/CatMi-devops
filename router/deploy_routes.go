package router

import (
	"CatMi-devops/controller/deploy"

	"github.com/gin-gonic/gin"
)

func InitDeployRoutes(r *gin.RouterGroup) gin.IRoutes {
	sshd := r.Group("/deploy/sshd")
	//// 开启jwt认证中间件
	//user.Use(authMiddleware.MiddlewareFunc())
	//// 开启casbin鉴权中间件
	//user.Use(middleware.CasbinMiddleware())
	{
		sshd.POST("/command", deploy.Sshd.Command)                     // 添加key
		sshd.POST("/sshdCommandParams", deploy.Sshd.SshdCommandParams) // 添加key
	}
	logs := r.Group("/deploy/logs")
	//// 开启jwt认证中间件
	//user.Use(authMiddleware.MiddlewareFunc())
	//// 开启casbin鉴权中间件
	//user.Use(middleware.CasbinMiddleware())
	{
		logs.POST("/commandParams_list", deploy.CommandLog.List) // 添加key
	}
	task := r.Group("/deploy/task")
	//// 开启jwt认证中间件
	//user.Use(authMiddleware.MiddlewareFunc())
	//// 开启casbin鉴权中间件
	//user.Use(middleware.CasbinMiddleware())
	{
		task.POST("/task_list", deploy.Task.List)     // 添加key
		task.POST("/task_info", deploy.Task.Info)     // 添加key
		task.POST("/task_del", deploy.Task.Delete)    // 添加key
		task.POST("/task_add", deploy.Task.Add)       // 添加key
		task.POST("/task_update", deploy.Task.Update) // 添加key

		task.POST("/task_env_list", deploy.TaskEnv.List)     // 添加key
		task.POST("/task_env_info", deploy.TaskEnv.Info)     // 添加key
		task.POST("/task_env_del", deploy.TaskEnv.Delete)    // 添加key
		task.POST("/task_env_add", deploy.TaskEnv.Add)       // 添加key
		task.POST("/task_env_update", deploy.TaskEnv.Update) // 添加key

	}

	template := r.Group("/deploy/template")
	//// 开启jwt认证中间件
	//user.Use(authMiddleware.MiddlewareFunc())
	//// 开启casbin鉴权中间件
	//user.Use(middleware.CasbinMiddleware())
	{
		template.POST("/template_list", deploy.Template.List)     // 添加key
		template.POST("/template_info", deploy.Template.Info)     // 添加key
		template.POST("/template_del", deploy.Template.Delete)    // 添加key
		template.POST("/template_add", deploy.Template.Add)       // 添加key
		template.POST("/template_update", deploy.Template.Update) // 添加key

	}

	return r
}
