package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/consts"

	"net/http"
)

// GetTagTree 查询完整的Tag知识树（一级+二级分类）
// @Summary 查询Tag分类列表
// @Description 返回所有一级分类及对应的二级分类列表
// @Tags Tag管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{ "code":200, "msg":"查询成功", "data":[]util.PrimaryTag }
// @Router /api/tag/tree [get]
func GetTagTree(c *gin.Context) {
	// 直接返回全局知识树数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询Tag知识树成功",
		"data": consts.KnowledgeTree,
	})
}
