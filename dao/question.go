package dao

import (
	"strconv"

	"github.com/vaynedu/exam_system/model"
	"gorm.io/gorm"
)

type QuestionDao struct {
	db *gorm.DB
}

func NewQuestionDao(db *gorm.DB) *QuestionDao {
	return &QuestionDao{
		db: db,
	}
}

// CreateQuestion 创建题目
func (q *QuestionDao) CreateQuestion(question *model.ExamQuestion) error {
	return q.db.Omit("CreatedAt").Create(&question).Error
}

// CreateQuestionsInBatches 批量创建题目
func (q *QuestionDao) CreateQuestionsInBatches(questions []model.ExamQuestion, batchSize int) error {
	return q.db.CreateInBatches(&questions, batchSize).Error
}

// GetRandomQuestions 随机获取指定数量的题目
func (q *QuestionDao) GetRandomQuestions(limit int) ([]model.ExamQuestion, error) {
	var questions []model.ExamQuestion
	err := q.db.Order("RAND()").Limit(limit).Find(&questions).Error
	return questions, err
}

// GetRandomQuestionsByTag 根据标签随机获取指定数量的题目
func (q *QuestionDao) GetRandomQuestionsByTag(tag, secondTag string, limit int) ([]model.ExamQuestion, error) {
	var questions []model.ExamQuestion
	query := q.db.Order("RAND()").Limit(limit)

	// 如果指定了标签，则添加标签过滤条件
	if tag != "" && secondTag != "" {
		query = query.Where("tag = ? AND second_tag = ?", tag, secondTag)
	} else if tag != "" {
		query = query.Where("tag = ?", tag)
	} else if secondTag != "" {
		query = query.Where("second_tag = ?", secondTag)
	}

	err := query.Find(&questions).Error
	return questions, err
}

// UpdateQuestion 更新题目
func (q *QuestionDao) UpdateQuestion(question *model.ExamQuestion) error {
	return q.db.Model(&model.ExamQuestion{}).Where("id = ?", question.ID).Updates(question).Error
}

// DeleteQuestion 删除题目
func (q *QuestionDao) DeleteQuestion(id uint) error {
	return q.db.Delete(&model.ExamQuestion{}, id).Error
}

// GetAllQuestions 获取全部题目
func (q *QuestionDao) GetAllQuestions() ([]model.ExamQuestion, error) {
	var questions []model.ExamQuestion
	err := q.db.Find(&questions).Error
	return questions, err
}

// GetQuestionsByIDList 根据ID列表获取题目
func (q *QuestionDao) GetQuestionsByIDList(ids []uint) ([]model.ExamQuestion, error) {
	var questions []model.ExamQuestion
	if len(ids) == 0 {
		return questions, nil
	}
	err := q.db.Where("id IN ?", ids).Find(&questions).Error
	return questions, err
}

// GetQuestionsByFilter 根据筛选条件获取题目
func (q *QuestionDao) GetQuestionsByFilter(tag, secondTag, questionType, keyword string) ([]model.ExamQuestion, error) {
	var questions []model.ExamQuestion
	query := q.db.Model(&model.ExamQuestion{})

	// 添加筛选条件
	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	if secondTag != "" {
		query = query.Where("second_tag = ?", secondTag)
	}
	if questionType != "" {
		// 转换题型为数字
		typeInt, err := strconv.Atoi(questionType)
		if err == nil {
			query = query.Where("question_type = ?", int8(typeInt))
		}
	}
	if keyword != "" {
		// 在题干中搜索关键词
		query = query.Where("question_title LIKE ?", "%"+keyword+"%")
	}

	err := query.Find(&questions).Error
	return questions, err
}
