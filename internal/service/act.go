package service

import (
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
	"entrytask/internal/repository"
)

type IActService interface {
	CreateActivity(act *model.Activity) error
	UpdateActivity(id uint64, act *model.Activity) (*model.Activity, error)
	GetActivityDetail(id uint64) (*model.ActivityDetail, error)
	DeleteActivity(id uint64) error
	ListActivity(page, size int32, where map[string]interface{}) ([]*model.Activity, error)
	GetActivityUser(page, size int32, actId uint64) ([]*model.User, error)
	CheckUserJoinActivity(actId uint64, userId uint64) (bool, error)
	CheckUserJoinActivities(actId []uint64, userId uint64) ([]uint64, error)
}

type ActService struct {
	Log    logger.ILogger               `inject:""`
	Repo   repository.IActivityRepo     `inject:""`
	AtRepo repository.IActivityTypeRepo `inject:""`
	AuRepo repository.IActivityUserRepo `inject:""`
}

func (a *ActService) CreateActivity(act *model.Activity) error {
	a.Log.Infof("[ActivityService]create activity(%v)", act)
	// check type
	if _, err := a.AtRepo.GetActivityTypeById(act.ActType); err != nil {
		a.Log.Errorf("[ActivityService]create activity(%v) fail, get act type fail, err: %v", *act, err)
		return error2.ActCreateTypeNotFoundError
	}
	err := a.Repo.InsertActivity(act)
	if err != nil {
		a.Log.Errorf("[ActivityService]create activity(%v) fail, err: %v", *act, err)
		return nil
	}
	return nil
}

func (a *ActService) GetActivityDetail(actId uint64) (*model.ActivityDetail, error) {
	result, err := a.Repo.GetActivityDetail(actId)
	if err != nil {
		a.Log.Errorf("[ActivityService]get activity(%d) detail fail, err: %v", actId, err)
		return nil, err
	} else {
		return result, nil
	}
}

func (a *ActService) CheckUserJoinActivity(actId uint64, userId uint64) (bool, error) {
	isIn, err := a.AuRepo.CheckActivityUserByActIdAndUserId(actId, userId)
	if err != nil {
		a.Log.Errorf("[ActivityService] check user(%d) and activity(%d) fail, err: %v", userId, actId, err)
		return false, err
	}
	return isIn, nil
}

func (a *ActService) UpdateActivity(id uint64, act *model.Activity) (*model.Activity, error) {
	a.Log.Infof("[ActivityService]update activity(%d) to %v", id, *act)
	oldAct, err := a.Repo.GetActivityById(id)
	if err != nil {
		a.Log.Errorf("[ActivityService]update activity(%d) fail, get object fail, err: %v", id, err)
		return nil, err
	}
	if oldAct.ActType != act.ActType {
		if _, err := a.AtRepo.GetActivityTypeById(act.ActType); err != nil {
			a.Log.Errorf("[ActivityService]update activity(%v) fail, get act type fail, err: %v", *act, err)
			return nil, error2.ActCreateTypeNotFoundError
		}
	}
	if oldAct.CompareAndSwap(act) {
		err = a.Repo.UpdateActivity(oldAct)
		if err != nil {
			a.Log.Errorf("[ActivityService]update activity(%d) to (%v) fail, err: %v", id, *act, err)
			return nil, err

		}
	}
	return oldAct, nil
}

func (a *ActService) DeleteActivity(id uint64) error {
	a.Log.Infof("[ActivityService]delete activity(%d)", id)
	// todo 删除相关评论？
	err := a.Repo.DeleteActivity(id)
	if err != nil {
		a.Log.Errorf("[ActivityService]delete activity(%d) fail, err: %v", id, err)
	}
	return err
}

func (a *ActService) GetActivityUser(page, size int32, actId uint64) ([]*model.User, error) {
	act, err := a.Repo.GetActivityById(actId)
	if err != nil {
		a.Log.Errorf("[ActivityService]get activity(%d) user fail, get object fail, err: %v", act, err)
		return nil, err
	}
	users, err := a.Repo.ListActivityUserById(page, size, actId)
	if err != nil {
		a.Log.Errorf("[ActivityService]get activity(%d) user fail, err: %v", act, err)
		return nil, err
	}
	return users, nil
}

func (a *ActService) ListActivity(page, size int32, where map[string]interface{}) ([]*model.Activity, error) {
	result, err := a.Repo.ListActivity(page, size, where)
	if err != nil {
		a.Log.Errorf("[ActivityService]list activity fail, where(%v), err: %v", where, err)
		return nil, err
	}
	return result, nil
}

func (a *ActService) CheckUserJoinActivities(actId []uint64, userId uint64) ([]uint64, error) {
	aus, err := a.AuRepo.ListActivityUserByUserAndActs(userId, actId)
	if err != nil {
		a.Log.Errorf("[ActivityService]check user(%d) join activities(%v) fail,err: %v", userId, actId, err)
		return nil, err
	}
	result := make([]uint64, 0)
	for _, t := range aus {
		result = append(result, t.ActId)
	}
	return result, nil
}
