package benchmark

import (
	"bytes"
	"context"
	"encoding/json"
	"entrytask/internal/model"
	"fmt"
	"github.com/jinzhu/gorm"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
)

const (
	userSize   = 10000
	clientSize = 200
)

var (
	users   []*model.User
	clients *HttpClientPool
)

func GetUserUserName(num int, conn *gorm.DB) error {
	return conn.Table(model.User{}.TableName()).Select("id,username").Limit(num).Find(&users).Error
}

func initUsers() error {
	users = make([]*model.User, 0)
	conn := initDb()
	if err := GetUserUserName(userSize, conn); err != nil {
		return err
	}
	return nil
}

type HttpClientPool struct {
	clientPool *pool.ObjectPool
}

func NewHttpClientPool(login bool, size int) *HttpClientPool {
	ctx := context.Background()
	config := pool.ObjectPoolConfig{
		MaxTotal:           size,
		MaxIdle:            size,
		BlockWhenExhausted: true,
	}
	return &HttpClientPool{
		clientPool: pool.NewObjectPool(ctx, &httpClientFactory{login: login}, &config),
	}
}

func destroyHttpClients() {
	ctx := context.Background()
	clients.clientPool.Close(ctx)
}

func initForBenchmark(login bool) {
	clients = NewHttpClientPool(login, clientSize)
	users = make([]*model.User, 0)

	initUsers()
	initHttpClients(login)
}

func BenchmarkLogin(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)

	defer fmt.Printf("benchmark login parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			iu := rand.Intn(userSize)
			u := users[iu]
			userLoginReq := UserLoginReq{
				Username: u.Username,
				Password: userPassword,
			}
			data, errData := json.Marshal(userLoginReq)
			if errData != nil {
				logrus.Panicf("json error:%v", u)
				b.Skipped()
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}
			reqUrl := httpServerAddr + "/login"
			res, err := client.Post(reqUrl, "application/json", bytes.NewBuffer(data))
			if err != nil {
				logrus.Errorf("login request close err:%v", err)
			}
			if res.StatusCode != http.StatusOK {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					_ = clients.clientPool.ReturnObject(ctx, obj)
					continue
				}
			} else {
				_, err := io.Copy(ioutil.Discard, res.Body)
				err = res.Body.Close()
				if err != nil {
					b.Error(err)
					logrus.Errorf("close body err:%v", err)
				}
			}
			_ = clients.clientPool.ReturnObject(ctx, obj)
		}
	})
}
