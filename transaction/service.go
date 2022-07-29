package transaction

import (
	"crowdfund/campaign"
	"errors"
)

type Service interface {
	GetByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
}

type service struct {
	repository         Respository
	campaignRepository campaign.Repository
}

func NewService(repository Respository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetByCampaignID(input GetCampaignTransactionsInput) ([]Transaction, error) {
	newCampaign, err := s.campaignRepository.FindByID(input.ID)

	if err != nil {
		return []Transaction{}, err
	}

	if newCampaign.User.ID != input.User.ID {
		return []Transaction{}, errors.New("not an owner of the campaign")
	}

	transactions, err := s.repository.FindByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.FindByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}
