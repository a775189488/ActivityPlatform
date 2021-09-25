package repository

import (
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
	"github.com/jinzhu/gorm"
)

type IActivityRepo interface {
	GetActivityById(id uint64) *model.Activity
	InsertActivity(act *model.Activity) bool
	UpdateActivity(act *model.Activity) bool
	DeleteActivity(id uint64) bool
	ListActivity(page, size int32, total *int32, where interface{}) []*model.Activity
	//ListActivityUserById(page, size int32, total *int32, actId uint64) []*model.User
	GetActivityDetail(id uint64) (*model.ActivityDetail, error)
}

type ActivityRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
}

func (r *ActivityRepo) GetActivityById(id uint64) *model.Activity {
	var act model.Activity
	if err := r.BaseRepo.FirstByID(&act, id); err != nil {
		if err != gorm.ErrRecordNotFound {
			r.Log.Errorf("[ActRepo]get act(%d) fail, err: %v", id, err)
		} else {
			r.Log.Infof("[ActRepo]get act(%d) not found", id)
		}
	}
	return &act
}

func (r *ActivityRepo) InsertActivity(act *model.Activity) bool {
	if err := r.BaseRepo.Create(act); err != nil {
		r.Log.Errorf("[ActRepo]insert act(%v) fail, err: %v", *act, err)
		return false
	}
	return true
}

func (r *ActivityRepo) UpdateActivity(act *model.Activity) bool {
	if err := r.BaseRepo.Source.DB().Model(&act).Update(act).Error; err != nil {
		r.Log.Errorf("[ActRepo]update act(%v) fail, err: %v", *act, err)
		return false
	}
	return true
}

func (r *ActivityRepo) DeleteActivity(id uint64) bool {
	act := model.Activity{}
	where := &model.Activity{Id: id}
	if count, err := r.BaseRepo.DeleteByWhere(&act, where); err != nil {
		r.Log.Errorf("[ActRepo]delete act(%d) fail, err: %v", id, err)
		return false
	} else {
		return count > 0
	}
}

func (r *ActivityRepo) ListActivity(page, size int32, total *int32, where interface{}) []*model.Activity {
	var acts []*model.Activity
	if err := r.BaseRepo.GetPages(&model.Activity{}, &acts, page, size, total, where); err != nil {
		r.Log.Errorf("[ActRepo]list act fail, condition(%v), err: %v", where, err)
	}
	return acts
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

func (r *ActivityRepo) GetUserByActivityId(actId uint64, page, size int32, total *int32) ([]*model.User, error) {
	var users []*model.User
	where := map[string]interface{}{"act_id": actId}
	db := r.BaseRepo.Source.DB().Table(model.ActivityUser{}.TableName()).Select("user_tab.*").Joins(
		"left join user_tab on act_user_tab.user_id=user_tab.id").Where(where)
	err := db.Count(total).Error
	if err != nil {
		return nil, err
	}
	if *total == 0 {
		return users, nil
	}
	err = db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

//func (r *ActivityRepo) ListActivityUserById(page, size int32, total *int32, actId uint64) []*model.User {
//	var acts []*model.User
//	var actUsers []*model.ActivityUser
//	cond := map[string]interface{} {
//		"act_id": actId,
//	}
//	if err := r.BaseRepo.GetPages(&model.ActivityUser{}, &actUsers, page, size, total, cond); err != nil {
//		r.Log.Errorf("[ActRepo]list act fail, condition(%v), err: %v", cond, err)
//	}
//	if len(actUsers) == 0 {
//		return acts
//	}
//	var userIds []uint64
//	userIdMap := map[uint64]struct{} {}
//	for _, au := range actUsers {
//		_, ok := userIdMap[au.UserId]
//		if !ok {
//			userIds = append(userIds, au.UserId)
//			userIdMap[au.UserId] = struct{}{}
//		}
//	}
//}
