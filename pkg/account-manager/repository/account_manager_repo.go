package repository

import (
	models "linkaja.com/e-wallet/lib/base_models"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
)

// AccountManagerRepo : Account Manager Repository Interface
type AccountManagerRepo interface {
	GetAccountInfo(accNmbr int) models.Result
	UpdateBalance(account model.CustomerAccount) bool
	CheckUserAndPassword(cust *model.Customer) models.Result
}
