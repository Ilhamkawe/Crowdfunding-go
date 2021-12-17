package payment

import (
	"crowdfunding-TA/campaign"
	// "crowdfunding-TA/transaction"
	"crowdfunding-TA/user"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type service struct {
	// TransactionRepository transaction.Repository
	CampaignRepository campaign.Repository
}

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
	// ProcessPayment(input transaction.TransactionNotificationInput) error
}

func NewService(CampaignRepository campaign.Repository) *service {
	return &service{CampaignRepository}
}

func (s *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-r_sQ4bvvoYrkeot4PdJNa1XL"
	midclient.ClientKey = "SB-Mid-client-DIX41Cb-1atPEusP"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		}, CustomerDetail: &midtrans.CustDetail{
			FName: user.Name,
			Email: user.Email,
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)

	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}
