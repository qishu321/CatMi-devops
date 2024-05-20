package router

import (
	"CatMi-devops/controller/system"
	"CatMi-devops/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func InitGroupRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	group := r.Group("/group")
	// 开启jwt认证中间件
	group.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	group.Use(middleware.CasbinMiddleware())
	{
		group.GET("/list", system.Group.List)
		group.GET("/tree", system.Group.GetTree)
		group.POST("/add", system.Group.Add)
		group.POST("/update", system.Group.Update)
		group.POST("/delete", system.Group.Delete)
		group.POST("/adduser", system.Group.AddUser)
		group.POST("/removeuser", system.Group.RemoveUser)

		group.GET("/useringroup", system.Group.UserInGroup)
		group.GET("/usernoingroup", system.Group.UserNoInGroup)
	}

	return r
}
