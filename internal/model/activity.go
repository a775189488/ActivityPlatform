package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Activity struct {
	Id          uint64 `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	BeginAt     uint64 `json:"begin_at"`
	EndAt       uint64 `json:"end_at"`
	Description string `json:"description"`
	Creator     uint64 `json:"creator"`
	ActType     uint64 `json:"act_type"`
	Address     string `json:"address"`
	Statue      int    `json:"status"`
	CreateTime  uint64 `json:"create_time"`
	UpdateTime  uint64 `json:"update_time"`
}

type ActivityDetail struct {
	Activity
	ActTypeName string `json:"act_type_name"`
}

func (Activity) TableName() string {
	return "act_tab"
}

func (a *Activity) BeforeCreate(scope *gorm.Scope) error {
	err := scope.SetColumn("CreateTime", time.Now().Unix())
	err = scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (a *Activity) BeforeUpdate(scope *gorm.Scope) error {
	err := scope.SetColumn("UpdateTime", time.Now().Unix())
	return err
}

func (a *Activity) CompareAndSwap(act *Activity) bool {
	isSwap := false
	if a.Title != act.Title {
		a.Title = act.Title
		isSwap = true
	}
	if a.Description != act.Description {
		a.Description = act.Description
		isSwap = true
	}
	if a.ActType != act.ActType {
		a.ActType = act.ActType
		isSwap = true
	}
	if a.BeginAt != act.BeginAt {
		a.BeginAt = act.BeginAt
		isSwap = true
	}
	if a.EndAt != act.EndAt {
		a.EndAt = act.EndAt
		isSwap = true
	}
	return isSwap
}
