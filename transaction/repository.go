package transaction

import (
	"gorm.io/gorm"
)

type Respository interface {
	FindByCampaignID(id int) ([]Transaction, error)
	FindByUserID(id int) ([]Transaction, error)
	FindByID(id int) (Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
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

func (r *repository) FindByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) FindByID(id int) (Transaction, error) {
	transaction := Transaction{}
	err := r.db.Where("id = ?", id).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
