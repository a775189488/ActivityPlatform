package main

import (
	db2 "entrytask/internal/common/db"
	"entrytask/internal/common/logger"
	"entrytask/internal/common/utils"
	conf2 "entrytask/internal/conf"
	"entrytask/internal/model"
	"entrytask/internal/repository"
	"entrytask/internal/service"
	"flag"
	"github.com/facebookgo/inject"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	TestUserPasswordStr = "PNlwSCDOLpqb6tW40QSbNA=="
)

var configPath string
var objNum int

// 1 user 2 activityType 3 activity
var actionType int

const (
	UserAction       = 1
	ActTypeAction    = 2
	ActAction        = 3
	ActCommentAction = 4
)

func initCmdLineFlag() {
	flag.StringVar(&configPath, "config", "../etc/conf.yml", "configuration file")
	flag.StringVar(&configPath, "c", "../etc/conf.yml", "configuration file")

	flag.IntVar(&actionType, "type", 0, "test user num")

	flag.IntVar(&objNum, "num", 0, "test num")
	flag.Parse()
}

func toolInit() {
	initCmdLineFlag()
	conf2.Confinit(configPath)
}

func genSigleTestActivity() *model.Activity {
	return &model.Activity{
		Title:       utils.RandString(20),
		Description: utils.RandString(20),
		Address:     utils.RandString(10),
		ActType:     uint64(rand.Intn(10) + 1),
		BeginAt:     uint64(rand.Intn(100000) + 1),
		EndAt:       uint64(rand.Intn(100000) + 1),
	}
}

func genSigleTestUser() *model.User {
	return &model.User{
		Aliasname: utils.RandString(10),
		Username:  utils.RandString(10),
		Password:  TestUserPasswordStr,
		Email:     utils.RandString(10),
		Headpic:   utils.RandString(10),
		Role:      0,
	}
}

func genSigleTestActType() *model.ActivityType {
	return &model.ActivityType{
		Name:   utils.RandString(10),
		Parent: 0,
	}
}

func genSigleTestActComment() *model.ActivityComment {
	return &model.ActivityComment{
		UserId:  22,
		ActId:   12,
		Message: utils.RandString(10),
		Parent:  0,
	}
}

func createParallel(num int, createFun func(num int)) {
	var begin = time.Now()
	cpuCount := runtime.NumCPU()
	grRunCount := num / cpuCount
	if num%cpuCount != 0 {
		grRunCount += 1
	}
	var wg sync.WaitGroup
	for i := 0; i < cpuCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			createFun(grRunCount)
		}()
	}
	wg.Wait()
	end := time.Since(begin)
	userService.Log.Infof("[TestTool] create %d obj, create func(%v), use time %v", grRunCount*cpuCount, createFun, end)
}

func createTestAct(num int) {
	var failCount, successCount int
	for i := 0; i < num; i++ {
		tAct := genSigleTestActivity()
		if err := actService.CreateActivity(tAct); err != nil {
			failCount++
			userService.Log.Errorf("[TestTool] create act(%v) fail, err: %v", err)
		} else {
			successCount++
		}
	}
}

func createTestUser(num int) {
	var failCount, successCount int
	for i := 0; i < num; i++ {
		tUser := genSigleTestUser()
		if err := userService.CreateUser(tUser); err != nil {
			failCount++
			userService.Log.Errorf("[TestTool] create user(%v) fail, err: %v", err)
		} else {
			successCount++
		}
	}
}

func createTestActionType(num int) {
	var failCount, successCount int
	for i := 0; i < num; i++ {
		tUser := genSigleTestActType()
		if err := actTypeService.CreateActType(tUser); err != nil {
			failCount++
			userService.Log.Errorf("[TestTool] create activity type(%v) fail, err: %v", err)
		} else {
			successCount++
		}
	}
}

func createTestActionComment(num int) {
	var failCount, successCount int
	for i := 0; i < num; i++ {
		ac := genSigleTestActComment()
		if err := actCommentService.CreateActComment(ac); err != nil {
			failCount++
			userService.Log.Errorf("[TestTool] create activity comment(%v) fail, err: %v", err)
		} else {
			successCount++
		}
	}
}

var userService service.UserService
var actTypeService service.ActTypeService
var actCommentService service.ActCommentService
var actService service.ActService

func main() {
	toolInit()

	userService = service.UserService{}
	actTypeService = service.ActTypeService{}
	actCommentService = service.ActCommentService{}
	actService = service.ActService{}
	db := db2.DbMysql{}
	zap := logger.Logger{}
	var injector inject.Graph
	if err := injector.Provide(
		&inject.Object{Value: &db},
		&inject.Object{Value: &zap},
		&inject.Object{Value: &userService},
		&inject.Object{Value: &actTypeService},
		&inject.Object{Value: &actCommentService},
		&inject.Object{Value: &actService},
		&inject.Object{Value: &repository.UserRepo{}},
		&inject.Object{Value: &repository.ActivityRepo{}},
		&inject.Object{Value: &repository.ActivityUserRepo{}},
		&inject.Object{Value: &repository.ActivityTypeRepo{}},
		&inject.Object{Value: &repository.ActivityCommentRepo{}},
	); err != nil {
		log.Fatal("inject fatal: ", err)
	}
	if err := injector.Populate(); err != nil {
		log.Fatal("injector fatal: ", err)
	}
	zap.Init()
	if err := db.Connect(); err != nil {
		log.Fatal("db fatal:", err)
	}

	switch actionType {
	case UserAction:
		createParallel(objNum, createTestUser)
	case ActTypeAction:
		createParallel(objNum, createTestActionType)
	case ActAction:
		createParallel(objNum, createTestAct)
	case ActCommentAction:
		createParallel(objNum, createTestActionComment)
	default:
		log.Fatalf("invalid action type %d", actionType)
	}
}
