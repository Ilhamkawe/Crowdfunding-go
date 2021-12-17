package campaign

import "strings"

type CampaignFormatter struct {
	ID               int                   `json:"id"`
	UserID           int                   `json:"user_id"`
	Name             string                `json:"name"`
	ShortDescription string                `json:"short_description"`
	ImageURL         string                `json:"image_url"`
	GoalAmount       int                   `json:"goal_amount"`
	CurrentAmount    int                   `json:"current_amount"`
	Cattegory        string                `json:"cattegory"`
	Author           string                `json:"author"`
	Slug             string                `json:"slug"`
	Status           string                `json:"status"`
	User             CampaignUserFormatter `json:"user"`
}

// memformat data campaign
func FormatCampaign(campaign Campaign) CampaignFormatter {
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.Cattegory = campaign.Cattegory
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.ImageURL = ""
	campaignFormatter.Author = campaign.User.Name
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.Status = campaign.Status
	user := campaign.User

	CampaignUserFormatter := CampaignUserFormatter{}
	CampaignUserFormatter.Name = user.Name
	CampaignUserFormatter.ImageUrl = user.AvatarFileName

	campaignFormatter.User = CampaignUserFormatter

	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	return campaignFormatter
}

// memformat data list campaign (hasil dari FormatCampaign)
func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {

	campaignsFormatter := []CampaignFormatter{}
	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}
	return campaignsFormatter
}

type CampaignDetailFormatter struct {
	ID               int                       `json:"id" `
	Name             string                    `json:"name"`
	ShortDescription string                    `json:"short_description"`
	Description      string                    `json:"description"`
	ImageUrl         string                    `json:"image_url"`
	GoalAmount       int                       `json:"goal_amount"`
	CurrentAmount    int                       `json:"current_amount"`
	BackerCount      int                       `json:"backer_count"`
	UserID           int                       `json:"user_id"`
	Cattegory        string                    `json:"cattegory"`
	Slug             string                    `json:"slug"`
	Status           string                    `json:"status"`
	Attachment       string                    `json:"attachment"`
	Rewards          []CampaignRewardFormatter `json:"reward"`
	User             CampaignUserFormatter     `json:"user"`
	Images           []CampaignImageFormatter  `json:"images"`
}

type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type CampaignImageFormatter struct {
	ID        int    `json:"id"`
	ImageUrl  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type CampaignRewardFormatter struct {
	ID          int      `json:"id"`
	Description string   `json:"description"`
	Perks       []string `json:"perks"`
	MinDonate   int      `json:"min_donate"`
}

// memformat data campaign reward
func FormatCampaignReward(rewards []CampaignReward) []CampaignRewardFormatter {
	campaignRewardsFormatter := []CampaignRewardFormatter{}
	for _, reward := range rewards {
		campaignRewardFormatter := CampaignRewardFormatter{reward.ID, reward.Description, strings.Split(reward.Perks, ","), reward.MinDonate}
		campaignRewardsFormatter = append(campaignRewardsFormatter, campaignRewardFormatter)
	}

	return campaignRewardsFormatter
}

// memformat data campaign detail
func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.ImageUrl = ""
	campaignDetailFormatter.BackerCount = campaign.BackerCount
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.UserID = campaign.UserID
	campaignDetailFormatter.Cattegory = campaign.Cattegory
	campaignDetailFormatter.Slug = campaign.Slug
	campaignDetailFormatter.Status = campaign.Status
	campaignDetailFormatter.Attachment = campaign.Attachment
	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageUrl = campaign.CampaignImages[0].FileName
	}

	user := campaign.User
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageUrl = user.AvatarFileName

	campaignDetailFormatter.User = campaignUserFormatter

	// set Cmpaign images
	images := []CampaignImageFormatter{}
	for _, image := range campaign.CampaignImages {
		campaignImageFormatter := CampaignImageFormatter{}
		campaignImageFormatter.ImageUrl = image.FileName
		campaignImageFormatter.ID = image.ID

		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
			campaignDetailFormatter.ImageUrl = image.FileName
		}
		campaignImageFormatter.IsPrimary = isPrimary

		images = append(images, campaignImageFormatter)

	}

	campaignDetailFormatter.Images = images

	// set CampaingReward
	rewards := []CampaignRewardFormatter{}
	for _, reward := range campaign.CampaignRewards {
		campaignRewardFormatter := CampaignRewardFormatter{}
		campaignRewardFormatter.ID = reward.ID
		campaignRewardFormatter.Description = reward.Description

		// set reward perks
		var perks []string = strings.Split(reward.Perks, ",")
		campaignRewardFormatter.Perks = perks

		campaignRewardFormatter.MinDonate = reward.MinDonate

		rewards = append(rewards, campaignRewardFormatter)
	}

	campaignDetailFormatter.Rewards = rewards

	return campaignDetailFormatter
}
