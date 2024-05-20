package router

import (
	"CatMi-devops/middleware"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"CatMi-devops/controller/system"

)

func InitRoleRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	role := r.Group("/role")
	// 开启jwt认证中间件
	role.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	role.Use(middleware.CasbinMiddleware())
	{
		role.GET("/list", system.Role.List)
		role.POST("/add", system.Role.Add)
		role.POST("/update", system.Role.Update)
		role.POST("/delete", system.Role.Delete)

		role.GET("/getmenulist", system.Role.GetMenuList)
		role.GET("/getapilist", system.Role.GetApiList)
		role.POST("/updatemenus", system.Role.UpdateMenus)
		role.POST("/updateapis", system.Role.UpdateApis)
	}
	return r
}
