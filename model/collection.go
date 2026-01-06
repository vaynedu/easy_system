package model

import (
	"time"
)

// ExamQuestionCollection 收藏题目模型
type ExamQuestionCollection struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionID uint      `json:"question_id" gorm:"not null;uniqueIndex:uk_question_id"`
	Tag        string    `json:"tag" gorm:"not null;index:idx_tag"`
	SecondTag  string    `json:"second_tag" gorm:"not null;index:idx_second_tag"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	// 关联关系
	Question *ExamQuestion `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
}

// TableName 指定表名
func (ExamQuestionCollection) TableName() string {
	return "exam_question_collection"
}
