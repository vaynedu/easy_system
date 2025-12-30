package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/service"
)

// GetSpecialQuestionsByFilter 根据专项分类筛选条件获取题目列表
func GetSpecialQuestionsByFilter(c *gin.Context) {
	// 获取查询参数
	tag := c.Query("tag")
	secondTag := c.Query("second_tag")
	questionType := c.Query("type")
	keyword := c.Query("keyword")

	// 分页参数
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}
	size, _ := strconv.Atoi(c.Query("size"))
	if size <= 0 || size > 100 {
		size = 10
	}

	// 调用Service层获取题目列表
	questions, total, err := service.GetQuestionsByFilterService(tag, secondTag, questionType, keyword, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "获取专项题目列表失败：" + err.Error(),
			"code": 500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"questions": questions,
			"total":     total,
			"page":      page,
			"size":      size,
		},
	})
}
