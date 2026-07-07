package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/koki/tabgraph/internal/api"
	"github.com/koki/tabgraph/internal/config"
	"github.com/koki/tabgraph/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	sqlDB, err := db.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer sqlDB.Close()

	router := api.NewServer(sqlDB)

	addr := "localhost:" + cfg.Port
	url := "http://" + addr
	fmt.Printf("tabgraph listening on %s\n", url)

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()

	log.Fatal(http.ListenAndServe(addr, router))
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return
	}
	exec.Command(cmd, args...).Start() //nolint:errcheck
}
