package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/service"
)

// 校验请求参数是否合法
func validateGenerateAIQuestionRequest(req *service.GenerateAIQuestionRequest) error {
	return service.ValidateGenerateAIQuestionRequest(req)
}

// GenerateAIQuestion 生成AI题目
func GenerateAIQuestion(c *gin.Context) {

	ctx := c.Request.Context()

	var req service.GenerateAIQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数解析失败：" + err.Error(),
		})
		return
	}

	// 校验请求参数是否合法
	if err := validateGenerateAIQuestionRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数校验失败：" + err.Error(),
		})
		return
	}

	// 调用Service层生成AI题目
	questions, err := service.GenerateAIQuestionService(ctx, req.QuestionType, req.Tag, req.SecondTag, req.Count, req.Requirements)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成AI题目失败：" + err.Error(),
		})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "AI题目生成成功，但无有效题目",
			"data": questions,
		})
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "AI题目生成成功",
		"data": questions,
	})
}
