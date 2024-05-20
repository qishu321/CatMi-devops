package router

import (
	"CatMi-devops/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"CatMi-devops/controller/system"

)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	user := r.Group("/user")
	// 开启jwt认证中间件
	user.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	user.Use(middleware.CasbinMiddleware())
	{
		user.GET("/info", system.User.GetUserInfo)                   // 暂时未完成
		user.GET("/list", system.User.List)                          // 用户列表
		user.POST("/add", system.User.Add)                           // 添加用户
		user.POST("/update", system.User.Update)                     // 更新用户
		user.POST("/delete", system.User.Delete)                     // 删除用户
		user.POST("/changePwd", system.User.ChangePwd)               // 修改用户密码
		user.POST("/changeUserStatus", system.User.ChangeUserStatus) // 修改用户状态
	}
	return r
}
