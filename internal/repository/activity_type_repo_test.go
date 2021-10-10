package repository

import (
	"entrytask/internal/model"
	"testing"
)

var activityTypeRepo ActivityTypeRepo

func TestActivityTypeRepo_InsertActivityType(t *testing.T) {
	actType := &model.ActivityType{
		Name:   "test",
		Parent: 0,
	}

	if activityTypeRepo.InsertActivityType(actType) != nil {
		t.Fatalf("insert activity(%v) type faile", actType)
	}

	t.Cleanup(func() {
		if activityTypeRepo.DeleteActivityType(actType.Id) != nil {
			t.Fatalf("delete activity(%v) type faile", actType)
		}
	})

	newObj, err := activityTypeRepo.GetActivityTypeById(actType.Id)
	if err != nil {
		t.Fatalf("get activity(%d) type fail, err: %v", actType.Id, err)
	}
	if newObj.Name != actType.Name {
		t.Fatalf("need (%s), actual(%s)", actType.Name, newObj.Name)
	}
}

func TestActivityTypeRepo_GetActCountByActType(t *testing.T) {
	var actList []*model.Activity
	for i := 0; i < 15; i++ {
		act := &model.Activity{
			Title:       RandString(10),
			BeginAt:     111111,
			EndAt:       22222,
			Address:     RandString(10),
			Description: RandString(10),
			Creator:     1,
			ActType:     999,
		}
		actList = append(actList, act)
		if activityRepo.InsertActivity(act) != nil {
			t.Fatalf("insert act(%v) fail", *act)
		}
	}
	t.Cleanup(func() {
		for _, a := range actList {
			if activityRepo.DeleteActivity(a.Id) != nil {
				t.Fatalf("delete act(%v) fail", a.Id)
			}
		}
	})

	count, err := activityTypeRepo.GetActCountByActType(999)
	if err != nil {
		t.Fatalf("get activity count by activity type(999) fail, err: %v", err)
	}
	if count != len(actList) {
		t.Fatalf("hope for %d, actual %d", len(actList), count)
	}
}
