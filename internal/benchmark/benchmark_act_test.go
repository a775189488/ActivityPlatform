package benchmark

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
)

const (
	ActivitySize = 1000000
)

func BenchmarkListActivity(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := 200
	b.SetParallelism(parallelism)

	defer fmt.Printf("benchmark list activity parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			param := "?page=1&size=10"
			reqUrl := httpServerAddr + "/act" + param
			res, err := client.Get(reqUrl)
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

func BenchmarkGetActivity(b *testing.B) {
	initForBenchmark(false)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := 200
	b.SetParallelism(parallelism)

	defer fmt.Printf("benchmark get activity parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			client := obj.(*httpClient).client
			actId := rand.Intn(ActivitySize) + 1
			reqUrl := fmt.Sprintf("%s/act/%d", httpServerAddr, actId)
			res, err := client.Get(reqUrl)
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

func BenchmarkGetActivityOnAuth(b *testing.B) {
	initForBenchmark(true)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := 200
	b.SetParallelism(parallelism)

	defer fmt.Printf("benchmark get activity on auth parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			mCli := obj.(*httpClient)
			client := mCli.client
			actId := rand.Intn(ActivitySize) + 1
			reqUrl := fmt.Sprintf("%s/v1/act/%d", httpServerAddr, actId)
			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}

			req.Header.Add(JwtHeaderKey, JwtHeaderValuePrefix+mCli.token)
			res, err := client.Do(req)
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

func BenchmarkListActivityOnAuth(b *testing.B) {
	initForBenchmark(true)
	defer destroyHttpClients()

	ctx := context.Background()
	b.ResetTimer()
	parallelism := 200
	b.SetParallelism(parallelism)

	defer fmt.Printf("benchmark list activity on auth parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			obj, err := clients.clientPool.BorrowObject(ctx)
			if err != nil {
				b.Error(err)
				continue
			}
			mCli := obj.(*httpClient)
			client := mCli.client
			param := "?page=1&size=10"
			reqUrl := httpServerAddr + "/v1/act" + param
			req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
			if err != nil {
				b.Error(err)
				_ = clients.clientPool.ReturnObject(ctx, obj)
				continue
			}

			req.Header.Add(JwtHeaderKey, JwtHeaderValuePrefix+mCli.token)
			res, err := client.Do(req)
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
