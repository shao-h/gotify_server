package mode

import (
	"log"

	"github.com/gin-gonic/gin"
)

const (
	//Dev for development mode
	Dev = "dev"
	//Prod for production mode
	Prod = "prod"
)

var mode = Dev

//Set sets the new mode
func Set(newMode string) {
	mode = newMode
	updateGinMode()
}

//Get returns the current mode
func Get() string {
	return mode
}

func updateGinMode() {
	switch mode {
	case Dev:
		gin.SetMode(gin.DebugMode)
	case Prod:
		gin.SetMode(gin.ReleaseMode)
	default:
		log.Fatal("unknown mode")
	}
}
