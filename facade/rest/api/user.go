package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/facade/rest/common"
)

// Register
func PostUserHandler(c *gin.Context) {
	payload := dto.NewUser{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, err)
		return
	}

	user, err := service.CreateUser(payload)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
