package dao

import (
	"github.com/vaynedu/exam_system/model"
	"gorm.io/gorm"
)

// CollectionDao 收藏题目DAO
type CollectionDao struct {
	db *gorm.DB
}

// NewCollectionDao 创建收藏题目DAO实例
func NewCollectionDao(db *gorm.DB) *CollectionDao {
	return &CollectionDao{
		db: db,
	}
}

// CreateCollection 创建收藏
func (d *CollectionDao) CreateCollection(collection *model.ExamQuestionCollection) error {
	return d.db.Create(collection).Error
}

// DeleteCollection 删除收藏
func (d *CollectionDao) DeleteCollection(questionID uint) error {
	return d.db.Where("question_id = ?", questionID).Delete(&model.ExamQuestionCollection{}).Error
}

// GetCollectionByQuestionID 根据题目ID获取收藏
func (d *CollectionDao) GetCollectionByQuestionID(questionID uint) (*model.ExamQuestionCollection, error) {
	var collection model.ExamQuestionCollection
	err := d.db.Where("question_id = ?", questionID).First(&collection).Error
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

// GetCollectionList 获取收藏列表
func (d *CollectionDao) GetCollectionList(tag, secondTag string, page, size int) ([]*model.ExamQuestionCollection, int64, error) {
	var collections []*model.ExamQuestionCollection
	var total int64
	
	query := d.db.Model(&model.ExamQuestionCollection{})
	
	// 筛选条件
	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	if secondTag != "" {
		query = query.Where("second_tag = ?", secondTag)
	}
	
	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 分页查询
	offset := (page - 1) * size
	err := query.Preload("Question").Offset(offset).Limit(size).Order("created_at DESC").Find(&collections).Error
	if err != nil {
		return nil, 0, err
	}
	
	return collections, total, nil
}

// BatchGetCollectionStatus 批量获取题目收藏状态
func (d *CollectionDao) BatchGetCollectionStatus(questionIDs []uint) (map[uint]bool, error) {
	var collections []*model.ExamQuestionCollection
	err := d.db.Where("question_id IN ?", questionIDs).Find(&collections).Error
	if err != nil {
		return nil, err
	}
	
	// 构建结果映射
	result := make(map[uint]bool)
	for _, collection := range collections {
		result[collection.QuestionID] = true
	}
	
	return result, nil
}
