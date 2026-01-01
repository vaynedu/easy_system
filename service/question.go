package service

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/consts"
	"github.com/vaynedu/exam_system/dao"
	"github.com/vaynedu/exam_system/model"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// AddQuestionService 新增题目服务
func AddQuestionService(question *model.ExamQuestion) error {
	// 1. 题型合法性校验（0/1/2）
	if question.QuestionType != 0 && question.QuestionType != 1 && question.QuestionType != 2 {
		return errors.New("题型无效！仅支持0（选择题）、1（填空题）、2（问答题）")
	}

	// 2. 通用校验：题干和正确答案不能为空
	if question.QuestionTitle == "" || question.CorrectAnswer == "" {
		return errors.New("题干和正确答案不能为空！")
	}

	// 3. 选择题专属校验
	if question.QuestionType == 0 {
		if question.OptionA == "" || question.OptionB == "" || question.OptionC == "" || question.OptionD == "" {
			return errors.New("选择题的选项A-D不能为空！")
		}
		// 校验选择题答案
		validAnswers := []string{"A", "B", "C", "D"}
		isValid := false
		for _, a := range validAnswers {
			if strings.ToUpper(question.CorrectAnswer) == a {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("选择题正确答案只能是A/B/C/D！")
		}
	}

	// 【核心修改】使用新的Tag校验函数（适配知识树）
	if question.Tag != "" {
		// 1. 校验一级分类是否合法
		if !consts.IsValidPrimaryTag(question.Tag) {
			return errors.New("一级分类无效！请从合法分类中选择（算法、系统设计、数据存储、高频考点）")
		}
		// 2. 一级分类存在时，二级分类不能为空
		if question.SecondTag == "" {
			return errors.New("填写一级分类后，必须填写对应的二级分类！")
		}
		// 3. 校验二级分类是否属于该一级分类
		if !consts.IsSecondaryOfPrimary(question.Tag, question.SecondTag) {
			return errors.New(fmt.Sprintf("二级分类「%s」不属于一级分类「%s」！", question.SecondTag, question.Tag))
		}
	} else {
		// 4. 一级分类为空时，二级分类必须为空
		if question.SecondTag != "" {
			return errors.New("未填写一级分类时，禁止填写二级分类！")
		}
	}

	// 4.题目上传方式
	question.UploadType = consts.QuestionImportTypeManual

	// 5. 调用DAO层插入数据
	return dao.NewQuestionDao(config.DB).CreateQuestion(question)
}

// GetRandomQuestionsService 随机获取题目服务
func GetRandomQuestionsService(tag, secondTag string, limit int) ([]model.ExamQuestion, error) {
	// 校验标签参数
	if tag != "" {
		if !IsValidPrimaryTag(tag) {
			return nil, errors.New("一级分类无效")
		}
		if secondTag == "" {
			return nil, errors.New("当指定一级分类时，二级分类不能为空")
		}
		if !IsSecondaryOfPrimary(tag, secondTag) {
			return nil, errors.New("二级分类与一级分类不匹配")
		}
	} else {
		if secondTag != "" {
			return nil, errors.New("未指定一级分类时，不能单独指定二级分类")
		}
	}

	// 调用DAO层获取随机题目
	return dao.NewQuestionDao(config.DB).GetRandomQuestionsByTag(tag, secondTag, limit)
}

// IsValidPrimaryTag 验证一级标签是否有效
func IsValidPrimaryTag(tag string) bool {
	// return consts.IsValidPrimaryTag(tag)
	return consts.IsValidPrimaryTag(tag)
}

// IsSecondaryOfPrimary 验证二级标签是否属于一级标签
func IsSecondaryOfPrimary(primary, secondary string) bool {
	return consts.IsSecondaryOfPrimary(primary, secondary)
}

// GetQuestionsByFilterService 根据筛选条件获取题目列表服务
func GetQuestionsByFilterService(tag, secondTag, questionType, keyword string, page, size int) ([]model.ExamQuestion, int64, error) {
	offset := (page - 1) * size

	// 构建查询条件
	query := config.DB.Model(&model.ExamQuestion{})

	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	if secondTag != "" {
		query = query.Where("second_tag = ?", secondTag)
	}
	if questionType != "" {
		if typeInt, err := strconv.Atoi(questionType); err == nil {
			query = query.Where("question_type = ?", typeInt)
		}
	}
	if keyword != "" {
		query = query.Where("question_title LIKE ? OR correct_answer LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	var questions []model.ExamQuestion
	if err := query.Offset(offset).Limit(size).Order("id DESC").Find(&questions).Error; err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// GetQuestionByIDService 根据ID获取题目详情服务
func GetQuestionByIDService(id uint) (*model.ExamQuestion, error) {
	var question model.ExamQuestion
	if err := config.DB.Where("id = ?", id).First(&question).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &model.ExamQuestion{}, nil // 返回空对象表示未找到
		}
		return nil, err
	}
	return &question, nil
}

// UpdateQuestionService 更新题目服务
func UpdateQuestionService(question *model.ExamQuestion) error {
	// 题目校验（复用AddQuestionService的校验逻辑）
	if question.QuestionType != 0 && question.QuestionType != 1 && question.QuestionType != 2 {
		return errors.New("题型无效！仅支持0（选择题）、1（填空题）、2（问答题）")
	}

	if question.QuestionTitle == "" || question.CorrectAnswer == "" {
		return errors.New("题干和正确答案不能为空！")
	}

	if question.QuestionType == 0 {
		if question.OptionA == "" || question.OptionB == "" || question.OptionC == "" || question.OptionD == "" {
			return errors.New("选择题的选项A-D不能为空！")
		}
		validAnswers := []string{"A", "B", "C", "D"}
		isValid := false
		for _, a := range validAnswers {
			if strings.ToUpper(question.CorrectAnswer) == a {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("选择题正确答案只能是A/B/C/D！")
		}
	}

	// 标签校验
	if question.Tag != "" {
		if !IsValidPrimaryTag(question.Tag) {
			return errors.New("一级分类无效")
		}
		if question.SecondTag == "" {
			return errors.New("当设置一级分类时，二级分类不能为空")
		}
		if !IsSecondaryOfPrimary(question.Tag, question.SecondTag) {
			return errors.New("二级分类与一级分类不匹配")
		}
	} else {
		if question.SecondTag != "" {
			return errors.New("未设置一级分类时，不能单独设置二级分类")
		}
	}

	// 调用DAO层更新
	return dao.NewQuestionDao(config.DB).UpdateQuestion(question)
}

// DeleteQuestionService 删除题目服务
func DeleteQuestionService(id uint) error {
	return dao.NewQuestionDao(config.DB).DeleteQuestion(id)
}

// ImportExcelQuestions 解析Excel并导入题目（核心业务逻辑）
func ImportExcelQuestions(fileReader io.Reader) (successCount, failCount, invalidRow int, err error) {
	// 1. 解析Excel文件
	excelFile, err := excelize.OpenReader(fileReader)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("解析Excel失败：%w", err)
	}

	// 2. 读取工作表数据
	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("读取Excel数据失败：%w", err)
	}

	// 3. 校验数据行数（至少表头+1行数据）
	if len(rows) <= 1 {
		return 0, 0, 0, errors.New("excel无有效数据（需包含表头+至少1行题目）")
	}

	// 4. 解析&校验每行数据
	var questions []model.ExamQuestion
	invalidRow = 0
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		// 单行列数校验
		if len(row) < 9 {
			failCount++
			invalidRow = getExcelRowNum(i)
			continue
		}

		// 提取并格式化数据
		question, rowErr := parseAndValidateRow(row, i)
		if rowErr != nil {
			failCount++
			invalidRow = getExcelRowNum(i)
			continue
		}

		questions = append(questions, *question)
		successCount++
	}

	// 5. 批量插入数据库（调用DAO层）
	if len(questions) > 0 {
		questionDao := dao.NewQuestionDao(config.DB)
		if err := questionDao.CreateQuestionsInBatches(questions, 100); err != nil {
			return successCount, failCount, invalidRow, fmt.Errorf("批量插入失败：%w", err)
		}
	}

	return successCount, failCount, invalidRow, nil
}

// parseAndValidateRow 解析并校验单行数据
func parseAndValidateRow(row []string, rowIdx int) (*model.ExamQuestion, error) {
	// 提取字段（trim空格）
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

	// 1. 题型转换&校验
	typeInt, err := strconv.Atoi(typeStr)
	if err != nil || (typeInt != 0 && typeInt != 1 && typeInt != 2) {
		return nil, errors.New("题型无效（仅支持0/1/2：选择/填空/问答）")
	}

	// 2. 通用必填项校验
	if title == "" || answer == "" {
		return nil, errors.New("题干/正确答案不能为空")
	}

	// 3. 选择题专属校验
	if typeInt == 0 {
		if optA == "" || optB == "" || optC == "" || optD == "" {
			return nil, errors.New("选择题选项A-D不能为空")
		}
		validAnswers := []string{"A", "B", "C", "D"}
		isValidAnswer := false
		for _, a := range validAnswers {
			if answer == a {
				isValidAnswer = true
				break
			}
		}
		if !isValidAnswer {
			return nil, errors.New("选择题答案仅支持A/B/C/D")
		}
	}

	// 4. 标签校验（调用consts层）
	if err := validateTagRelation(tag, secondTag); err != nil {
		return nil, err
	}

	// 构造题目对象
	return &model.ExamQuestion{
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
		UploadType:     consts.QuestionImportTypeExcel,
	}, nil
}

// 辅助函数：获取Excel实际行号（索引+1）
func getExcelRowNum(idx int) int {
	return idx + 1
}

// validateTagRelation 标签关联校验（抽离复用）
func validateTagRelation(primaryTag, secondaryTag string) error {
	if primaryTag != "" {
		// 校验一级标签有效性
		if !consts.IsValidPrimaryTag(primaryTag) {
			return errors.New("一级分类标签无效")
		}
		// 一级标签存在时，二级标签不能为空
		if secondaryTag == "" {
			return errors.New("一级分类存在时，二级分类不能为空")
		}
		// 校验二级标签归属
		if !consts.IsSecondaryOfPrimary(primaryTag, secondaryTag) {
			return errors.New("二级分类不属于当前一级分类")
		}
	} else {
		// 一级标签为空时，二级标签也必须为空
		if secondaryTag != "" {
			return errors.New("一级分类为空时，二级分类不能为空")
		}
	}
	return nil
}

// ExportExcelQuestionRequest 导出题目请求参数结构体
type ExportExcelQuestionRequest struct {
	IDs          []uint `json:"ids"`        // 指定题目ID列表
	ExportAll    bool   `json:"export_all"` // 是否导出全部
	Tag          string `json:"tag"`        // 一级分类
	SecondTag    string `json:"second_tag"` // 二级分类
	QuestionType string `json:"type"`       // 题型
	Keyword      string `json:"keyword"`    // 关键词搜索
}

// ExportExcelQuestionService 导出Excel题目的服务函数
func ExportExcelQuestionService(req ExportExcelQuestionRequest) ([]model.ExamQuestion, error) {
	var questions []model.ExamQuestion
	var err error

	// 根据不同条件获取题目
	if len(req.IDs) > 0 {
		// 根据ID列表导出
		questions, err = dao.NewQuestionDao(config.DB).GetQuestionsByIDList(req.IDs)
		if err != nil {
			return nil, fmt.Errorf("根据ID列表获取题目失败：%v", err)
		}
	} else if req.ExportAll {
		// 导出全部题目
		questions, err = dao.NewQuestionDao(config.DB).GetAllQuestions()
		if err != nil {
			return nil, fmt.Errorf("获取全部题目失败：%v", err)
		}
	} else {
		// 根据筛选条件导出
		questions, err = dao.NewQuestionDao(config.DB).GetQuestionsByFilter(req.Tag, req.SecondTag, req.QuestionType, req.Keyword)
		if err != nil {
			return nil, fmt.Errorf("根据筛选条件获取题目失败：%v", err)
		}
	}

	return questions, nil
}
