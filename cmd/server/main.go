package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	conf2 "entrytask/internal/conf"
	"entrytask/internal/routers"
)

var configPath string

func initCmdLineFlag() {
	flag.StringVar(&configPath, "config", "../etc/conf.yml", "configuration file")
	flag.StringVar(&configPath, "c", "../etc/conf.yml", "configuration file")
	flag.Parse()
}

func main() {
	initCmdLineFlag()
	conf2.Confinit(configPath)
	router := routers.InitRouter()
	conf := conf2.Config.Server
	s := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        router,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
