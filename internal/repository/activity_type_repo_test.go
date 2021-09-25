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
		if activityTypeRepo.DeleteActivityType(actType.Id) == false {
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
