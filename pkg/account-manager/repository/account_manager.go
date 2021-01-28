package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	models "linkaja.com/e-wallet/lib/base_models"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
)

type accountManagerRepo struct {
	dbUser *gorm.DB
}

// NewAccountManagerRepo : Account Manager Repo Builder
func NewAccountManagerRepo(db *gorm.DB) AccountManagerRepo {
	return &accountManagerRepo{
		dbUser: db,
	}
}

func (r *accountManagerRepo) GetAccountInfo(accNmbr int) models.Result {
	var accInfo model.AccountInfo
	if err := r.dbUser.Where("ca.account_number = ?", accNmbr).
		Select("ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance").
		Joins("join customers as c on ca.customer_number = c.customer_number").
		Table("customer_accounts ca").Find(&accInfo).Error; err != nil {
		fmt.Println(fmt.Sprintf("Failed to get customer account info (%s)", err.Error()))
		return models.Result{Error: err}
	}
	return models.Result{Data: accInfo}
}

func (r *accountManagerRepo) UpdateBalance(account model.CustomerAccount) bool {
	// update sender balance
	if err := r.dbUser.Model(model.CustomerAccount{}).
		Where("account_number = ?", account.AccountNumber).
		Update("account_balance", account.Balance).Error; err != nil {
		fmt.Println(fmt.Sprintf("Failed to update customer balance (%s)", err.Error()))
		return false
	}

	return true
}
