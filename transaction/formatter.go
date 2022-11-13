package transaction

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/user"
	"time"
)

type CampaignTransactionFormatter struct {
	ID         int       `json:"id"`
	RewardID   int       `json:"reward_id"`
	CampaignID int       `json:"campaign_id"`
	Amount     int       `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	User       user.User
	Reward     campaign.CampaignReward
	Campaign   campaign.Campaign
}

func FormatTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	formatter.User = transaction.User
	formatter.RewardID = transaction.RewardID
	formatter.CampaignID = transaction.CampaignID
	formatter.Reward = transaction.Reward
	formatter.Campaign = transaction.Campaign

	return formatter
}

func FormatTransactions(transaction []Transaction) []CampaignTransactionFormatter {
	if len(transaction) == 0 {
		return []CampaignTransactionFormatter{}
	}

	transactionsFormatter := []CampaignTransactionFormatter{}
	for _, trx := range transaction {
		transactionFormatter := FormatTransaction(trx)
		transactionsFormatter = append(transactionsFormatter, transactionFormatter)
	}

	return transactionsFormatter
}

type PaginateTransactions struct {
	Limit         int                            `json:"limit"`
	Page          int                            `json:"page"`
	PageCount     int                            `json:"page_count"`
	CountCampaign int                            `json:"count_campaign"`
	Transactions  []CampaignTransactionFormatter `json:"campaigns"`
}

type UserTransactionFormatter struct {
	ID        int               `json:"name"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageUrl = ""

	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageUrl = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = campaignFormatter

	return formatter
}

func FormatUserTransactions(transaction []Transaction) []UserTransactionFormatter {
	if len(transaction) == 0 {
		return []UserTransactionFormatter{}
	}

	var transactionFormatters []UserTransactionFormatter

	for _, trx := range transaction {
		transactionFormatter := FormatUserTransaction(trx)
		transactionFormatters = append(transactionFormatters, transactionFormatter)
	}

	return transactionFormatters
}

type TransactionFormatter struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentUrl string `json:"payment_url"`
}

func FormatPaymentTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.CampaignID = transaction.CampaignID
	formatter.Amount = transaction.Amount
	formatter.Code = transaction.Code
	formatter.PaymentUrl = transaction.PaymentUrl
	formatter.Status = transaction.Status

	return formatter
}

type CollectedCampaignFormatter struct {
	Name     string `json:"campaign_name"`
	UserName string `json:"name"`
	NoRek    string `json:"no_rekening"`
	Bank     string `json:"bank"`
	Total    int    `json:"total"`
}

func FormatPDFCollectAbleTrx(c []CollectCampaign) []CollectedCampaignFormatter {
	formatters := []CollectedCampaignFormatter{}
	formatter := CollectedCampaignFormatter{}
	for _, data := range c {

		formatter.Name = data.Campaign.Name
		formatter.UserName = data.User.Name
		formatter.NoRek = data.NoRekening
		formatter.Bank = data.Bank
		formatter.Total = data.Campaign.CurrentAmount
		formatters = append(formatters, formatter)
	}
	return formatters
}
