package repository

import (
	"math/rand"
	"testing"
	"time"

	"entrytask/internal/model"
)

var userRepo UserRepo

func TestUserRepo_InsertUser(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  "iam",
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}

	if err := userRepo.InsertUser(userObj); err != nil {
		t.Fatal("insert user fail")
	}

	newObj, err := userRepo.GetUserById(userObj.Id)
	if err != nil {
		t.Fatal("find user obj nil!")
	}
	if newObj.Username != userObj.Username {
		t.Fatalf("need username(%s), actually(%s)", userObj.Username, newObj.Username)
	}

	t.Cleanup(func() {
		if err := userRepo.DeleteUser(userObj.Id); err != nil {
			t.Fatalf("delete user(%v) fail", userObj)
		}
	})
}

func TestUserRepo_GetUserByUsername(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  "iam",
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}

	if err := userRepo.InsertUser(userObj); err != nil {
		t.Fatal("insert user fail")
	}

	newObj, err := userRepo.GetUserByUsername(userObj.Username)
	if err != nil {
		t.Fatal("find user obj nil!")
	}
	if newObj.Id != userObj.Id {
		t.Fatalf("need username(%d), actually(%d)", userObj.Id, newObj.Id)
	}

	t.Cleanup(func() {
		if err := userRepo.DeleteUser(userObj.Id); err != nil {
			t.Fatalf("delete user(%v) fail", userObj)
		}
	})
}

func TestUserRepo_DeleteUser(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  "iam",
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}

	if err := userRepo.InsertUser(userObj); err != nil {
		t.Fatal("insert user fail")
	}

	if err := userRepo.DeleteUser(userObj.Id); err != nil {
		t.Fatalf("delete user(%v) fail", userObj)
	}

	_, err := userRepo.GetUserById(userObj.Id)
	if err == nil {
		t.Fatal("find user obj fail!")
	}
}

func TestUserRepo_ListUser(t *testing.T) {
	var userList []*model.User
	for i := 0; i < 100; i++ {
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
			t.Fatalf("insert user(%v) fail", obj)
		}
	}

	var total int32
	result, err := userRepo.ListUser(0, 50, &total, map[string]interface{}{})
	if err != nil {
		t.Fatalf("list user fail!")
	}
	if len(result) != 50 {
		t.Fatalf("need 50, actually %d", len(result))
	}
	if total != 100 {
		t.Fatalf("need 100, actually %d", total)
	}

	t.Cleanup(func() {
		for _, u := range userList {
			if err := userRepo.DeleteUser(u.Id); err != nil {
				t.Fatalf("delete user(%v) fail", u)
			}
		}
	})
}

func TestUserRepo_UpdateUser(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  "iam",
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}

	if err := userRepo.InsertUser(userObj); err != nil {
		t.Fatal("insert user fail")
	}
	userObj.Aliasname = "bbb"
	if err := userRepo.UpdateUser(userObj); err != nil {
		t.Fatalf("update user fail")
	}
	newObj, err := userRepo.GetUserById(userObj.Id)
	if err != nil {
		t.Fatalf("get user by id fail")
	}
	if newObj.Aliasname != userObj.Aliasname {
		t.Fatalf("need user.aliasname(%s), actual(%s)", userObj.Aliasname, newObj.Aliasname)
	}
	t.Cleanup(func() {
		if err := userRepo.DeleteUser(userObj.Id); err != nil {
			t.Fatalf("delete user(%v) fail", userObj)
		}
	})
}

func TestUserRepo_CheckUser(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  "iam",
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}
	if err := userRepo.InsertUser(userObj); err != nil {
		t.Fatal("insert user fail")
	}
	if userRepo.CheckUser(userObj.Username, userObj.Password) != nil {
		t.Fatalf("check user fail")
	}
	t.Cleanup(func() {
		if err := userRepo.DeleteUser(userObj.Id); err != nil {
			t.Fatalf("delete user(%v) fail", userObj)
		}
	})
}

func TestUserRepo_CheckUserPasswordFail(t *testing.T) {
	userObj := &model.User{
		Aliasname: "tony",
		Username:  "iam",
		Password:  "1234456",
		Email:     "123@qq.com",
		Headpic:   "/home/a.jpg",
		Role:      0,
	}
	if userRepo.InsertUser(userObj) != nil {
		t.Fatal("insert user fail")
	}
	if userRepo.CheckUser(userObj.Username, "test") == nil {
		t.Fatalf("check user fail")
	}
	t.Cleanup(func() {
		if err := userRepo.DeleteUser(userObj.Id); err != nil {
			t.Fatalf("delete user(%v) fail", userObj)
		}
	})
}

func TestUserRepo_ListUserActivity(t *testing.T) {
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

	var total int32
	acts, err := userRepo.ListUserActivity(1, 10, &total, userObj.Id)
	if err != nil {
		t.Fatalf("list user(%d) activity fail, err: %v", userObj.Id, err)
	}
	if int(total) != len(actList) {
		t.Fatalf("hope for %d users, actual %d", len(actList), total)
	}
	if len(acts) != 10 {
		t.Fatalf("hope for %d users, actual %d", 10, len(acts))
	}
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
