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
	UserAction    = 1
	ActTypeAction = 2
	ActAction     = 3
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

func createTestUser(num int) {
	var failCount, successCount int
	start := time.Now()

	for i := 0; i < num; i++ {
		tUser := genSigleTestUser()
		if err := userService.CreateUser(tUser); err != nil {
			failCount++
			userService.Log.Errorf("[TestTool] create user(%v) fail, err: %v", err)
		} else {
			successCount++
		}
	}

	elapsed := time.Since(start)
	userService.Log.Infof("[TestTool] hope for create %d user, actual success %d,"+
		" fail %d, waste time %v", num, successCount, failCount, elapsed)
}

func createTestActionType(num int) {
	var failCount, successCount int
	start := time.Now()

	for i := 0; i < num; i++ {
		tUser := genSigleTestActType()
		if err := actTypeService.CreateActType(tUser); err != nil {
			failCount++
			userService.Log.Errorf("[TestTool] create activity type(%v) fail, err: %v", err)
		} else {
			successCount++
		}
	}

	elapsed := time.Since(start)
	userService.Log.Infof("[TestTool] hope for create %d activity type, actual success %d,"+
		" fail %d, waste time %v", num, successCount, failCount, elapsed)
}

var userService service.UserService
var actTypeService service.ActTypeService

func main() {
	toolInit()

	userService = service.UserService{}
	actTypeService = service.ActTypeService{}
	db := db2.DbMysql{}
	zap := logger.Logger{}
	var injector inject.Graph
	if err := injector.Provide(
		&inject.Object{Value: &db},
		&inject.Object{Value: &zap},
		&inject.Object{Value: &userService},
		&inject.Object{Value: &actTypeService},
		&inject.Object{Value: &repository.UserRepo{}},
		&inject.Object{Value: &repository.ActivityTypeRepo{}},
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
		createTestUser(objNum)
	case ActTypeAction:
		createTestActionType(objNum)
	default:
		log.Fatalf("invalid action type %d", actionType)
	}
}
