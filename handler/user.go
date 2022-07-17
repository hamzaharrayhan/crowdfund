package handler

import (
	"crowdfund/auth"
	"crowdfund/helper"
	"crowdfund/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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
		response := helper.JSONResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := u.authService.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.JSONResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formattedResponse := user.FormatterUserResponse(newUser, token)
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

	token, err := u.authService.GenerateToken(loggedInUser.ID)

	if err != nil {
		response := helper.JSONResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formattedResponse := user.FormatterUserResponse(loggedInUser, token)
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

func (u *userHandler) UploadUserAvatar(c *gin.Context) {
	//get input
	//panggil service handle input
	//service panggil db repository
	//db menyimpan gambar, mencari user dgn token/id, updateuser
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_Uploaded": false}
		response := helper.JSONResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := "images/" + file.Filename
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_Uploaded": false}
		response := helper.JSONResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := 1
	_, err = u.userService.SaveAvatar(userID, path)

	if err != nil {
		data := gin.H{"is_Uploaded": false}
		response := helper.JSONResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.JSONResponse("Avatar is successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
