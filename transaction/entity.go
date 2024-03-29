package transaction

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/user"
	"time"

	"github.com/leekchan/accounting"
)

type CollectCampaign struct {
	ID          int
	CampaignID  int
	UserID      int
	AccountName string
	NoRekening  string
	Bank        string
	Status      string
	User        user.User
	Campaign    campaign.Campaign
}

type Transaction struct {
	ID         int
	CampaignID int
	RewardID   int
	Email      string
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentUrl string
	User       user.User
	Campaign   campaign.Campaign
	Reward     campaign.CampaignReward
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (t Transaction) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Decimal: ",", Thousand: "."}
	return ac.FormatMoney(t.Amount)
}
