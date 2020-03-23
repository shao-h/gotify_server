package main

import (
	"fmt"
	"gotify_server/config"
	"gotify_server/database"
	"gotify_server/mode"
	"gotify_server/model"
	"gotify_server/router"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	// Version the version of Gotify.
	Version = "0.1"
	// Commit the git commit hash of this version.
	Commit = "unknown"
	// BuildDate the date on which this binary was build.
	BuildDate = "2020-03-18"
	// Mode the build mode
	Mode = mode.Dev
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	mode.Set(Mode)
	vInfo := &model.VersionInfo{Version: Version, Commit: Commit, BuildDate: BuildDate}
	fmt.Println("Starting Gotify version", vInfo.Version+"@"+BuildDate)

	rand.Seed(time.Now().UnixNano())

	conf := config.Get()
	if err := os.MkdirAll(conf.UploadedImagesDir, 0755); err != nil {
		log.Fatal(err)
	}

	db, err := database.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	engine := router.Create(db, vInfo, conf)
	addr := fmt.Sprintf("%s:%d", conf.Server.ListenAddr, conf.Server.Port)
	engine.Run(addr)
}
