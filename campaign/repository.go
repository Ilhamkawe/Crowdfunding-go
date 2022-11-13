package campaign

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Rewards(id int) ([]CampaignReward, error)
	Limit(num int) ([]Campaign, error)
	Find(campaignName string, campaignCattegory string) ([]Campaign, error)
	FindAll() ([]Campaign, error)
	FindAllApproved() ([]Campaign, error)
	FindByStatus(status string) ([]Campaign, error)
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
	CreateActivity(campaignActivity CampaignActivity) (CampaignActivity, error)
	UpdateActivity(campaignActivity CampaignActivity) (CampaignActivity, error)
	DeleteActivity(activityID int) (bool, error)
	FindActivity(activityID int) (CampaignActivity, error)
	FindAllActivityByCampaignID(campaignID int) ([]CampaignActivity, error)
	Paginate(limit int, offset int) ([]Campaign, error)
	isCollectAbleByDate() (bool, error)
	isCollectAbleByAmount() (bool, error)
	CreateCattegory(cattegory Cattegory) (Cattegory, error)
	DeleteCattegory(id int) (bool, error)
	FindAllCattegory() ([]Cattegory, error)
	FindPaginate(campaignName string, campaignCattegory string, limit int, offset int) ([]Campaign, error)
	FindByNamePaginate(campaignName string, limit int, offset int) ([]Campaign, error)
	FindByCattegoryPaginate(campaignCattegory string, limit int, offset int) ([]Campaign, error)
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
	err := r.db.Where("campaign_id = ?", id).Order("min_donate asc").Find(&rewards).Error
	if err != nil {
		return rewards, err
	}

	return rewards, nil
}

// mengambil data campaign dengan limit
func (r *repository) Limit(num int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Preload("User").Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status = ?", "Berjalan").Preload("User").Limit(num).Order("backer_count desc").Find(&campaigns).Error
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
func (r *repository) FindByStatus(status string) ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status = ?", status).Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
func (r *repository) FindAllApproved() ([]Campaign, error) {
	var campaigns []Campaign

	// select * from campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status = ?", "Berjalan").Preload("User").Find(&campaigns).Error
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
	err := r.db.Where("name like ?", "%"+campaignName+"%").Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign berdasarkan cattegory
func (r *repository) FindByCattegory(campaignCattegory string) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("cattegory = ?", campaignCattegory).Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status = ?", "Berjalan").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign berdasarkan nama dan cattegory
func (r *repository) Find(campaignName string, campaignCattegory string) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("name like ?", "%"+campaignName+"%").Preload("CampaignImages", "campaign_images.is_primary = 1").Where("cattegory = ?", campaignCattegory).Where("status <> ? AND status <> ?", "Pending", "Tertolak").Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign yang dibuat oleh user tertentu berdasarkan userID
func (r *repository) FindByUserID(UserID int) ([]Campaign, error) {
	var campaigns []Campaign

	// select * from campaign where user_id = ?
	err := r.db.Where("user_id = ?", UserID).Preload("User").Preload("CampaignImages", "campaign_images.is_primary = 1").Order("created_at desc").Find(&campaigns).Error

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("CampaignImages").Preload("CampaignRewards").Preload("User").Where("id = ?", ID).Where("status = ? OR status = ?", "Berjalan", "Dicairkan").Find(&campaign).Error

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

// ! Campaign Reward Repository
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

// ! Campaign Image Repository
func (r *repository) CreateImage(campaignImage CampaignImage) (CampaignImage, error) {

	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}

	return campaignImage, nil
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

// ! Activity Repository

func (r *repository) CreateActivity(campaignActivity CampaignActivity) (CampaignActivity, error) {
	err := r.db.Create(&campaignActivity).Error
	if err != nil {
		return campaignActivity, err
	}
	return campaignActivity, nil
}

func (r *repository) UpdateActivity(campaignActivity CampaignActivity) (CampaignActivity, error) {
	err := r.db.Save(&campaignActivity).Error
	if err != nil {
		return campaignActivity, err
	}
	return campaignActivity, nil
}

func (r *repository) DeleteActivity(activityID int) (bool, error) {
	err := r.db.Delete(CampaignActivity{}, activityID).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
func (r *repository) FindActivity(activityID int) (CampaignActivity, error) {
	var campaignActivity CampaignActivity
	err := r.db.Preload("Campaign").Find(&campaignActivity, activityID).Error
	if err != nil {
		return campaignActivity, err
	}
	return campaignActivity, nil
}

func (r *repository) FindAllActivityByCampaignID(campaignID int) ([]CampaignActivity, error) {
	var campaignActivities []CampaignActivity
	err := r.db.Where("campaign_id = ?", campaignID).Preload("Campaign").Order("id desc").Find(&campaignActivities).Error

	if err != nil {
		return campaignActivities, err
	}
	return campaignActivities, nil

}

func (r *repository) Paginate(limit int, offset int) ([]Campaign, error) {

	var campaign []Campaign

	err := r.db.Preload("User").Preload("CampaignImages", "campaign_images.is_primary = 1").Limit(limit).Offset(offset).Where("status = ?", "Berjalan").Find(&campaign).Error

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (r *repository) isCollectAbleByDate() (bool, error) {

	var campaign Campaign

	err := r.db.Model(&campaign).Where("collectable = ? AND finish_at = ?", false, time.Now().Local().Format("2006-01-02 15:04:05")).Updates(Campaign{Collectable: true}).Error
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repository) isCollectAbleByAmount() (bool, error) {
	err := r.db.Model(&Campaign{}).Where("collectable = ? AND goal_amount <= current_amount", false).Update("collectable", true).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

// ! Repository Cattegory Campaign

func (r *repository) CreateCattegory(cattegory Cattegory) (Cattegory, error) {
	err := r.db.Save(&cattegory).Error

	if err != nil {
		return cattegory, err
	}

	return cattegory, nil
}
func (r *repository) DeleteCattegory(id int) (bool, error) {
	err := r.db.Delete(&Cattegory{}, id).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
func (r *repository) FindAllCattegory() ([]Cattegory, error) {
	var cattegories []Cattegory
	err := r.db.Find(&cattegories).Error
	if err != nil {
		return []Cattegory{}, err
	}

	return cattegories, nil
}

// mengambil data campaign berdasarkan nama dan cattegory
func (r *repository) FindPaginate(campaignName string, campaignCattegory string, limit int, offset int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("name like ?", "%"+campaignName+"%").Preload("CampaignImages", "campaign_images.is_primary = 1").Where("cattegory = ?", campaignCattegory).Where("status = ?", "Berjalan").Limit(limit).Offset(offset).Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign berdasarkan nama
func (r *repository) FindByNamePaginate(campaignName string, limit int, offset int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("name like ?", "%"+campaignName+"%").Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status = ?", "Berjalan").Limit(limit).Offset(offset).Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// mengambil data campaign berdasarkan cattegory
func (r *repository) FindByCattegoryPaginate(campaignCattegory string, limit int, offset int) ([]Campaign, error) {
	var campaigns []Campaign
	err := r.db.Where("cattegory = ?", campaignCattegory).Preload("CampaignImages", "campaign_images.is_primary = 1").Where("status = ?", "Berjalan").Limit(limit).Offset(offset).Preload("User").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
