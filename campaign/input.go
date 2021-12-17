package campaign

import "crowdfunding-TA/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type GetLimitDataInput struct {
	Limit int `uri:"limit"`
}

type GetUserCampaign struct {
	User user.User
}

type SearchCampaignInput struct {
	Name      string `json:"name"`
	Cattegory string `json:"cattegory" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string `form:"name" binding:"required"`
	ShortDescription string `form:"short_description" binding:"required"`
	Description      string `form:"description" binding:"required"`
	GoalAmount       int    `form:"goal_amount" binding:"required"`
	Cattegory        string `form:"cattegory" binding:"required"`
	Path             string `form:"path"`
	User             user.User
}

type UpdateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Cattegory        string `json:"cattegory" binding:"required"`
	User             user.User
}

type CreateCampaignImageInput struct {
	CampaignID int  `form:"campaign_id" binding:"required"`
	IsPrimary  bool `form:"is_primary" bindind:"required"`
	User       user.User
}

type CreateCampaignRewardInput struct {
	CampaignID  int    `json:"campaign_id" binding:"required"`
	Description string `json:"description" binding:"required"`
	Perks       string `json:"perks" binding:"required"`
	MinDonate   int    `json:"min_donate" binding:"required"`
	User        user.User
}

type DeleteCampaignRewardInput struct {
	CampaignID int `json:"campaign_id" binding:"required"`
	RewardID   int `json:"reward_id" binding:"required"`
	User       user.User
}
type DeleteCampaignImageInput struct {
	CampaignID int `json:"campaign_id" binding:"required"`
	ImageID    int `json:"image_id" binding:"required"`
	User       user.User
}
type UpdateAttachmentInput struct {
	CampaignID int    `form:"campaign_id" binding:"required"`
	Action     string `form:"action" binding:"required"`
	Path       string `form:"path"`
	User       user.User
}
type FormCampaignInput struct {
	Name             string `form:"name" binding:"required"`
	ShortDescription string `form:"short_description" binding:"required"`
	Description      string `form:"description" binding:"required"`
	GoalAmount       int    `form:"goal_amount" binding:"required"`
	Perks            string `form:"perks" binding:"required"`
	UserID           int    `form:"user_id" binding:"required"`
	Users            []user.User
}

type FormCampaignUpdate struct {
	ID               int
	Name             string `form:"name" binding:"required"`
	ShortDescription string `form:"short_description" binding:"required"`
	Description      string `form:"description" binding:"required"`
	GoalAmount       int    `form:"goal_amount" binding:"required"`
	Perks            string `form:"perks" binding:"required"`
	Users            []user.User
}
