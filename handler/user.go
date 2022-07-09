package handler

import (
	"crowdfund/helper"
	"crowdfund/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (u *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.ResponseValidationError(err)

		errorMessages := gin.H{"errors": errors}
		response := helper.JSONResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	newUser, err := u.userService.RegisterUser(input)
	if err != nil {
		response := helper.JSONResponse("Register account failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
	}
	formattedResponse := user.FormatterUserResponse(newUser, "tokentoken")
	response := helper.JSONResponse("Account has been created", 200, "success", formattedResponse)
	c.JSON(http.StatusOK, response)
}
