package handler

import (
	"crowdfund/campaign"
	"crowdfund/helper"
	"crowdfund/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))
	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.JSONResponse("Failed to retrieve campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.JSONResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusBadRequest, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.JSONResponse("Failed to retrieve campaign with corresponding ID", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignItem, err := h.service.GetCampaignByID(inputID.ID)
	if err != nil {
		response := helper.JSONResponse("Failed to retrieve campaign with corresponding ID", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.JSONResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignItem))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var inputNewCampaign campaign.CreateCampaignInput
	err := c.ShouldBind(&inputNewCampaign)
	if err != nil {
		errors := helper.ResponseValidationError(err)
		errorMessages := gin.H{"errors": errors}
		response := helper.JSONResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessages)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	inputNewCampaign.User = c.MustGet("currentUser").(user.User)
	newCampaign, err := h.service.CreateCampaign(inputNewCampaign)
	if err != nil {
		response := helper.JSONResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.JSONResponse("New campaign successfully added", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var idCampaign campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&idCampaign)
	if err != nil {
		response := helper.JSONResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
	}

	var inputUpdateCampaign campaign.CreateCampaignInput
	err = c.ShouldBind(&inputUpdateCampaign)
	if err != nil {
		errors := helper.ResponseValidationError(err)
		errorMessages := gin.H{"errors": errors}
		response := helper.JSONResponse("Failed to update campaign", http.StatusBadRequest, "error", errorMessages)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	inputUpdateCampaign.User = c.MustGet("currentUser").(user.User)

	updatedCampaign, err := h.service.UpdateCampaign(idCampaign, inputUpdateCampaign)
	if err != nil {
		response := helper.JSONResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.JSONResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) SaveCampaignImage(c *gin.Context) {
	var campaignImageInput campaign.CreateCampaignImageInput

	err := c.ShouldBind(&campaignImageInput)

	if err != nil {
		if err != nil {
			errors := helper.ResponseValidationError(err)
			errorMessages := gin.H{"errors": errors}
			response := helper.JSONResponse("Failed to upload campaign image", http.StatusBadRequest, "error", errorMessages)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_Uploaded": false}
		response := helper.JSONResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	campaignImageInput.User = currentUser

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_Uploaded": false}
		response := helper.JSONResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(campaignImageInput, path)
	if err != nil {
		data := gin.H{"is_Uploaded": false}
		response := helper.JSONResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.JSONResponse("Succes to upload campaign image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
