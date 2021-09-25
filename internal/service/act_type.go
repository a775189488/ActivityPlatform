package service

import (
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
	"entrytask/internal/repository"
)

type IActTypeService interface {
	CreateActType(at *model.ActivityType) error
}

type ActTypeService struct {
	Log  logger.ILogger               `inject:""`
	Repo repository.IActivityTypeRepo `inject:""`
}

func (a *ActTypeService) CreateActType(at *model.ActivityType) error {
	a.Log.Infof("[ActTypeService] create activity type (%v)", *at)
	// check parent type exist
	if at.Parent != 0 {
		_, err := a.Repo.GetActivityTypeById(at.Parent)
		if err != nil {
			if err == error2.ActTypeNotFoundError {
				a.Log.Errorf("[ActTypeService] create activity type(%v) fail, can not find parent(%d)", *at, at.Parent)
				return error2.ActTypeParentNotFoundError
			}
			return err
		}
	}
	err := a.Repo.InsertActivityType(at)
	if err != nil {
		a.Log.Errorf("[ActTypeService] create activity type(%v) fail, err: %v", *at, err)
		return err
	}
	return nil
}
