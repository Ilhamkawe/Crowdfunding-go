package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
}

type repository struct {
	db *gorm.DB
}

// mendeklarasikan repository untuk campaign
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// mengambil semua data pada tabel campaign
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	// select * from campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// mengambil data campaign yang dibuat oleh user tertentu berdasarkan userID
func (r *repository) FindByUserID(UserID int) ([]Campaign, error) {
	var campaigns []Campaign

	// select * from campaign where user_id = ?
	err := r.db.Where("user_id = ?", UserID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
