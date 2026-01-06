package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/service"
)

// CreateCollection 创建收藏
func CreateCollection(c *gin.Context) {
	// 解析请求参数
	type CreateCollectionRequest struct {
		QuestionID uint   `json:"question_id" binding:"required"`
		Tag        string `json:"tag" binding:"required"`
		SecondTag  string `json:"second_tag" binding:"required"`
	}

	var req CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数解析失败：" + err.Error(),
		})
		return
	}

	// 调用服务创建收藏
	if err := service.CreateCollectionService(req.QuestionID, req.Tag, req.SecondTag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "收藏成功",
	})
}

// DeleteCollection 删除收藏
func DeleteCollection(c *gin.Context) {
	// 解析请求参数
	questionIDStr := c.Query("question_id")
	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "题目ID格式错误",
		})
		return
	}

	// 调用服务删除收藏
	if err := service.DeleteCollectionService(uint(questionID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "取消收藏成功",
	})
}

// GetCollectionStatus 获取收藏状态
func GetCollectionStatus(c *gin.Context) {
	// 解析请求参数
	questionIDStr := c.Query("question_id")
	questionID, err := strconv.ParseUint(questionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "题目ID格式错误",
		})
		return
	}

	// 调用服务获取收藏状态
	isCollected, err := service.GetCollectionStatusService(uint(questionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取收藏状态失败：" + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":        200,
		"msg":         "获取收藏状态成功",
		"is_collected": isCollected,
	})
}

// BatchGetCollectionStatus 批量获取收藏状态
func BatchGetCollectionStatus(c *gin.Context) {
	// 解析请求参数
	questionIDsStr := c.QueryArray("question_ids")
	questionIDs := make([]uint, 0, len(questionIDsStr))
	for _, idStr := range questionIDsStr {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			continue
		}
		questionIDs = append(questionIDs, uint(id))
	}

	// 调用服务批量获取收藏状态
	statusMap, err := service.BatchGetCollectionStatusService(questionIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量获取收藏状态失败：" + err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "批量获取收藏状态成功",
		"data": statusMap,
	})
}

// GetCollectionList 获取收藏列表
func GetCollectionList(c *gin.Context) {
	// 解析请求参数
	tag := c.Query("tag")
	secondTag := c.Query("second_tag")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	// 调用服务获取收藏列表
	collections, total, err := service.GetCollectionListService(tag, secondTag, page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取收藏列表失败：" + err.Error(),
		})
		return
	}

	// 构造返回数据
	var result []map[string]interface{}
	for _, collection := range collections {
		result = append(result, gin.H{
			"id":         collection.ID,
			"question_id": collection.QuestionID,
			"tag":        collection.Tag,
			"second_tag":  collection.SecondTag,
			"created_at":  collection.CreatedAt,
			"question":     collection.Question,
		})
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取收藏列表成功",
		"data": gin.H{
			"collections": result,
			"total":       total,
			"page":        page,
			"size":        size,
		},
	})
}
