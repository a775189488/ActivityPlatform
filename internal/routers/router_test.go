package routers

import (
	"entrytask/internal/conf"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

	os.Exit(m.Run())
}

func TestInitRouter(t *testing.T) {
	engine := InitRouter()
	ts := httptest.NewServer(engine)
	// shutdown
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/ping", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected status code:200, got:%v", resp.StatusCode)
	}
}
