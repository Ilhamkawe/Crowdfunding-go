package campaign

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	GetUserCampaignByID(inputID GetCampaignDetailInput, inputUser GetUserCampaign) (Campaign, error)
	GetAllCampaign() ([]Campaign, error)
	GetCampaignByIDWoStatus(id int) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData UpdateCampaignInput) (Campaign, error)
	UpdateAttachment(input UpdateAttachmentInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
	SaveCampaignReward(input CreateCampaignRewardInput) (CampaignReward, error)
	Limit(num int) ([]Campaign, error)
	GetRewards(input GetCampaignDetailInput) ([]CampaignReward, error)
	SearchCampaign(input SearchCampaignInput) ([]Campaign, error)
	DeleteReward(input DeleteCampaignRewardInput) (bool, error)
	DeleteImage(input DeleteCampaignImageInput) (bool, error)
	ChangeStatus(Status string, ID int) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetRewards(input GetCampaignDetailInput) ([]CampaignReward, error) {
	var rewards []CampaignReward

	rewards, err := s.repository.Rewards(input.ID)
	if err != nil {
		return rewards, err
	}

	return rewards, err
}

func (s *service) Limit(num int) ([]Campaign, error) {
	var campaigns []Campaign

	campaigns, err := s.repository.Limit(num)
	if err != nil {
		return campaigns, err
	}

	return campaigns, err
}

func (s *service) SearchCampaign(input SearchCampaignInput) ([]Campaign, error) {
	var campaigns []Campaign
	var err error

	if input.Name != "" {
		if input.Cattegory != "semua" {
			// mencari berdasarkan nama dan cattegory
			campaigns, err = s.repository.Find(input.Name, input.Cattegory)
		} else {
			// mencari berdasarkan nama
			campaigns, err = s.repository.FindByName(input.Name)
		}
	} else if input.Cattegory != "" {
		// mencari berdadsarkan cattegory
		campaigns, err = s.repository.FindByCattegory(input.Cattegory)
	} else {
		// tidak mencari apa apa
		campaigns = []Campaign{}
	}

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

//  func yang digunakan untuk mengambil campaign, jika parameter > 0 maka akan menampilkan campaign yang dibuat user tertentu, jika diisi 0 maka akan menampilkan semua campaign
func (s *service) GetCampaigns(UserID int) ([]Campaign, error) {
	var campaigns []Campaign

	if UserID != 0 {
		campaigns, err := s.repository.FindByUserID(UserID)
		if err != nil {
			return campaigns, err
		}

		return campaigns, err
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetAllCampaign() ([]Campaign, error) {
	var campaigns []Campaign

	campaigns, err := s.repository.FindAllWoStatus()
	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

// mengambil data campaign tertentu yang dimiliki user
func (s *service) GetCampaignByIDWoStatus(id int) (Campaign, error) {
	var campaign Campaign

	campaign, err := s.repository.FindUserCampaign(id)

	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

// mengambil data campaign tertentu yang dimiliki user dan berstatus berjalan
func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	var campaign Campaign

	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

// mengambil semua data campaign yang dimiliki user
func (s *service) GetUserCampaignByID(inputID GetCampaignDetailInput, inputUser GetUserCampaign) (Campaign, error) {
	var campaign Campaign
	campaign, err := s.repository.FindUserCampaign(inputID.ID)

	if err != nil {
		return campaign, err
	}

	if inputUser.User.ID != campaign.UserID {
		return campaign, errors.New("anda bukan owner campaign ini")
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID
	campaign.Cattegory = input.Cattegory
	campaign.Status = "Pending"
	campaign.Attachment = input.Path
	// slug
	slugCanditate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCanditate)

	newCampaign, err := s.repository.Save(campaign)

	if err != nil {
		return newCampaign, err
	}

	return newCampaign, err
}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData UpdateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindUserCampaign(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if inputData.User.ID != campaign.UserID {
		return campaign, errors.New("anda bukan owner campaign ini")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.GoalAmount = inputData.GoalAmount

	updateCampaign, err := s.repository.Update(campaign)

	if err != nil {
		return updateCampaign, err
	}

	return updateCampaign, nil
}

func (s *service) ChangeStatus(Status string, ID int) (Campaign, error) {

	campaign, err := s.repository.FindUserCampaign(ID)

	if err != nil {
		return campaign, err
	}

	campaign.Status = Status

	fmt.Println(campaign)

	updateCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updateCampaign, err
	}

	return updateCampaign, nil

}

func (s *service) UpdateAttachment(input UpdateAttachmentInput) (Campaign, error) {
	campaign, err := s.repository.FindUserCampaign(input.CampaignID)
	if err != nil {
		return campaign, err
	}

	if input.User.ID != campaign.UserID {
		return campaign, errors.New("anda bukan owner campaign ini")
	}

	if strings.ToUpper(input.Action) == "DELETE" {
		// mengubah attachment menjadi string kosong
		campaign.Attachment = " "
	} else if strings.ToUpper(input.Action) == "UPLOAD" {
		// upload file baru dan mengubah record attachment dengan nama file yg baru
		campaign.Attachment = input.Path
	} else {
		return campaign, errors.New("aksi tidak ditemukan")
	}

	fmt.Println(input)
	updateAttachment, err := s.repository.UpdateAttachment(campaign)
	if err != nil {
		return updateAttachment, err
	}

	return updateAttachment, nil
}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if input.User.ID != campaign.UserID {
		return CampaignImage{}, errors.New("anda bukan owner campaign ini")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.FileName = fileLocation
	campaignImage.IsPrimary = isPrimary

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil

}

func (s *service) SaveCampaignReward(input CreateCampaignRewardInput) (CampaignReward, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignReward{}, err
	}

	if input.User.ID != campaign.UserID {
		return CampaignReward{}, errors.New("anda bukan owner campaign ini")
	}

	campaignReward := CampaignReward{}
	campaignReward.CampaignID = input.CampaignID
	campaignReward.Description = input.Description
	campaignReward.MinDonate = input.MinDonate
	campaignReward.Perks = input.Perks

	newCampaignReward, err := s.repository.CreateReward(campaignReward)

	if err != nil {
		return newCampaignReward, err
	}

	return newCampaignReward, nil
}

func (s *service) DeleteReward(input DeleteCampaignRewardInput) (bool, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return false, err
	}

	if input.User.ID != campaign.UserID {
		return false, errors.New("anda bukan owner campaign ini")
	}

	_, err = s.repository.DeleteReward(input.RewardID)
	if err != nil {
		return false, err
	}

	return true, nil
}
func (s *service) DeleteImage(input DeleteCampaignImageInput) (bool, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return false, err
	}

	if input.User.ID != campaign.UserID {
		return false, errors.New("anda bukan owner campaign ini")
	}

	_, err = s.repository.DeleteImage(input.ImageID)
	if err != nil {
		return false, err
	}

	return true, nil
}
