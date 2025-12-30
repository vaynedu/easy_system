package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/model"
)

// GetStatistics 获取系统统计信息
func GetStatistics(c *gin.Context) {
	// 获取总题目数
	var totalQuestions int64
	config.DB.Model(&model.ExamQuestion{}).Count(&totalQuestions)

	// 获取各题型数量
	var choiceQuestions int64
	config.DB.Model(&model.ExamQuestion{}).Where("question_type = ?", 0).Count(&choiceQuestions)

	var fillQuestions int64
	config.DB.Model(&model.ExamQuestion{}).Where("question_type = ?", 1).Count(&fillQuestions)

	var essayQuestions int64
	config.DB.Model(&model.ExamQuestion{}).Where("question_type = ?", 2).Count(&essayQuestions)

	// 获取分类统计
	var tags []string
	config.DB.Model(&model.ExamQuestion{}).Distinct().Pluck("tag", &tags)

	tagStats := make(map[string]int)
	for _, tag := range tags {
		if tag != "" {
			var count int64
			config.DB.Model(&model.ExamQuestion{}).Where("tag = ?", tag).Count(&count)
			tagStats[tag] = int(count)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_questions":  totalQuestions,
			"choice_questions": choiceQuestions,
			"fill_questions":   fillQuestions,
			"essay_questions":  essayQuestions,
			"tag_statistics":   tagStats,
		},
	})
}

// GetQuestionTypeCount 按题型获取题目数量
func GetQuestionTypeCount(c *gin.Context) {
	questionTypeStr := c.Query("type")
	questionType, err := strconv.Atoi(questionTypeStr)
	if err != nil {
		questionType = -1 // 表示获取全部
	}

	var count int64
	query := config.DB.Model(&model.ExamQuestion{})

	if questionType >= 0 && questionType <= 2 {
		query = query.Where("question_type = ?", questionType)
	}

	query.Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"count": count,
			"type":  questionType,
		},
	})
}

// GetTagStatistics 获取分类统计
func GetTagStatistics(c *gin.Context) {
	var tags []string
	config.DB.Model(&model.ExamQuestion{}).Distinct().Pluck("tag", &tags)

	tagStats := make(map[string]int)
	for _, tag := range tags {
		if tag != "" {
			var count int64
			config.DB.Model(&model.ExamQuestion{}).Where("tag = ?", tag).Count(&count)
			tagStats[tag] = int(count)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": tagStats,
	})
}
