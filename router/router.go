package router

import (
	"gotify_server/api"
	"gotify_server/api/stream"
	"gotify_server/auth"
	"gotify_server/config"
	"gotify_server/database"
	"gotify_server/model"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

//Create creates the gin engine with all routes.
func Create(db *database.GormDatabase, vInfo *model.VersionInfo, conf *config.Configuration) *gin.Engine {
	r := gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 4 << 10

	authentication := auth.Auth{DB: db}
	clientHandler := api.ClientAPI{DB: db}
	applicationHandler := api.ApplicationAPI{
		DB:       db,
		ImageDir: conf.UploadedImagesDir,
	}
	pingPeriod := time.Duration(conf.Server.WebSocket.PingPeriod) * time.Second
	pongTimeout := time.Duration(conf.Server.WebSocket.PongTimeout) * time.Second
	streamHandler := stream.New(pingPeriod, pongTimeout)
	messageHandler := api.MessageAPI{Notifier: streamHandler, DB: db}

	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, vInfo)
	})

	// r.POST("/message", authentication.RequireAppToken(), messageHanler.CreateMessage)
	r.POST("/message", authentication.RequireAppToken(), messageHandler.CreateMessage)

	clientAuth := r.Group("")
	clientAuth.Use(authentication.RequireClient())
	{
		app := clientAuth.Group("/application")
		{
			app.GET("", applicationHandler.GetApplications)
			app.POST("", applicationHandler.CreateApplication)
			app.POST("/:id/image", applicationHandler.UploadApplicationImage)
		}

		client := clientAuth.Group("client")
		{
			client.POST("", clientHandler.CreateClient)
		}

		clientAuth.GET("/stream", streamHandler.Handle)
	}

	authAdmin := r.Group("/user")
	authAdmin.Use(authentication.RequireAdmin())
	{
		authAdmin.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			log.Println(name)
			user, err := db.GetUserByName(name)
			log.Println(user, err)
		})

	}

	return r
}
