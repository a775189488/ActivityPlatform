package service

import (
	"entrytask/internal/common/logger"
	"entrytask/internal/model"
	"entrytask/internal/repository"
)

type IActService interface {
	CreateActivity(at *model.ActivityType) error
}

type ActService struct {
	Log  logger.ILogger           `inject:""`
	Repo repository.IActivityRepo `inject:""`
}

func (a *ActService) CreateActivity(at *model.ActivityType) error {
	return nil
}
