package service

import (
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/conf"
	"entrytask/internal/model"
	"entrytask/internal/repository"
)

type IActTypeService interface {
	CreateActType(at *model.ActivityType) error
	ListActType(page, size int32) ([]*model.ActivityType, int, error)
	DeleteActType(id uint64) error
	UpdateActType(id uint64, at *model.ActivityType) (*model.ActivityType, error)
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

func (a *ActTypeService) ListActType(page, size int32) ([]*model.ActivityType, int, error) {
	if size == 0 {
		size = int32(conf.Config.App.PageSize)
	}
	var total int32
	users, err := a.Repo.ListActivityType(page, size, &total, nil)
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

func (a *ActTypeService) UpdateActType(id uint64, at *model.ActivityType) (*model.ActivityType, error) {
	a.Log.Infof("[ActTypeService]update activity type(%d) to (%v)", id, *at)
	oldAt, err := a.Repo.GetActivityTypeById(id)
	if err != nil {
		a.Log.Errorf("[ActTypeService]update activity type by id(%d)fail, err: %v", id, err)
		return nil, err
	}
	if !oldAt.CompareAndSwap(at) {
		return oldAt, nil
	}
	err = a.Repo.UpdateActivityType(oldAt)
	if err != nil {
		a.Log.Errorf("[ActTypeService]update activity type (%v) fail, err: %v", at, err)
		return nil, err
	}
	return oldAt, nil
}

func (a *ActTypeService) DeleteActType(id uint64) error {
	// todo 子类型检查
	a.Log.Infof("[ActTypeService]delete activity type by id(%d)", id)

	actCount, err := a.Repo.GetActCountByActType(id)
	if err != nil {
		a.Log.Errorf("[ActTypeService]list activity by typeid(%d)fail, err: %v", id, err)
		return err
	}
	if actCount > 0 {
		a.Log.Errorf("[ActTypeService] it has %d activities in this type(%d), forbidden to delete!", actCount, id)
		return error2.ActTypeDeleteError
	}

	err = a.Repo.DeleteActivityType(id)
	if err != nil {
		a.Log.Errorf("[ActTypeService]delete activity type by id(%d)fail, err: %v", id, err)
	}
	return err
}
