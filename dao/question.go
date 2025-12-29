package dao

import (
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
