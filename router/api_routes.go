package router

import (
	"CatMi-devops/controller/system"
	"CatMi-devops/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	api := r.Group("/api")
	// 开启jwt认证中间件
	api.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	api.Use(middleware.CasbinMiddleware())
	{
		api.GET("/tree", system.Api.GetTree)
		api.GET("/list", system.Api.List)
		api.POST("/add", system.Api.Add)
		api.POST("/update", system.Api.Update)
		api.POST("/delete", system.Api.Delete)
	}

	return r
}
