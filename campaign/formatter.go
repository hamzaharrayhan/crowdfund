package campaign

import (
	"strings"
)

type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func FormatCampaign(campaign Campaign) CampaignFormatter {
	formattedCampaign := CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		ImageURL:         "",
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
	}

	if len(campaign.CampaignImages) > 0 {
		formattedCampaign.ImageURL = campaign.CampaignImages[0].FileName
	}

	return formattedCampaign
}

func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	formattedCampaigns := []CampaignFormatter{}
	for _, campaign := range campaigns {
		formattedCampaigns = append(formattedCampaigns, FormatCampaign(campaign))
	}
	return formattedCampaigns
}

type CampaignDetailFormatter struct {
	ID               int                             `json:"id"`
	Name             string                          `json:"name"`
	ShortDescription string                          `json:"short_description"`
	ImageURL         string                          `json:"image_url"`
	GoalAmount       int                             `json:"goal_amount"`
	CurrentAmount    int                             `json:"current_amount"`
	UserID           int                             `json:"user_id"`
	Slug             string                          `json:"slug"`
	Description      string                          `json:"description"`
	User             UserCampaignDetailFormatter     `json:"user"`
	Perks            []string                        `json:"perks"`
	Images           []ImagesCampaignDetailFormatter `json:"images"`
}

type UserCampaignDetailFormatter struct {
	Name     string
	ImageURL string `json:"image_url"`
}

type ImagesCampaignDetailFormatter struct {
	ImageURL  string
	IsPrimary bool
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	formattedCampaign := CampaignDetailFormatter{
		ID:               campaign.ID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		UserID:           campaign.UserID,
		Slug:             campaign.Slug,
		Description:      campaign.Description,
	}

	if len(campaign.CampaignImages) > 0 {
		formattedCampaign.ImageURL = campaign.CampaignImages[0].FileName
	}

	perks := strings.Split(campaign.Perks, ",")
	for _, perk := range perks {
		formattedCampaign.Perks = append(formattedCampaign.Perks, strings.TrimSpace(perk))
	}

	user := campaign.User
	userFormat := UserCampaignDetailFormatter{
		Name:     user.Name,
		ImageURL: user.AvatarFileName,
	}

	formattedCampaign.User = userFormat

	images := ImagesCampaignDetailFormatter{}
	imagesList := []ImagesCampaignDetailFormatter{}

	for _, image := range campaign.CampaignImages {
		images.ImageURL = image.FileName
		images.IsPrimary = image.IsPrimary
		imagesList = append(imagesList, images)
	}

	formattedCampaign.Images = imagesList

	return formattedCampaign
}
