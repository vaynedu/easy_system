package model

import "time"

const ExamQuestionsTableName = "exam_questions"

// ExamQuestion 题目模型（适配GORM）
type ExamQuestion struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionType   int8      `gorm:"column:question_type;not null" json:"question_type"` // tinyint对应int8
	QuestionTitle  string    `gorm:"column:question_title;type:varchar(500);not null" json:"question_title"`
	OptionA        string    `gorm:"column:option_a;type:varchar(200);default:''" json:"option_a"`
	OptionB        string    `gorm:"column:option_b;type:varchar(200);default:''" json:"option_b"`
	OptionC        string    `gorm:"column:option_c;type:varchar(200);default:''" json:"option_c"`
	OptionD        string    `gorm:"column:option_d;type:varchar(200);default:''" json:"option_d"`
	CorrectAnswer  string    `gorm:"column:correct_answer;type:varchar(1000);not null" json:"correct_answer"`
	AnswerAnalysis string    `gorm:"column:answer_analysis;type:varchar(2000);default:''" json:"answer_analysis"`
	QuestionRemark string    `gorm:"column:question_remark;type:varchar(500);default:''" json:"question_remark"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`               // 自动生成创建时间
	Tag            string    `gorm:"column:tag;type:varchar(50);default:''" json:"tag"`                // 对应一级分类（KnowledgeTree.Name）
	SecondTag      string    `gorm:"column:second_tag;type:varchar(100);default:''" json:"second_tag"` // 对应二级分类（KnowledgeTree.SecondTag）
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	UploadType     int8      `gorm:"column:upload_type;not null" json:"upload_type"` // 题目录入方式，默认0=手动 1=excel表格 2=豆包AI 3=阿里AI 4=云雾AI
}

// TableName 指定表名（GORM默认复数，需显式指定）
func (ExamQuestion) TableName() string {
	return "exam_questions"
}
