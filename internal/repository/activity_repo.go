package repository

import (
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
)

type IActivityRepo interface {
	GetActivityById(id uint64) (*model.Activity, error)
	InsertActivity(act *model.Activity) error
	UpdateActivity(act *model.Activity) error
	DeleteActivity(id uint64) error
	ListActivity(page, size int32, where interface{}) ([]*model.Activity, error)
	ListActivityUserById(page, size int32, actId uint64) ([]*model.User, error)
	GetActivityDetail(id uint64) (*model.ActivityDetail, error)
}

type ActivityRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
}

func (r *ActivityRepo) GetActivityById(id uint64) (*model.Activity, error) {
	var act model.Activity
	if err := r.BaseRepo.FirstByID(&act, id); err != nil {
		return nil, err
	}

	return &act, nil
}

func (r *ActivityRepo) InsertActivity(act *model.Activity) error {
	if err := r.BaseRepo.Create(act); err != nil {
		r.Log.Errorf("[ActRepo]insert act(%v) fail, err: %v", *act, err)
		return err
	}
	return nil
}

func (r *ActivityRepo) UpdateActivity(act *model.Activity) error {
	if err := r.BaseRepo.Source.DB().Model(&act).Update(act).Error; err != nil {
		r.Log.Errorf("[ActRepo]update act(%v) fail, err: %v", *act, err)
		return err
	}
	return nil
}

func (r *ActivityRepo) DeleteActivity(id uint64) error {
	act := model.Activity{}
	where := &model.Activity{Id: id}
	if _, err := r.BaseRepo.DeleteByWhere(&act, where); err != nil {
		r.Log.Errorf("[ActRepo]delete act(%d) fail, err: %v", id, err)
		return err
	} else {
		return nil
	}
}

func (r *ActivityRepo) ListActivity(page, size int32, where interface{}) ([]*model.Activity, error) {
	var acts []*model.Activity
	if err := r.BaseRepo.GetPagesNotCount(&model.Activity{}, &acts, page, size, where); err != nil {
		r.Log.Errorf("[ActRepo]list act fail, condition(%v), err: %v", where, err)
		return nil, err
	}
	return acts, nil
}

func (r *ActivityRepo) GetActivityDetail(id uint64) (*model.ActivityDetail, error) {
	var act model.ActivityDetail
	err := r.BaseRepo.Source.DB().Table(model.Activity{}.TableName()).Select("act_tab.*, act_type_tab.name as act_type_name").
		Joins("left join act_type_tab on act_type_tab.id=act_tab.act_type").First(&act, id).Error
	if err != nil {
		return nil, err
	}
	return &act, err
}

func (r *ActivityRepo) ListActivityUserById(page, size int32, actId uint64) ([]*model.User, error) {
	var users []*model.User
	where := map[string]interface{}{"act_id": actId}
	db := r.BaseRepo.Source.DB().Table(model.ActivityUser{}.TableName()).Select("user_tab.*").Joins(
		"left join user_tab on act_user_tab.user_id=user_tab.id").Where(where)
	err := db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
