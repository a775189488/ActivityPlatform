package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"entrytask/internal/model"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
)

const (
	httpServerAddr = "http://127.0.0.1:8080"

	userPassword = "PNlwSCDOLpqb6tW40QSbNA=="
)

type httpClient struct {
	client *http.Client
}

type httpClientFactory struct {
	login bool
}

func (f *httpClientFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	return pool.NewPooledObject(
			&httpClient{
				client: initHttpClients(f.login),
			}),
		nil
}

func (f *httpClientFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	// do destroy
	myObj := object.Object.(*httpClient)
	logrus.Debugf("sessoin in poll destroyed, ctx:%v", ctx)
	myObj.client.CloseIdleConnections()
	return nil
}

func (f *httpClientFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	// do validate
	logrus.Debugf("sessoin in pool destroyed, ctx:%v", ctx)
	return true
}

func (f *httpClientFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	logrus.Debugf("session in pool activate, ctx:%v", ctx)
	return nil
}

func (f *httpClientFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do passivate(put into idle list)
	logrus.Debugf("session in pool passivate, ctx:%v", ctx)
	return nil
}

func initHttpClients(login bool) *http.Client {
	i := rand.Intn(userSize)
	client := getClient()
	// 登录
	if login {
		clientLogin(client, users[i])
	}
	return client
}

const (
	MaxConnsPerHost     int = 1
	MaxIdleConns        int = 0
	MaxIdleConnsPerHost int = 0
)

func getClient() *http.Client {
	cookieJar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   1 * time.Second,
				KeepAlive: 90 * time.Second,
			}).DialContext,
			MaxConnsPerHost:     MaxConnsPerHost,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		},
		Jar: cookieJar,
	}
	return client
}

type UserLoginReq struct {
	Username string `json:"username";`
	Password string `json:"password";`
}

func clientLogin(client *http.Client, u *model.User) {
	userLoginReq := UserLoginReq{
		Username: u.Username,
		Password: userPassword,
	}
	data, errData := json.Marshal(userLoginReq)
	if errData != nil {
		logrus.Panicf("json error:%v", u)
		return
	}
	var err error
	reqUrl := httpServerAddr + "/login"
	//req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(data))
	res, err := client.Post(reqUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logrus.Errorf("login request close err:%v", err)
	}
	body, errBody := ioutil.ReadAll(res.Body)
	if errBody != nil {
		logrus.Errorf("get body err:%v", err)
		return
	}
	err = res.Body.Close()
	if err != nil {
		logrus.Errorf("close body err:%v", err)
	}
	logrus.Debugf("post login response:%v", string(body[:]))
}

const (
	dbUser = "root"
	dbPwd  = "Admin@123"
	dbHost = "127.0.0.1:3306"
	dbName = "entrytask_activity_platform_db"
)

func initDb() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPwd, dbHost, dbName))
	if err != nil {
		log.Fatal("connecting mysql error: ", err)
	}
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	db.LogMode(false)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db
}
