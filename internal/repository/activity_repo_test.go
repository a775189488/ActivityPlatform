package repository

import (
	"entrytask/internal/model"
	"testing"
)

var activityRepo ActivityRepo

func TestActivityRepo_InsertActivity(t *testing.T) {
	act := &model.Activity{
		Title:       "test",
		BeginAt:     111111,
		EndAt:       22222,
		Address:     "ShangHai",
		Description: "xxx",
		Creator:     1,
		ActType:     1,
	}

	if activityRepo.InsertActivity(act) == false {
		t.Fatalf("insert act(%v) fail", *act)
	}
	t.Cleanup(func() {
		if activityRepo.DeleteActivity(act.Id) == false {
			t.Fatalf("delete act(%d) fail", act.Id)
		}
	})
	newObj := activityRepo.GetActivityById(act.Id)
	if newObj.Title != act.Title {
		t.Fatalf("need (%s), actual (%s)", act.Title, newObj.Title)
	}
}

func TestActivityRepo_GetActivityDetail(t *testing.T) {
	actType := &model.ActivityType{
		Name:   "test_tt",
		Parent: 0,
	}
	if err := activityTypeRepo.InsertActivityType(actType); err != nil {
		t.Fatalf("create activity type faile, err: %v", err)
	}
	t.Cleanup(func() {
		activityTypeRepo.DeleteActivityType(actType.Id)
	})
	act := &model.Activity{
		Title:       "test",
		BeginAt:     111111,
		EndAt:       22222,
		Address:     "ShangHai",
		Description: "xxx",
		Creator:     1,
		ActType:     actType.Id,
	}
	if activityRepo.InsertActivity(act) == false {
		t.Fatalf("insert act(%v) fail", *act)
	}
	t.Cleanup(func() {
		if activityRepo.DeleteActivity(act.Id) == false {
			t.Fatalf("delete act(%d) fail", act.Id)
		}
	})
	newObj, err := activityRepo.GetActivityDetail(act.Id)
	if err != nil {
		t.Fatalf("get act(%d) fail, err: %v", act.Id, err)
	}
	if newObj.ActTypeName != actType.Name {
		t.Fatalf("need (%s), actual (%s), obj(%v)", actType.Name, newObj.ActTypeName, *newObj)
	}
}

func TestActivityRepo_GetUserByActivityId(t *testing.T) {
	act := &model.Activity{
		Title:       "test",
		BeginAt:     111111,
		EndAt:       22222,
		Address:     "ShangHai",
		Description: "xxx",
		Creator:     1,
		ActType:     1,
	}
	if activityRepo.InsertActivity(act) == false {
		t.Fatalf("insert act(%v) fail", *act)
	}
	t.Cleanup(func() {
		if activityRepo.DeleteActivity(act.Id) == false {
			t.Fatalf("delete act(%d) fail", act.Id)
		}
	})

	var userList []*model.User
	for i := 0; i < 15; i++ {
		obj := &model.User{
			Aliasname: RandString(5),
			Username:  RandString(5),
			Password:  RandString(5),
			Email:     RandString(5),
			Headpic:   RandString(5),
			Role:      0,
		}
		userList = append(userList, obj)
		if err := userRepo.InsertUser(obj); err != nil {
			t.Fatalf("insert user(%v) fail, err: %v", obj, err)
		}
		actUserObj := &model.ActivityUser{ActId: act.Id, UserId: obj.Id}
		if err := activityUserRepo.InsertActivityUser(actUserObj); err != nil {
			t.Fatalf("insert act user(%v) fail, err: %v", actUserObj, err)
		}
	}
	t.Cleanup(func() {
		for _, u := range userList {
			if err := userRepo.DeleteUser(u.Id); err != nil {
				t.Fatalf("delete user(%v) fail", u)
			}
		}
	})
	t.Cleanup(func() {
		if err := activityUserRepo.DeleteActivityUserByActId(act.Id); err != nil {
			t.Fatalf("delete act user(%d) fail, err: %v", act.Id, err)
		}
	})

	var total int32
	actUsers, err := activityRepo.GetUserByActivityId(act.Id, 1, 10, &total)
	if err != nil {
		t.Fatalf("get user by activity(%d) fail, err: %v", act.Id, err)
	}
	if int(total) != len(userList) {
		t.Fatalf("hope for %d users, actual %d", len(userList), total)
	}
	if len(actUsers) != 10 {
		t.Fatalf("hope for %d users, actual %d", 10, len(actUsers))
	}
}
