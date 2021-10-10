package repository

import (
	"entrytask/internal/model"
	"testing"
)

var activityUserRepo ActivityUserRepo

func TestActivityUserRepo_ListActivityUser(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  RandString(10),
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}
	if userRepo.InsertUser(userObj) != nil {
		t.Fatal("insert user fail")
	}
	t.Cleanup(func() {
		if err := userRepo.DeleteUser(userObj.Id); err != nil {
			t.Fatalf("delete user(%v) fail", userObj)
		}
	})

	var actList []*model.Activity
	for i := 0; i < 15; i++ {
		act := &model.Activity{
			Title:       RandString(10),
			BeginAt:     111111,
			EndAt:       22222,
			Address:     RandString(10),
			Description: RandString(10),
			Creator:     1,
			ActType:     1,
		}
		actList = append(actList, act)
		if activityRepo.InsertActivity(act) != nil {
			t.Fatalf("insert act(%v) fail", *act)
		}
		actUserObj := &model.ActivityUser{ActId: act.Id, UserId: userObj.Id}
		if err := activityUserRepo.InsertActivityUser(actUserObj); err != nil {
			t.Fatalf("insert act user(%v) fail, err: %v", actUserObj, err)
		}
	}
	t.Cleanup(func() {
		if err := activityUserRepo.DeleteActivityUserByUserId(userObj.Id); err != nil {
			t.Fatalf("delete act user(%d) fail, err: %v", userObj.Id, err)
		}
	})
	t.Cleanup(func() {
		for _, a := range actList {
			if activityRepo.DeleteActivity(a.Id) != nil {
				t.Fatalf("delete act(%v) fail", a.Id)
			}
		}
	})

	var actIds []uint64
	for _, a := range actList {
		actIds = append(actIds, a.Id)
	}
	result, err := activityUserRepo.ListActivityUser(userObj.Id, actIds)
	if err != nil {
		t.Fatalf("list activity user fail, err: %v", err)
	}
	if len(result) != len(actList) {
		t.Fatalf("hope for %d result, actual %d", len(actList), len(result))
	}
}
