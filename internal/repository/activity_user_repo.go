package repository

import (
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
)

type IActivityUserRepo interface {
	InsertActivityUser(act *model.ActivityUser) error
	DeleteActivityUserByActId(id uint64) error
	DeleteActivityUserByUserId(id uint64) error
	DeleteActivityUserByActAndUser(userId uint64, actId uint64) error
	CheckActivityUserByActIdAndUserId(actId uint64, userId uint64) (bool, error)
	ListActivityUserByUserAndActs(userId uint64, actIds []uint64) ([]*model.ActivityUser, error)
}

type ActivityUserRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
}

func (a *ActivityUserRepo) ListActivityUserByUserAndActs(userId uint64, actIds []uint64) ([]*model.ActivityUser, error) {
	var result []*model.ActivityUser
	err := a.BaseRepo.Source.DB().Model(model.ActivityUser{}).Where("user_id = ? and act_id in (?)", userId, actIds).Find(&result).Error
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (a *ActivityUserRepo) CheckActivityUserByActIdAndUserId(actId uint64, userId uint64) (bool, error) {
	where := map[string]interface{}{
		"act_id":  actId,
		"user_id": userId,
	}
	var total int32
	if err := a.BaseRepo.GetCount(&model.ActivityUser{}, where, &total); err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (a *ActivityUserRepo) InsertActivityUser(actUser *model.ActivityUser) error {
	return a.BaseRepo.Create(actUser)
}

func (a *ActivityUserRepo) DeleteActivityUserByActId(id uint64) error {
	var actUser model.ActivityUser
	where := &model.ActivityUser{ActId: id}
	if _, err := a.BaseRepo.DeleteByWhere(&actUser, where); err != nil {
		return err
	} else {
		return nil
	}
}

func (a *ActivityUserRepo) DeleteActivityUserByUserId(id uint64) error {
	var actUser model.ActivityUser
	where := &model.ActivityUser{UserId: id}
	if _, err := a.BaseRepo.DeleteByWhere(&actUser, where); err != nil {
		return err
	} else {
		return nil
	}
}

func (a *ActivityUserRepo) DeleteActivityUserByActAndUser(userId uint64, actId uint64) error {
	var actUser model.ActivityUser
	where := &model.ActivityUser{ActId: actId, UserId: userId}
	if _, err := a.BaseRepo.DeleteByWhere(&actUser, where); err != nil {
		return err
	} else {
		return nil
	}
}
