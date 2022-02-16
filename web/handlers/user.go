package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/web/common"
)

// Register
func PostUserHandler(c *gin.Context) {
	payload := dto.CreateUserPayload{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, err)
		return
	}

	user, err := bll.CreateUser(payload)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
