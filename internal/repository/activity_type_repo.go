package repository

import (
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
	"github.com/jinzhu/gorm"
)

type IActivityTypeRepo interface {
	GetActivityTypeById(id uint64) (*model.ActivityType, error)
	ListActivityType(page, size int32, total *int32, where interface{}) ([]*model.ActivityType, error)
	InsertActivityType(actType *model.ActivityType) error
	UpdateActivityType(actType *model.ActivityType) error
	DeleteActivityType(id uint64) error
	GetActCountByActType(id uint64) (int, error)
}

type ActivityTypeRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
}

func (a *ActivityTypeRepo) GetActivityTypeById(id uint64) (*model.ActivityType, error) {
	var actType model.ActivityType
	if err := a.BaseRepo.FirstByID(&actType, id); err != nil {
		if err != gorm.ErrRecordNotFound {
			a.Log.Errorf("[ActTypeRepo]get actType(%d) fail, err: %v", id, err)
			return nil, error2.ActTypeNotFoundError
		} else {
			a.Log.Infof("[ActTypeRepo]get actType(%d) not found", id)
			return nil, err
		}
	}
	return &actType, nil
}

func (a *ActivityTypeRepo) ListActivityType(page, size int32, total *int32, where interface{}) ([]*model.ActivityType, error) {
	var users []*model.ActivityType
	if err := a.BaseRepo.GetPages(&model.ActivityType{}, &users, page, size, total, where); err != nil {
		a.Log.Errorf("[ActTypeRepo]list actType fail, condition(%v), err: %v", where, err)
		return nil, err
	}
	return users, nil
}

func (a *ActivityTypeRepo) InsertActivityType(actType *model.ActivityType) error {
	if err := a.BaseRepo.Create(actType); err != nil {
		a.Log.Errorf("[ActTypeRepo]insert actType(%v) fail, err: %v", *actType, err)
		return err
	}
	return nil
}

func (a *ActivityTypeRepo) UpdateActivityType(actType *model.ActivityType) error {
	if err := a.BaseRepo.Source.DB().Model(&actType).Update(actType).Error; err != nil {
		a.Log.Errorf("[ActTypeRepo]update actType(%v) fail, err: %v", *actType, err)
		return err
	}
	return nil
}

func (a *ActivityTypeRepo) DeleteActivityType(id uint64) error {
	user := model.ActivityType{}
	where := &model.ActivityType{Id: id}
	if _, err := a.BaseRepo.DeleteByWhere(&user, where); err != nil {
		a.Log.Errorf("[ActTypeRepo]delete actType(%d) fail, err: %v", id, err)
		return err
	}
	return nil
}

func (a *ActivityTypeRepo) GetActCountByActType(id uint64) (int, error) {
	var total int32
	where := &model.Activity{ActType: id}
	if err := a.BaseRepo.GetCount(model.Activity{}, where, &total); err != nil {
		return 0, err
	}
	return int(total), nil
}
