package conf

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	Config *Conf
)

func Confinit(configPath string) {
	Config = getConf(configPath)
	log.Println("[Setting] Config init success")
}

type Conf struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	App      App      `yaml:"app"`
}

type Server struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
}

type App struct {
	RunMode     string `yaml:"run-mode"`
	PageSize    int    `yaml:"page-size"`
	IdentityKey string `yaml:"identity-key"`
	LogPath     string `yaml:"log-path"`
	AesKey      string `yaml:"aes-key"`
	FilePath    string `yaml:"file-path"`
}

type Database struct {
	Type        string `yaml:"type"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Name        string `yaml:"name"`
	SqlOutput   bool   `json:"sql-output"`
	MaxIdleConn int    `json:"max-idle-conn"`
	MaxOpenConn int    `json:"max-open-conn"`
}

func getConf(configPath string) *Conf {
	var c *Conf
	// todo 做成传入式
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("[Setting] config error: %v", err)
	}
	err = yaml.UnmarshalStrict(file, &c)
	if err != nil {
		log.Fatalf("[Setting] yaml unmarshal error: %v", err)
	}
	return c
}
