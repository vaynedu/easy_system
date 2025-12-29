package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/consts"
	"github.com/vaynedu/exam_system/dao"
	"github.com/vaynedu/exam_system/model"
	"github.com/vaynedu/exam_system/service"
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

	// 调用Service层处理业务逻辑
	if err := service.AddQuestionService(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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
		tag := strings.TrimSpace(row[9])        // 一级分类
		secondTag := strings.TrimSpace(row[10]) // 二级分类

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

		// 【核心修改】使用新的Tag校验函数
		if tag != "" {
			if !consts.IsValidPrimaryTag(tag) {
				failCount++
				invalidRow = i + 1
				continue
			}
			if secondTag == "" {
				failCount++
				invalidRow = i + 1
				continue
			}
			if !consts.IsSecondaryOfPrimary(tag, secondTag) {
				failCount++
				invalidRow = i + 1
				continue
			}
		} else {
			if secondTag != "" {
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
			Tag:            tag,
			SecondTag:      secondTag,
		})
		successCount++
	}

	// GORM批量插入（每100条一批，避免单次插入过多）
	if len(questions) > 0 {
		if err := dao.NewQuestionDao(config.DB).CreateQuestionsInBatches(questions, 100); err != nil {
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
	// 获取请求参数中的tag和second_tag
	tag := c.Query("tag")
	secondTag := c.Query("second_tag")

	// 调用Service层获取随机题目
	questions, err := service.GetRandomQuestionsService(tag, secondTag, 10)
	if err != nil {
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
