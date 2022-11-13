package transaction

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetByCampaignID(ID int) ([]Transaction, error)
	GetByUserID(userID int) ([]Transaction, error)
	GetByID(ID int) (Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
	FindAll() ([]Transaction, error)
	FindAllByReward(Rid int, Cid int) ([]Transaction, error)
	FindAllByRewardPaginate(Rid int, Cid int, limit int, offset int) ([]Transaction, error)
	CollectAmount(input CollectCampaign) (CollectCampaign, error)
	FindCollectDataByID(id int) ([]CollectCampaign, error)
	FindAllCollectData() ([]CollectCampaign, error)
	FindAllPendingCollectData() ([]CollectCampaign, error)
	FindCollectDataByCID(id int) (CollectCampaign, error)
	UpdateCollect(collect CollectCampaign) (CollectCampaign, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(ID int) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Preload("User").Preload("Campaign").Where("campaign_id = ?", ID).Where("status = ?", "paid").Order("id desc").Find(&transactions).Error

	if err != nil {
		return transactions, err
	}

	return transactions, err
}

func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("created_at desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByID(ID int) (Transaction, error) {
	var transaction Transaction

	err := r.db.Where("id = ?", ID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) FindAll() ([]Transaction, error) {

	var transaction []Transaction

	err := r.db.Preload("Campaign").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, err

}

func (r *repository) FindAllByReward(Rid int, Cid int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("Campaign").Preload("Reward").Preload("User").Where("reward_id = ? AND campaign_id = ? AND status = ?", Rid, Cid, "paid").Order("created_at desc").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, err

}

func (r *repository) FindAllByRewardPaginate(Rid int, Cid int, limit int, offset int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("Campaign").Preload("Reward").Preload("User").Where("reward_id = ? AND campaign_id = ? AND status = ?", Rid, Cid, "paid").Limit(limit).Offset(offset).Order("created_at desc").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, err

}

func (r *repository) CollectAmount(input CollectCampaign) (CollectCampaign, error) {
	err := r.db.Create(&input).Error

	if err != nil {
		return input, err
	}

	return input, nil
}

func (r *repository) FindCollectDataByID(id int) ([]CollectCampaign, error) {
	var collectData []CollectCampaign
	err := r.db.Where("campaign_id = ?", id).Find(&collectData).Error

	if err != nil {
		return []CollectCampaign{}, err
	}

	return collectData, nil
}

func (r *repository) FindCollectDataByCID(id int) (CollectCampaign, error) {
	var collectData CollectCampaign
	err := r.db.Preload("Campaign").Preload("User").Where("id = ?", id).Find(&collectData).Error

	if err != nil {
		return CollectCampaign{}, err
	}
	fmt.Println(err)
	fmt.Println(collectData)
	return collectData, nil
}

func (r *repository) FindAllCollectData() ([]CollectCampaign, error) {
	var collectData []CollectCampaign
	err := r.db.Preload("Campaign").Preload("User").Order("Status asc").Find(&collectData).Error

	if err != nil {
		return []CollectCampaign{}, err
	}

	return collectData, nil
}

func (r *repository) FindAllPendingCollectData() ([]CollectCampaign, error) {
	var collectData []CollectCampaign
	err := r.db.Preload("Campaign").Preload("User").Where("status = ?", "Pending").Find(&collectData).Error

	if err != nil {
		return []CollectCampaign{}, err
	}

	return collectData, nil
}

func (r *repository) UpdateCollect(collect CollectCampaign) (CollectCampaign, error) {
	err := r.db.Save(&collect).Error
	if err != nil {
		return collect, err
	}
	return collect, nil
}
