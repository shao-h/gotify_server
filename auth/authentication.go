package auth

import (
	"errors"
	"gotify_server/auth/password"
	"gotify_server/model"

	"github.com/gin-gonic/gin"
)

const (
	headerName = "X-Gotify-Key"
)

//Database interface for encapsulating database access.
type Database interface {
	GetApplicationByToken(string) (*model.Application, error)
	GetUserByName(name string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetClientByToken(token string) (*model.Client, error)
}

//Auth is the provider for authentication middleware
type Auth struct {
	DB Database
}

type authenticate func(tokenID string, user *model.User) (authenticated, success bool, userID uint, err error)

//RequireAdmin returns a gin middleware which requires a client token or basic authentication header to be
//supplied with the request. Also the authenticated user must be an administrator.
func (a *Auth) RequireAdmin() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, bool, uint, error) {
		if user != nil {
			return true, user.Admin, user.ID, nil
		}
		token, err := a.DB.GetClientByToken(tokenID)
		if err != nil {
			return false, false, 0, err
		} else if token != nil {
			usr, err := a.DB.GetUserByID(token.UserID)
			if err != nil {
				return false, false, token.UserID, err
			} else if usr != nil {
				return true, usr.Admin, token.UserID, nil
			}
		}
		return false, false, 0, nil
	})
}

//RequireClient returns a gin middleware which requires a client token or basic authenticaton header to be
//supplied with the request.
func (a *Auth) RequireClient() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, bool, uint, error) {
		if user != nil {
			return true, true, user.ID, nil
		}
		client, err := a.DB.GetClientByToken(tokenID)
		if err != nil {
			return false, false, 0, err
		} else if client != nil {
			return true, true, client.UserID, nil
		}
		return false, false, 0, nil
	})
}

//RequireAppToken returns a gin middleware which requires an application token to be supplied.
func (a *Auth) RequireAppToken() gin.HandlerFunc {
	return a.requireToken(func(tokenID string, user *model.User) (bool, bool, uint, error) {
		if user != nil {
			return true, false, 0, nil
		}
		app, err := a.DB.GetApplicationByToken(tokenID)
		if err != nil {
			return false, false, 0, err
		} else if app != nil {
			return true, true, app.UserID, nil
		}
		return false, false, 0, nil
	})
}

func (a *Auth) requireToken(auth authenticate) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := a.tokenFromQueryOrHeader(c)
		user, err := a.userFromBasicAuth(c)

		if err != nil {
			c.AbortWithError(500, errors.New("an error occured while authenticating user"))
			return
		}

		if token != "" || user != nil {
			authenticated, ok, userID, err := auth(token, user)
			if err != nil {
				c.AbortWithError(500, errors.New("an error occured while authenticating user"))
				return
			} else if ok {
				//an administrator
				RegisterAuthentication(c, user, userID, token)
				c.Next()
				return
			} else if authenticated {
				c.AbortWithError(403, errors.New("you are not allowed to access this api"))
				return
			}
		}
		c.AbortWithError(401, errors.New("you need to provide a valid access token or user credentials to access this api"))
	}
}

func (a *Auth) tokenFromQueryOrHeader(c *gin.Context) (token string) {
	token = c.Query("token")
	if token != "" {
		return
	}
	return c.GetHeader(headerName)
}

func (a *Auth) userFromBasicAuth(c *gin.Context) (*model.User, error) {
	name, pwd, ok := c.Request.BasicAuth()
	if ok {
		user, err := a.DB.GetUserByName(name)
		if err != nil {
			return nil, err
		}
		if user != nil && password.ComparePassword(user.Pass, []byte(pwd)) {
			return user, nil
		}
	}
	return nil, nil
}
