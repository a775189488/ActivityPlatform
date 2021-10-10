package repository

import (
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
	"github.com/jinzhu/gorm"
)

type IActivityCommentRepo interface {
	GetCommentById(id uint64) (*model.ActivityComment, error)
	ListComment(page, size int32, total *int32, where interface{}) ([]*model.ActivityComment, error)
	InsertComment(comment *model.ActivityComment) error
	UpdateComment(comment *model.ActivityComment) error
	DeleteComment(id uint64) error
}
type ActivityCommentRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
}

func (a *ActivityCommentRepo) GetCommentById(id uint64) (*model.ActivityComment, error) {
	var actComment model.ActivityComment
	if err := a.BaseRepo.FirstByID(&actComment, id); err != nil {
		if err != gorm.ErrRecordNotFound {
			a.Log.Errorf("[ActCommentRepo]get actComment(%d) fail, err: %v", id, err)
		} else {
			a.Log.Infof("[ActCommentRepo]get actComment(%d) not found", id)
		}
		return nil, err
	}

	return &actComment, nil
}

func (a *ActivityCommentRepo) ListComment(page, size int32, total *int32, where interface{}) ([]*model.ActivityComment, error) {
	var actComments []*model.ActivityComment
	if err := a.BaseRepo.GetPages(&model.ActivityComment{}, &actComments, page, size, total, where); err != nil {
		a.Log.Errorf("[ActCommentRepo]list actComment fail, condition(%v), err: %v", where, err)
		return nil, err
	}
	return actComments, nil
}

func (a *ActivityCommentRepo) InsertComment(comment *model.ActivityComment) error {
	if err := a.BaseRepo.Create(comment); err != nil {
		a.Log.Errorf("[ActCommentRepo]insert actComment(%v) fail, err: %v", *comment, err)
		return err
	}
	return nil
}

func (a *ActivityCommentRepo) UpdateComment(comment *model.ActivityComment) error {
	if err := a.BaseRepo.Source.DB().Model(&comment).Update(comment).Error; err != nil {
		a.Log.Errorf("[ActCommentRepo]update actComment(%v) fail, err: %v", *comment, err)
		return err
	}
	return nil
}

func (a *ActivityCommentRepo) DeleteComment(id uint64) error {
	actComment := model.ActivityComment{}
	where := &model.ActivityComment{Id: id}
	if _, err := a.BaseRepo.DeleteByWhere(&actComment, where); err != nil {
		a.Log.Errorf("[ActCommentRepo]delete actComment(%d) fail, err: %v", id, err)
		return err
	}
	return nil
}
