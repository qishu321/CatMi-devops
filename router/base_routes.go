package router

import (
	"CatMi-devops/controller/system"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	base := r.Group("/base")
	{
		base.GET("ping", system.Demo)
		base.GET("getpasswd", system.Base.GetPasswd) // 将明文字符串转为MySQL识别的密码
		// 登录登出刷新token无需鉴权
		base.POST("/login", authMiddleware.LoginHandler)
		base.POST("/logout", authMiddleware.LogoutHandler)
		base.POST("/refreshToken", authMiddleware.RefreshHandler)
		base.GET("/dashboard", system.Base.Dashboard) // 系统首页展示数据
	}
	return r
}
