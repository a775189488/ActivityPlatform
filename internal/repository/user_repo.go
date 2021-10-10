package repository

import (
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
)

type IUserRepo interface {
	GetUserById(int uint64) (*model.User, error)
	ListUser(page, size int32, total *int32, where interface{}) ([]*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	CheckUser(loginName, password string) error
	InsertUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id uint64) error
	GetUserCountByUsername(username string) (int, error)
	ListUserActivity(page, size int32, userId uint64) ([]*model.Activity, error)
}

type UserRepo struct {
	Log      logger.ILogger `inject:""`
	BaseRepo BaseRepo       `inject:"inline"`
}

func (r *UserRepo) GetUserById(id uint64) (*model.User, error) {
	var user model.User
	if err := r.BaseRepo.FirstByID(&user, id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) InsertUser(user *model.User) error {
	return r.BaseRepo.Create(user)
}

func (r *UserRepo) UpdateUser(user *model.User) error {
	return r.BaseRepo.Source.DB().Model(&user).Update(user).Error
}

func (r *UserRepo) DeleteUser(id uint64) error {
	user := model.User{}
	where := &model.User{Id: id}
	if _, err := r.BaseRepo.DeleteByWhere(&user, where); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *UserRepo) ListUser(page, size int32, total *int32, where interface{}) ([]*model.User, error) {
	var users []*model.User
	if err := r.BaseRepo.GetPages(&model.User{}, &users, page, size, total, where); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepo) CheckUser(loginName, password string) error {
	var count int
	r.BaseRepo.Source.DB().Table(model.User{}.TableName()).Where("username = ? and password = ?", loginName, password).Count(&count)
	if count > 0 {
		return nil
	}
	r.BaseRepo.Source.DB().Table(model.User{}.TableName()).Where("username = ?", loginName).Count(&count)
	if count > 0 {
		return error2.UserPasswordError
	} else {
		return error2.UserNotFoundError
	}
}

func (r *UserRepo) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	where := model.User{Username: username}
	if err := r.BaseRepo.First(&where, &user, "*"); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserCountByUsername(username string) (int, error) {
	var count int
	r.BaseRepo.Source.DB().Table(model.User{}.TableName()).Where("username = ?", username).Count(&count)
	return count, nil
}

func (r *UserRepo) ListUserActivity(page, size int32, userId uint64) ([]*model.Activity, error) {
	var acts []*model.Activity
	where := map[string]interface{}{"user_id": userId}
	db := r.BaseRepo.Source.DB().Table(model.ActivityUser{}.TableName()).Select("act_tab.*").Joins(
		"left join act_tab on act_user_tab.act_id=act_tab.id").Where(where)
	err := db.Offset((page - 1) * size).Limit(size).Find(&acts).Error
	if err != nil {
		return nil, err
	}
	return acts, nil
}
