package campaign

import (
	"crowdfunding-TA/user"
	"time"

	"github.com/leekchan/accounting"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	BackerCount      int
	GoalAmount       int
	Cattegory        string
	CurrentAmount    int
	Slug             string
	Status           string
	Attachment       string
	Collectable      bool
	FinishAt         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	User             user.User
	CampaignImages   []CampaignImage
	CampaignRewards  []CampaignReward
}

func (c Campaign) GoalAmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Decimal: ",", Thousand: "."}
	return ac.FormatMoney(c.GoalAmount)
}
func (c Campaign) CurrentAmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp", Precision: 2, Decimal: ",", Thousand: "."}
	return ac.FormatMoney(c.CurrentAmount)
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CampaignReward struct {
	ID          int
	CampaignID  int
	Description string
	Perks       string
	MinDonate   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CampaignActivity struct {
	ID               int
	CampaignID       int
	Name             string
	ShortDescription string
	Description      string
	ImageUrl         string
	Slug             string
	Campaign         Campaign
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Cattegory struct {
	ID   int
	Name string
}
