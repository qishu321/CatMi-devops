package router

import (
	"CatMi-devops/controller/cmdb"
	"CatMi-devops/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"
)

func InitCmdbRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	key := r.Group("/cmdb/key")
	//// 开启jwt认证中间件
	key.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	key.Use(middleware.CasbinMiddleware())
	{

		key.POST("/add", cmdb.Key.Add)       // 添加key
		key.POST("/update", cmdb.Key.Update) // 添加key
		key.POST("/info", cmdb.Key.Info)     // 查看指定key
		key.POST("/list", cmdb.Key.List)     // 获取所有key
		key.POST("/delete", cmdb.Key.Delete) // 添加key

	}
	server := r.Group("/cmdb/server")
	// server.Use(authMiddleware.MiddlewareFunc())
	// // 开启casbin鉴权中间件
	// server.Use(middleware.CasbinMiddleware())

	{
		server.POST("/add", cmdb.Server.Add)                // 添加Server
		server.POST("/update", cmdb.Server.Update)          // 添加Server
		server.POST("/info", cmdb.Server.Info)              // 查看指定Server
		server.POST("/list", cmdb.Server.List)              // 获取所有Server
		server.POST("/enable_list", cmdb.Server.EnableList) // // 获取所有有效Server

		server.POST("/delete", cmdb.Server.Delete)               // 添加Server
		server.POST("/enable", cmdb.Server.Enabled)              // 获取Server是否可达
		server.POST("/enabledParams", cmdb.Server.EnabledParams) // 获取Server是否可达

	}
	group := r.Group("/cmdb/groups")
	group.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	group.Use(middleware.CasbinMiddleware())

	{
		group.POST("/add", cmdb.Group.Add)               // 添加Group
		group.POST("/update", cmdb.Group.Update)         // 添加Group
		group.POST("/info", cmdb.Group.Info)             // 查看指定Group
		group.POST("/list", cmdb.Group.List)             // 获取所有Group
		group.POST("/delete", cmdb.Group.Delete)         // 添加Group
		group.POST("/delGroupid", cmdb.Group.DelGroupId) // delGroupid解绑指定服务器

	}

	return r
}
