package api

import (
	"encoding/json"
	"fmt"
	"gotify_server/auth"
	"gotify_server/model"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//MessageDatabase for encapsulating database access.
type MessageDatabase interface {
	GetApplicationByToken(string) (*model.Application, error)
	CreateMessage(*model.Message) error
}

//Notifier notifies when new message received.
type Notifier interface {
	Notify(userID uint, message *model.MessageExternal)
}

//MessageAPI provides handlers for managing messages.
type MessageAPI struct {
	DB       MessageDatabase
	Notifier Notifier
}

//CreateMessage creates a message, application token is required
func (m *MessageAPI) CreateMessage(c *gin.Context) {
	msg := new(model.MessageExternal)
	if err := c.Bind(msg); err == nil {
		application, err := m.DB.GetApplicationByToken(auth.GetTokenID(c))
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		if application == nil {
			c.AbortWithError(400, fmt.Errorf("app with token %s not exists", auth.GetTokenID(c)))
			return
		}
		msg.ApplicationID = application.ID
		if strings.TrimSpace(msg.Title) == "" {
			msg.Title = application.Name
		}
		msg.Created = time.Now()
		msgInternal := toInternalMessage(msg)
		err = m.DB.CreateMessage(msgInternal)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		m.Notifier.Notify(auth.GetUserID(c), msg)
		c.JSON(200, msg)
	}
}

func toInternalMessage(msg *model.MessageExternal) *model.Message {
	res := &model.Message{
		ID:            msg.ID,
		ApplicationID: msg.ApplicationID,
		Message:       msg.Message,
		Title:         msg.Title,
		Priority:      msg.Priority,
		Created:       msg.Created,
	}
	if msg.Extras != nil {
		res.Extras, _ = json.Marshal(msg.Extras)
	}
	return res
}
