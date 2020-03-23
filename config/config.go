package config

import (
	"log"

	"github.com/jinzhu/configor"
)

//Configuration is stuff that can be configured
//externally per env variables or config file (config.yml).
type Configuration struct {
	Server struct {
		ListenAddr string `default:""`
		Port       int    `default:"80"`
		WebSocket  struct {
			PingPeriod  int
			PongTimeout int
		}
	}
	Database struct {
		Dialect    string
		Connection string
	}
	DefaultUser struct {
		Name string
		Pass string
	}
	PassStrength      int
	UploadedImagesDir string `default:"data/images"`
	PluginsDir        string
}

//Get returns the configuration extracted from env variables or config file.
func Get() *Configuration {
	conf := new(Configuration)
	err := configor.Load(conf, "config.yml")
	if err != nil {
		log.Fatal(err)
	}
	// addTrailingSlashToPaths(conf)
	return conf
}
