package transaction

import "crowdfunding-TA/user"

type GetCampaignTransactionInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}

type CreateTransactionInput struct {
	Amount     int    `json:"amount" binding:"required"`
	CampaignID int    `json:"campaign_id" binding:"required"`
	Email      string `json:"email"`
	RewardID   int    `json:"reward_id"`
	User       user.User
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

type CollectInput struct {
	CampaignID  int    `json:"campaign_id" binding:"required"`
	UserID      int    `json:"user_id" binding:"required"`
	AccountName string `json:"account_name" binding:"required"`
	NoRekening  string `json:"no_rekening" binding:"required"`
	Bank        string `json:"bank" binding:"required"`
	Status      string `json:"status"`
}
