package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/handler"
)

// InitRouter 初始化Gin路由
func InitRouter() *gin.Engine {
	r := gin.Default()

	// 配置跨域（解决前端本地访问的跨域问题）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源（开发环境）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API路由分组
	api := r.Group("/api")
	{
		api.POST("/addQuestion", handler.AddQuestion)                 // 新增题目
		api.POST("/importExcelQuestion", handler.ImportExcelQuestion) // Excel导入
		api.GET("/getRandom10", handler.GetRandom10Questions)         // 随机抽10题
	}

	return r
}
