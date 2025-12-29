package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/consts"
	"github.com/vaynedu/exam_system/dao"
	"github.com/vaynedu/exam_system/model"
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

	// 4. 调用DAO层插入数据
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
