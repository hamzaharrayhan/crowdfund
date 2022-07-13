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
		return
	}

	newUser, err := u.userService.RegisterUser(input)
	if err != nil {
		response := helper.JSONResponse("Register account failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	formattedResponse := user.FormatterUserResponse(newUser, "tokentoken")
	response := helper.JSONResponse("Account has been created", 200, "success", formattedResponse)
	c.JSON(http.StatusOK, response)
}

func (u *userHandler) LoginHandler(c *gin.Context) {
	var input user.LoginUserInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.ResponseValidationError(err)
		errorMessages := gin.H{"errors": errors}
		response := helper.JSONResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessages)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := u.userService.UserLogin(input)
	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.JSONResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formattedResponse := user.FormatterUserResponse(loggedInUser, "tokentoken")
	response := helper.JSONResponse("Successfully Logged In", http.StatusOK, "success", formattedResponse)
	c.JSON(http.StatusOK, response)
}

func (u *userHandler) CheckEmailAvailability(c *gin.Context) {
	var emailInput user.EmailAvailabilityInput
	err := c.ShouldBindJSON(&emailInput)

	if err != nil {
		errors := helper.ResponseValidationError(err)
		errorMessages := gin.H{"errors": errors}
		response := helper.JSONResponse("Failed in checking email availability", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	isEmailAvailable, err := u.userService.EmailAvailability(emailInput)
	if err != nil {
		errors := helper.ResponseValidationError(err)
		errorMessages := gin.H{"errors": errors}
		response := helper.JSONResponse("Failed in checking email availability", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
	}

	metaMessage := "Email has been registered"
	data := gin.H{"is_available": isEmailAvailable}
	if isEmailAvailable {
		metaMessage = "Email is Available"
	}
	response := helper.JSONResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
