package repository

import (
	"os"
	"testing"

	db2 "entrytask/internal/common/db"
	"entrytask/internal/common/logger"
	"entrytask/internal/conf"
)

func TestMain(m *testing.M) {
	conf.Config = &conf.Conf{
		conf.Server{},
		conf.Database{
			Type:     "mysql",
			User:     "root",
			Password: "Admin@123",
			Host:     "127.0.0.1:3306",
			Name:     "entrytask_activity_platform_db",
		},
		conf.App{
			LogPath: "",
		},
	}
	log := logger.Logger{}
	db := db2.DbMysql{}
	log.Init()
	if err := db.Connect(); err != nil {
		return
	}
	base := BaseRepo{&db, &log}
	userRepo = UserRepo{&log, base}
	activityRepo = ActivityRepo{&log, base}
	activityTypeRepo = ActivityTypeRepo{&log, base}
	activityCommentRepo = ActivityCommentRepo{&log, base}
	activityUserRepo = ActivityUserRepo{&log, base}

	os.Exit(m.Run())
}
