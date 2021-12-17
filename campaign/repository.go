package campaign

import "gorm.io/gorm"

type Repository interface {
	Rewards(id int) ([]CampaignReward, error)
	Limit(num int) ([]Campaign, error)
	Find(campaignName string, campaignCattegory string) ([]Campaign, error)
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	FindByID(ID int) (Campaign, error)
	FindUserCampaign(ID int) (Campaign, error)
	FindByName(campaignName string) ([]Campaign, error)
	FindByCattegory(campaignCattegory string) ([]Campaign, error)
	Save(campaign Campaign) (Campaign, error)
	Update(campaign Campaign) (Campaign, error)
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
	CreateImage(campaignImage CampaignImage) (CampaignImage, error)
	CreateReward(campaignReward CampaignReward) (CampaignReward, error)
	DeleteReward(rewardID int) (bool, error)
	DeleteImage(imageID int) (bool, error)
	UpdateAttachment(campaign Campaign) (Campaign, error)
	FindAllWoStatus() ([]Campaign, error)
}

type repository struct {
	db *gorm.DB
}

// mendeklarasikan repository untuk campaign
func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// mengambil data rewards berdasarkan campaign_id
func (r *repository) Rewards(id int) ([]CampaignReward, error) {
	var rewards []CampaignReward
	err := r.db.Where("campaign_id = ?", id).Find(&rewards).Error
	if err != nil {
		return rewards, err
	}

	return rewards, nil
}

// mengambil data campaign dengan limit
func (r *repository) Limit(num int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("User").Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Limit(num).Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// mengambil semua data pada tabel campaign
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	// select * from campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
func (r *repository) FindAllWoStatus() ([]Campaign, error) {
	var campaigns []Campaign

	// select * from campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// mengambil data campaign berdasarkan nama
func (r *repository) FindByName(campaignName string) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("name like ?", "%"+campaignName+"%").Preload("CampaignImages").Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign berdasarkan cattegory
func (r *repository) FindByCattegory(campaignCattegory string) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("cattegory = ?", campaignCattegory).Preload("CampaignImages").Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign berdasarkan nama dan cattegory
func (r *repository) Find(campaignName string, campaignCattegory string) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("name like ?", "%"+campaignName+"%").Preload("CampaignImages").Where("cattegory = ?", campaignCattegory).Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Find(&campaigns).Error
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

func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("CampaignImages").Preload("CampaignRewards").Preload("User").Where("id = ?", ID).Where("status <> ? AND status <> ?", "Pending", "Tertolak").Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) FindUserCampaign(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("CampaignImages").Preload("CampaignRewards").Preload("User").Where("id = ?", ID).Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}

	return campaign, err
}

func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) UpdateAttachment(campaign Campaign) (Campaign, error) {
	err := r.db.Model(&campaign).Updates(Campaign{Attachment: campaign.Attachment}).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *repository) CreateImage(campaignImage CampaignImage) (CampaignImage, error) {

	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
}

func (r *repository) CreateReward(campaignReward CampaignReward) (CampaignReward, error) {
	err := r.db.Create(&campaignReward).Error
	if err != nil {
		return campaignReward, err
	}
	return campaignReward, nil
}

func (r *repository) DeleteReward(rewardID int) (bool, error) {
	err := r.db.Delete(&CampaignReward{}, rewardID).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
func (r *repository) DeleteImage(imageID int) (bool, error) {
	err := r.db.Delete(&CampaignImage{}, imageID).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *repository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
