package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 配置API路由
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	// 设置CORS
	r.Use(corsMiddleware())

	// API版本组
	v1 := r.Group("/api/v1")
	{
		// 电影相关API
		v1.GET("/movies/:id", handler.GetMovie) // 获取单个电影
		v1.GET("/movies", handler.SearchMovies) // 搜索电影
	}

	return r
}

// CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
