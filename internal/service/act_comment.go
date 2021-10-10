package service

import (
	error2 "entrytask/internal/common/error"
	"entrytask/internal/common/logger"
	"entrytask/internal/conf"
	"entrytask/internal/model"
	"entrytask/internal/repository"
)

type IActCommentService interface {
	CreateActComment(ac *model.ActivityComment) error
	ListActCommentByActId(page, size int32, id uint64) ([]*model.ActivityComment, int, error)
	UpdateActComment(id uint64, ac *model.ActivityComment) (*model.ActivityComment, error)
	DeleteActComment(id uint64) error
}

type ActCommentService struct {
	Repo     repository.IActivityCommentRepo `inject:""`
	UserRepo repository.IUserRepo            `inject:""`
	ActRepo  repository.IActivityRepo        `inject:""`
	Log      logger.ILogger                  `inject:""`
}

func (a *ActCommentService) ListActCommentByActId(page, size int32, id uint64) ([]*model.ActivityComment, int, error) {
	if size == 0 {
		size = int32(conf.Config.App.PageSize)
	}
	var total int32
	where := map[string]interface{}{"act_id": id}
	users, err := a.Repo.ListComment(page, size, &total, where)
	if err != nil {
		return nil, 0, err
	}
	return users, int(total), nil
}

func (a *ActCommentService) CreateActComment(ac *model.ActivityComment) error {
	a.Log.Infof("[ActCommentService]create activity comment(%v)", *ac)
	_, err := a.UserRepo.GetUserById(ac.UserId)
	if err != nil {
		a.Log.Errorf("[ActCommentService]create activity comment(%v) fail, user not exist", *ac)
		return error2.ActCommentCreateUserNotExistError
	}
	_, err = a.ActRepo.GetActivityById(ac.ActId)
	if err != nil {
		a.Log.Errorf("[ActCommentService]create activity comment(%v) fail, act not exist", *ac)
		return error2.ActCommentCreateActNotExistError
	}
	if err := a.Repo.InsertComment(ac); err != nil {
		a.Log.Errorf("")
		return err
	}
	return nil
}

func (a *ActCommentService) UpdateActComment(id uint64, ac *model.ActivityComment) (*model.ActivityComment, error) {
	a.Log.Infof("[ActCommentService]update activity comment(%v)", *ac)
	oldAc, err := a.Repo.GetCommentById(id)
	if err != nil {
		a.Log.Errorf("[ActCommentService]get activity comment(%d) fail, err: %v", id, err)
		return nil, err
	}
	if !oldAc.CompareAndSwap(ac) {
		return oldAc, nil
	}
	if err := a.Repo.UpdateComment(ac); err != nil {
		a.Log.Errorf("[ActCommentService]update activity comment(%v) fail, err: %v", *ac, err)
		return nil, err
	}
	return oldAc, err
}

func (a *ActCommentService) DeleteActComment(id uint64) error {
	a.Log.Infof("[ActCommentService]delete activity comment(%d)", id)
	err := a.Repo.DeleteComment(id)
	if err != nil {
		a.Log.Errorf("[ActCommentService] delete activity comment fail, err: %v", err)
	}
	return err
}
