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
		api.POST("/exportExcelQuestion", handler.ExportExcelQuestion) // Excel导出
		api.GET("/getRandom10", handler.GetRandom10Questions)         // 随机抽10题
		api.GET("/tag/tree", handler.GetTagTree)                      // 获取标签树
		// 专项练习相关接口
		r.GET("/api/specialQuestions", handler.GetSpecialQuestionsByFilter)

		// 统计相关接口
		r.GET("/api/statistics", handler.GetStatistics)
		r.GET("/api/questionTypeCount", handler.GetQuestionTypeCount)
		r.GET("/api/tagStatistics", handler.GetTagStatistics)

		// 题库管理相关路由
		api.GET("/questions", handler.GetQuestionsByFilter) // 获取题目列表（带筛选）
		api.GET("/question/:id", handler.GetQuestionByID)   // 获取题目详情
		api.PUT("/question/:id", handler.UpdateQuestion)    // 更新题目
		api.DELETE("/question/:id", handler.DeleteQuestion) // 删除题目
	}

	return r
}
