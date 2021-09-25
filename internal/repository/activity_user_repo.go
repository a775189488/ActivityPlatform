package repository

import (
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
)

type IActivityUserRepo interface {
	InsertActivityUser(act *model.ActivityUser) error
	DeleteActivityUserByActId(id uint64) error
}

type ActivityUserRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
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
