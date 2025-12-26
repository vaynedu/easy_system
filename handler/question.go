package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/model"
	"github.com/xuri/excelize/v2"
)

// AddQuestion 手动新增题目接口
func AddQuestion(c *gin.Context) {
	// 定义接收参数的结构体
	var req model.ExamQuestion
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数解析失败：" + err.Error(),
		})
		return
	}

	// 1. 题型合法性校验（0/1/2）
	if req.QuestionType != 0 && req.QuestionType != 1 && req.QuestionType != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "题型无效！仅支持0（选择题）、1（填空题）、2（问答题）",
		})
		return
	}

	// 2. 通用校验：题干和正确答案不能为空
	if req.QuestionTitle == "" || req.CorrectAnswer == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "题干和正确答案不能为空！",
		})
		return
	}

	// 3. 选择题专属校验
	if req.QuestionType == 0 {
		if req.OptionA == "" || req.OptionB == "" || req.OptionC == "" || req.OptionD == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "选择题的选项A-D不能为空！",
			})
			return
		}
		// 校验选择题答案
		validAnswers := []string{"A", "B", "C", "D"}
		isValid := false
		for _, a := range validAnswers {
			if strings.ToUpper(req.CorrectAnswer) == a {
				isValid = true
				break
			}
		}
		if !isValid {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "选择题正确答案只能是A/B/C/D！",
			})
			return
		}
	}

	// 4. GORM插入数据，排除时间字段以使用数据库自动时间
	if err := config.DB.Omit("CreatedAt").Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "新增题目失败：" + err.Error(),
		})
		return
	}

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"msg":  fmt.Sprintf("新增题目成功，题目ID：%d", req.ID),
		"code": 200,
	})
}

// ImportExcelQuestion Excel批量导入题目接口
func ImportExcelQuestion(c *gin.Context) {
	// 接收上传的Excel文件
	file, err := c.FormFile("excelFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取Excel文件失败：" + err.Error(),
		})
		return
	}

	// 校验文件格式（仅.xlsx）
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "仅支持.xlsx格式的Excel文件！",
		})
		return
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "打开Excel文件失败：" + err.Error(),
		})
		return
	}
	defer src.Close()

	// 使用打开的文件作为io.Reader来创建Excel文件对象
	excelFile, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "解析Excel文件失败：" + err.Error(),
		})
		return
	}

	// 获取工作表数据
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "读取Excel数据失败：" + err.Error(),
		})
		return
	}

	// 校验数据行数（至少表头+1行数据）
	if len(rows) <= 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Excel文件中无有效题目数据（需包含表头+至少1行题目）",
		})
		return
	}

	// 解析Excel数据
	var questions []model.ExamQuestion
	successCount := 0
	failCount := 0
	invalidRow := 0

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		// 校验列数（至少9列：题型、题干、选项A-D、正确答案、解析、备注）
		if len(row) < 9 {
			failCount++
			invalidRow = i + 1 // Excel行号从1开始
			continue
		}

		// 提取并格式化数据
		typeStr := strings.TrimSpace(row[0])
		title := strings.TrimSpace(row[1])
		optA := strings.TrimSpace(row[2])
		optB := strings.TrimSpace(row[3])
		optC := strings.TrimSpace(row[4])
		optD := strings.TrimSpace(row[5])
		answer := strings.TrimSpace(strings.ToUpper(row[6]))
		analysis := strings.TrimSpace(row[7])
		remark := strings.TrimSpace(row[8])

		// 题型转换与校验
		typeInt, err := strconv.Atoi(typeStr)
		if err != nil || (typeInt != 0 && typeInt != 1 && typeInt != 2) {
			failCount++
			invalidRow = i + 1
			continue
		}

		// 通用校验
		if title == "" || answer == "" {
			failCount++
			invalidRow = i + 1
			continue
		}

		// 选择题专属校验
		if typeInt == 0 {
			if optA == "" || optB == "" || optC == "" || optD == "" {
				failCount++
				invalidRow = i + 1
				continue
			}
			validAnswers := []string{"A", "B", "C", "D"}
			isValid := false
			for _, a := range validAnswers {
				if answer == a {
					isValid = true
					break
				}
			}
			if !isValid {
				failCount++
				invalidRow = i + 1
				continue
			}
		}

		// 构造题目对象
		questions = append(questions, model.ExamQuestion{
			QuestionType:   int8(typeInt),
			QuestionTitle:  title,
			OptionA:        optA,
			OptionB:        optB,
			OptionC:        optC,
			OptionD:        optD,
			CorrectAnswer:  answer,
			AnswerAnalysis: analysis,
			QuestionRemark: remark,
		})
		successCount++
	}

	// GORM批量插入（每100条一批，避免单次插入过多）
	if len(questions) > 0 {
		if err := config.DB.CreateInBatches(&questions, 100).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "批量插入题目失败：" + err.Error(),
			})
			return
		}
	}

	// 返回导入结果
	msg := fmt.Sprintf("导入完成！成功：%d 道，失败：%d 道", successCount, failCount)
	if failCount > 0 {
		msg += fmt.Sprintf("（首个无效行：Excel第 %d 行）", invalidRow)
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  msg,
		"code": 200,
	})
}

// GetRandom10Questions 随机获取10道题接口
func GetRandom10Questions(c *gin.Context) {
	var questions []model.ExamQuestion

	// GORM随机查询10道题（ORDER BY RAND()适配MySQL）
	if err := config.DB.Order("RAND()").Limit(10).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "获取题目失败：" + err.Error(),
		})
		return
	}

	// 无题目时返回提示
	if len(questions) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "题库中暂无题目，请先录入！",
			"data": questions,
		})
		return
	}

	// 返回题目列表
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": questions,
	})
}
