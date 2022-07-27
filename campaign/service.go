package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(inputID int) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(campaignID GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(campaignImageInput CreateCampaignImageInput, imageLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignByID(inputID int) (Campaign, error) {
	campaign, err := s.repository.FindByID(inputID)
	if err != nil {
		return campaign, err
	}

	return campaign, err
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	inputtedCampaign := Campaign{}
	inputtedCampaign.Name = input.Name
	inputtedCampaign.ShortDescription = input.ShortDescription
	inputtedCampaign.Description = input.Description
	inputtedCampaign.GoalAmount = input.GoalAmount
	inputtedCampaign.Perks = input.Perks
	inputtedCampaign.UserID = input.User.ID

	// bikin slug
	slugName := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	inputtedCampaign.Slug = slug.Make(slugName)

	//simpan campaign as new campaign
	newCampaign, err := s.repository.Save(inputtedCampaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

func (s *service) UpdateCampaign(campaignID GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error) {

	campaign, err := s.repository.FindByID(campaignID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != input.User.ID {
		return campaign, errors.New("Not the owner of the campaign")
	}

	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}

	return updatedCampaign, nil
}

func (s *service) SaveCampaignImage(campaignImageInput CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaign, err := s.repository.FindByID(campaignImageInput.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != campaignImageInput.User.ID {
		return CampaignImage{}, errors.New("not the owner of the campaign")
	}

	campaignImage := CampaignImage{}
	isPrimary := 0

	if campaignImageInput.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImageAsNonPrimary(campaignImageInput.CampaignID)

		if err != nil {
			return CampaignImage{}, err
		}
	}
	campaignImage.CampaignID = campaignImageInput.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation
	savedImage, err := s.repository.CreateImage(campaignImage)

	if err != nil {
		return savedImage, err
	}

	return savedImage, nil
}
