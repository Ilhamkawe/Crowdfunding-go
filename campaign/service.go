package campaign

type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
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

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	var campaign Campaign

	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, err

}
