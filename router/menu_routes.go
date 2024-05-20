package router

import (
	"CatMi-devops/middleware"
	"CatMi-devops/controller/system"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	menu := r.Group("/menu")
	// 开启jwt认证中间件
	menu.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	menu.Use(middleware.CasbinMiddleware())
	{
		menu.GET("/tree", system.Menu.GetTree)
		menu.GET("/access/tree", system.Menu.GetAccessTree)
		menu.POST("/add", system.Menu.Add)
		menu.POST("/update", system.Menu.Update)
		menu.POST("/delete", system.Menu.Delete)
	}

	return r
}
