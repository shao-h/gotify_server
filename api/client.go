package api

import (
	"gotify_server/auth"
	"gotify_server/model"

	"github.com/gin-gonic/gin"
)

//ClientDatabase encapsulates database access
type ClientDatabase interface {
	CreateClient(*model.Client) error
	GetClientByToken(string) (*model.Client, error)
}

//ClientAPI provides handler for managing clients and applications
type ClientAPI struct {
	DB            ClientDatabase
	NotifyDeleted func(uint, string)
}

//CreateClient creates a client and returns the access token.
func (a *ClientAPI) CreateClient(c *gin.Context) {
	client := new(model.Client)
	if err := c.Bind(client); err == nil {
		client.Token, err = a.genClientToken()
		if err == nil {
			client.UserID = auth.GetUserID(c)
			err = a.DB.CreateClient(client)
			if err == nil {
				c.JSON(200, client)
				return
			}
		}
		c.AbortWithError(500, err)
	}
}

func (a *ClientAPI) genClientToken() (string, error) {
	for {
		tokenID := generateClientToken()
		client, err := a.DB.GetClientByToken(tokenID)
		if err != nil {
			return "", err
		}
		if client == nil {
			return tokenID, nil
		}
	}
}
