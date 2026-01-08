package service

import (
	"errors"

	"github.com/vaynedu/exam_system/config"
	"github.com/vaynedu/exam_system/dao"
	"github.com/vaynedu/exam_system/model"
	"gorm.io/gorm"
)

// CreateCollectionService 创建收藏
func CreateCollectionService(questionID uint, tag, secondTag string) error {
	// 参数校验
	if questionID == 0 {
		return errors.New("题目ID不能为空")
	}

	// 检查题目是否存在
	questionDao := dao.NewQuestionDao(config.DB)
	questions, err := questionDao.GetQuestionsByIDList([]uint{questionID})
	if err != nil {
		return err
	}
	if len(questions) == 0 {
		return errors.New("题目不存在")
	}

	// 检查是否已收藏
	collectionDao := dao.NewCollectionDao(config.DB)
	_, err = collectionDao.GetCollectionByQuestionID(questionID)
	if err == nil {
		return errors.New("题目已收藏")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建收藏
	collection := &model.ExamQuestionCollection{
		QuestionID: questionID,
		Tag:        tag,
		SecondTag:  secondTag,
	}

	return collectionDao.CreateCollection(collection)
}

// DeleteCollectionService 删除收藏
func DeleteCollectionService(questionID uint) error {
	if questionID == 0 {
		return errors.New("题目ID不能为空")
	}

	collectionDao := dao.NewCollectionDao(config.DB)
	return collectionDao.DeleteCollection(questionID)
}

// GetCollectionStatusService 获取收藏状态
func GetCollectionStatusService(questionID uint) (bool, error) {
	if questionID == 0 {
		return false, errors.New("题目ID不能为空")
	}

	collectionDao := dao.NewCollectionDao(config.DB)
	_, err := collectionDao.GetCollectionByQuestionID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// BatchGetCollectionStatusService 批量获取收藏状态
func BatchGetCollectionStatusService(questionIDs []uint) (map[uint]bool, error) {
	if len(questionIDs) == 0 {
		return map[uint]bool{}, nil
	}

	collectionDao := dao.NewCollectionDao(config.DB)
	return collectionDao.BatchGetCollectionStatus(questionIDs)
}

// GetCollectionListService 获取收藏列表
func GetCollectionListService(tag, secondTag string, page, size int) ([]*model.ExamQuestionCollection, int64, error) {
	collectionDao := dao.NewCollectionDao(config.DB)
	return collectionDao.GetCollectionList(tag, secondTag, page, size)
}
