package transaction

import (
	"crowdfunding-TA/campaign"
	"crowdfunding-TA/payment"
	"fmt"
	"math"
	"strconv"

	// "crowdfunding-TA/payment"
	"errors"
)

type Service interface {
	GetTransactionsByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
	GetTransactionsByUserID(UserID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
	FindAll() ([]Transaction, error)
	FindAllByReward(Rid int, Cid int, Uid int, input campaign.PaginateCampaignInput) (PaginateTransactions, error)
	CollectAmount(input CollectInput) (CollectCampaign, error)
	FindCollectData(id int) ([]CollectCampaign, error)
	FindAllCollectData() ([]CollectCampaign, error)
	FindCollectDataByCID(id int) (CollectCampaign, error)
	ChangeCollectStatus(Status string, ID int) (CollectCampaign, error)
}

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionsByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindByID(input.ID)

	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("anda bukan pemilik campaign ini")
	}

	transaction, err := s.repository.GetByCampaignID(input.ID)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *service) GetTransactionsByUserID(UserID int) ([]Transaction, error) {

	transaction, err := s.repository.GetByUserID(UserID)

	if err != nil {
		return []Transaction{}, err
	}

	return transaction, nil

}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Email = input.Email
	transaction.Status = "Pending"
	transaction.RewardID = input.RewardID

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
	}

	paymentUrl, err := s.paymentService.GetPaymentUrl(paymentTransaction, input.User)
	fmt.Println(err)
	if err != nil {
		return newTransaction, nil
	}

	newTransaction.PaymentUrl = paymentUrl

	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {

	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(transaction_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlemen" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)

	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount
		if !campaign.Collectable && campaign.CurrentAmount >= campaign.GoalAmount {
			campaign.Collectable = true
		}
		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) FindAll() ([]Transaction, error) {

	transaction, err := s.repository.FindAll()
	if err != nil {
		return transaction, err
	}

	return transaction, nil

}

func (s *service) FindAllByReward(Rid int, Cid int, Uid int, input campaign.PaginateCampaignInput) (PaginateTransactions, error) {
	transaction, err := s.repository.FindAllByReward(Rid, Cid)

	if err != nil {
		return PaginateTransactions{}, err
	}

	offset := (input.ActivePage * input.Limit) - input.Limit

	countCampaigns := len(transaction)
	pages := int(math.Ceil(float64(countCampaigns) / float64(input.Limit)))

	if input.ActivePage > pages {
		return PaginateTransactions{}, errors.New("data tidak ada")
	}

	paginateTransactions, err := s.repository.FindAllByRewardPaginate(Rid, Cid, input.Limit, offset)

	if err != nil {
		return PaginateTransactions{}, err
	}

	var PaginateTransaction PaginateTransactions
	PaginateTransaction.Limit = input.Limit
	PaginateTransaction.CountCampaign = countCampaigns
	PaginateTransaction.PageCount = pages
	PaginateTransaction.Page = input.ActivePage
	PaginateTransaction.Transactions = FormatTransactions(paginateTransactions)

	return PaginateTransaction, nil
}

func (s *service) CollectAmount(input CollectInput) (CollectCampaign, error) {
	collect := CollectCampaign{}
	collect.CampaignID = input.CampaignID
	collect.UserID = input.UserID
	collect.AccountName = input.AccountName
	collect.NoRekening = input.NoRekening
	collect.Bank = input.Bank
	collect.Status = "Pending"

	newCollect, err := s.repository.CollectAmount(collect)
	if err != nil {
		return newCollect, err
	}

	campaign, err := s.campaignRepository.FindUserCampaign(input.CampaignID)

	if err != nil {
		return CollectCampaign{}, err
	}

	campaign.Status = "Dicairkan"

	fmt.Println(campaign)

	_, err = s.campaignRepository.Update(campaign)
	if err != nil {
		return CollectCampaign{}, err
	}

	return newCollect, nil

}

func (s *service) FindCollectData(id int) ([]CollectCampaign, error) {

	collectData, err := s.repository.FindCollectDataByID(id)

	if err != nil {
		return []CollectCampaign{}, err
	}

	return collectData, nil

}
func (s *service) FindCollectDataByCID(id int) (CollectCampaign, error) {

	collectData, err := s.repository.FindCollectDataByCID(id)
	if err != nil {
		return CollectCampaign{}, err
	}

	return collectData, nil
}
func (s *service) FindAllCollectData() ([]CollectCampaign, error) {

	collectData, err := s.repository.FindAllCollectData()

	if err != nil {
		return []CollectCampaign{}, err
	}

	return collectData, nil

}

func (s *service) ChangeCollectStatus(Status string, ID int) (CollectCampaign, error) {
	collectData, err := s.repository.FindCollectDataByCID(ID)
	if err != nil {
		return CollectCampaign{}, err
	}

	collectData.Status = Status

	updateCollect, err := s.repository.UpdateCollect(collectData)
	if err != nil {
		return updateCollect, err
	}

	return updateCollect, nil
}
