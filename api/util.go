package api

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func withParam(c *gin.Context, name string, f func(uint)) {
	id, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil {
		c.AbortWithError(400, errors.New("invalid id"))
	}
	f(uint(id))
}
