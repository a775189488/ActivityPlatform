package service

import (
	"encoding/base64"
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/utils"
	"entrytask/internal/conf"
	"entrytask/internal/model"
	"entrytask/internal/repository"
	"github.com/jinzhu/gorm"
)

type IUserSerivce interface {
	CheckUser(username string, password string) error
	GetUserInfoByLoginName(username string) (*model.User, error)
	CreateUser(user *model.User) error
	ListUser(page, size int32) ([]*model.User, int, error)
	GetUserById(id uint64) (*model.User, error)
	UpdateUser(id uint64, user *model.User) (*model.User, error)
	DeleteUser(id uint64) error
	JoinActivity(userId uint64, actId uint64) error
	LeaveActivity(userId uint64, actId uint64) error
}

type UserService struct {
	Log     logger.ILogger               `inject:""`
	Repo    repository.IUserRepo         `inject:""`
	AuRepo  repository.IActivityUserRepo `inject:""`
	ActRepo repository.ActivityRepo      `inject:""`
}

func (u *UserService) DeleteUser(id uint64) error {
	u.Log.Infof("[UserService]delete user(%d)", id)
	// todo 考虑是否需要顺便删除所属活动、记录、comment等
	err := u.Repo.DeleteUser(id)
	if err != nil {
		u.Log.Errorf("[UserService]delete user by id(%d)fail, err: %v", id, err)
	}
	return err
}

func (u *UserService) UpdateUser(id uint64, user *model.User) (*model.User, error) {
	u.Log.Infof("[UserService]update user(%d) to (%v)", id, *user)
	oldUser, err := u.GetUserById(id)
	if err != nil {
		u.Log.Errorf("[UserService]Get user by id(%d)fail, err: %v", id, err)
		return nil, err
	}
	if !oldUser.CompareAndSwap(user) {
		return oldUser, nil
	}
	err = u.Repo.UpdateUser(oldUser)
	if err != nil {
		u.Log.Errorf("[UserService]update user(%v) fail, err: %v", user, err)
		return nil, err
	}
	return oldUser, nil
}

func (u *UserService) GetUserById(id uint64) (*model.User, error) {
	user, err := u.Repo.GetUserById(id)
	if err != nil {
		u.Log.Errorf("[UserService]Get user by id(%d)fail, err: %v", id, err)
		return nil, error2.UserNotFoundError
	}
	return user, nil
}

func (u *UserService) ListUser(page, size int32) ([]*model.User, int, error) {
	if size == 0 {
		size = int32(conf.Config.App.PageSize)
	}
	var total int32
	users, err := u.Repo.ListUser(page, size, &total, nil)
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

func (u *UserService) CheckUser(username string, password string) error {
	return u.Repo.CheckUser(username, password)
}

func (u *UserService) GetUserInfoByLoginName(username string) (*model.User, error) {
	user, err := u.Repo.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, error2.UserNotFoundError
		}
		return nil, err
	}
	return user, nil
}

func (u *UserService) CreateUser(user *model.User) error {
	u.Log.Infof("[UserService]Create user %v", *user)
	count, err := u.Repo.GetUserCountByUsername(user.Username)
	if err != nil {
		u.Log.Errorf("[UserService]Get user by username(%s)fail, err: %v", user.Username, err)
		return err
	}
	if count > 0 {
		u.Log.Errorf("[UserService]create user(%v) fail, username conflict!", user)
		return error2.UserNameConflictError
	}
	// handle password
	pwd, _ := base64.StdEncoding.DecodeString(user.Password)
	realPwd := utils.AesCbcDecrypt(pwd, []byte(conf.Config.App.AesKey))
	user.Password = utils.Md5Encrypt(realPwd)

	// todo check head picture useful ?

	if err := u.Repo.InsertUser(user); err != nil {
		u.Log.Errorf("[UserService]create user(%v) fail, err: %v", err)
		return err
	}
	return nil
}

func (u *UserService) JoinActivity(userId uint64, actId uint64) error {
	// todo 检查活动有效
	au := &model.ActivityUser{UserId: userId, ActId: actId}
	if err := u.AuRepo.InsertActivityUser(au); err != nil {
		return err
	}
	return nil
}

func (u *UserService) LeaveActivity(userId uint64, actId uint64) error {
	if err := u.AuRepo.DeleteActivityUserByActAndUser(userId, actId); err != nil {
		return err
	}
	return nil
}
