package transaction

import (
	"gorm.io/gorm"
)

type Respository interface {
	FindByCampaignID(id int) ([]Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByCampaignID(input int) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Preload("User").Where("campaign_id = ?", input).Order("id desc").Find(&transactions).Error
	if err != nil {
		return []Transaction{}, err
	}
	return transactions, nil
}
