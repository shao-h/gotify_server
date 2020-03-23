package api

import (
	"errors"
	"fmt"
	"gotify_server/auth"
	"gotify_server/model"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
)

//ApplicationDatabase interface for encapsulating database access.
type ApplicationDatabase interface {
	GetApplicationsByUser(uint) ([]*model.Application, error)
	CreateApplication(*model.Application) error
	GetApplicationByID(id uint) (*model.Application, error)
	GetApplicationByToken(string) (*model.Application, error)
	UpdateApplicationField(*model.Application, map[string]interface{}) error
}

//ApplicationAPI provides handlers for managing application.
type ApplicationAPI struct {
	DB       ApplicationDatabase
	ImageDir string
}

//GetApplications returns all applications a user has
func (a *ApplicationAPI) GetApplications(c *gin.Context) {
	userID := auth.GetUserID(c)
	apps, err := a.DB.GetApplicationsByUser(userID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	for _, app := range apps {
		withResovedImage(app)
	}
	c.JSON(200, apps)
}

//CreateApplication creates an application and returns the access token
func (a *ApplicationAPI) CreateApplication(c *gin.Context) {
	app := new(model.Application)
	if err := c.Bind(app); err == nil {
		app.Token, err = a.genAppToken()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		app.UserID = auth.GetUserID(c)
		err = a.DB.CreateApplication(app)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, withResovedImage(app))
	}
}

//UploadApplicationImage uploads an image for an application.
func (a *ApplicationAPI) UploadApplicationImage(c *gin.Context) {
	withParam(c, "id", func(id uint) {
		app, err := a.DB.GetApplicationByID(id)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		if app != nil && app.UserID == auth.GetUserID(c) {
			file, err := c.FormFile("file")
			if err == http.ErrMissingFile {
				c.AbortWithError(400, errors.New("file with key 'file' must be present"))
				return
			} else if err != nil {
				c.AbortWithError(500, err)
				return
			}
			head := make([]byte, 261)
			f, _ := file.Open()
			f.Read(head)
			if !filetype.IsImage(head) {
				c.AbortWithError(400, errors.New("file must be an image"))
				return
			}
			ext := filepath.Ext(file.Filename)
			imgName := a.genImgName(a.ImageDir, func() string {
				return generateImageName() + ext
			})
			err = c.SaveUploadedFile(file, filepath.Join(a.ImageDir, imgName))
			if err != nil {
				c.AbortWithError(500, err)
				return
			}

			if app.Image != "" {
				os.Remove(filepath.Join(a.ImageDir, app.Image))
			}
			err = a.DB.UpdateApplicationField(app, map[string]interface{}{"Image": imgName})
			if err != nil {
				c.AbortWithError(500, err)
				return
			}
			c.JSON(200, withResovedImage(app))
		} else {
			c.AbortWithError(400, fmt.Errorf("app with id %d not exists", id))
		}
	})
}

func withResovedImage(app *model.Application) *model.Application {
	if app.Image == "" {
		app.Image = "static/defaultapp.png"
	} else {
		app.Image = "image/" + app.Image
	}
	return app
}

//if it takes too long?
func (*ApplicationAPI) genImgName(imgDir string, gen func() string) string {
	for {
		name := gen()
		fileName := filepath.Join(imgDir, name)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			return name
		}
	}
}

func (a *ApplicationAPI) genAppToken() (string, error) {
	for {
		tokenID := generateApplicationToken()
		app, err := a.DB.GetApplicationByToken(tokenID)
		if err != nil {
			return "", err
		}
		if app == nil {
			return tokenID, nil
		}
	}
}
