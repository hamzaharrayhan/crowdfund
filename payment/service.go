package payment

import (
	"crowdfund/user"
	"strconv"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Service interface {
	GetPaymentUrl(transaction Transaction, user user.User) (string, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (service *service) GetPaymentUrl(transaction Transaction, user user.User) (string, error) {
	midtrans.ServerKey = "SB-Mid-server-RlnRwkZpKPnCXXfYpj53FrqD"
	midtrans.Environment = midtrans.Sandbox

	//Initiate client for Midtrans Snap
	var snapGateway = snap.Client{}
	snapGateway.New(midtrans.ServerKey, midtrans.Sandbox)

	// Initiate Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	snapTokenResp, err := snap.CreateTransaction(snapReq)
	if err != nil {
		return snapTokenResp.RedirectURL, err
	}
	return snapTokenResp.RedirectURL, nil
}
