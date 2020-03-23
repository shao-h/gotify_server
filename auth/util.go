package auth

import (
	"gotify_server/model"

	"github.com/gin-gonic/gin"
)

//RegisterAuthentication registers the userid, user and token
func RegisterAuthentication(c *gin.Context, user *model.User, userID uint, tokenID string) {
	c.Set("user", user)
	c.Set("userid", userID)
	c.Set("tokenid", tokenID)
}

//GetUserID returns user id which was previously registered by RegisterAuthentication
func GetUserID(c *gin.Context) uint {
	user := c.MustGet("user").(*model.User)
	if user != nil {
		return user.ID
	}
	userID := c.MustGet("userid").(uint)
	if userID == 0 {
		panic("both token and user are null")
	}
	return userID
}

//GetTokenID returns the token id
func GetTokenID(c *gin.Context) string {
	return c.MustGet("tokenid").(string)
}
