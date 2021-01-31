package usecase

import (
	models "linkaja.com/e-wallet/lib/base_models"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
	"linkaja.com/e-wallet/pkg/account-manager/model/dto"
	"linkaja.com/e-wallet/pkg/account-manager/repository"
)

type accountManagerUsecase struct {
	repo repository.AccountManagerRepo
}

// NewAccountManagerUsecase : account manager usecase builder
func NewAccountManagerUsecase(repo repository.AccountManagerRepo) AccountManagerUsecase {
	return &accountManagerUsecase{
		repo: repo,
	}
}

func (uc accountManagerUsecase) IsUserExist(cust *model.Customer) bool {
	status := true
	result := uc.repo.CheckUserAndPassword(cust)
	switch customer := result.Data.(type) {
	case *model.Customer:
		if customer == cust {
			status = true
		}
	default:
		status = false
	}
	return status
}

func (uc accountManagerUsecase) ViewAccountInfo(accNmbr int) models.Result {
	result := uc.repo.GetAccountInfo(accNmbr)
	if result.Error != nil {
		return models.Result{Error: result.Error}
	}

	return models.Result{Data: result.Data.(model.AccountInfo)}
}

func (uc accountManagerUsecase) TransferCredit(param dto.Param) string {
	sender := model.AccountInfo{}
	reciever := model.AccountInfo{}

	if result := uc.repo.GetAccountInfo(param.FromAccountNmbr); result.Error == nil {
		if result.Data != nil {
			sender = result.Data.(model.AccountInfo)
		}
	} else {
		return models.NotFound
	}

	if result := uc.repo.GetAccountInfo(param.ToAccountNmbr); result.Error == nil {
		if result.Data != nil {
			reciever = result.Data.(model.AccountInfo)
		}
	} else {
		return models.NotFound
	}

	if sender.Balance >= param.Amount {
		reciever.Balance += param.Amount
		sender.Balance -= param.Amount
	} else {
		return models.InsufBalance
	}

	// update reciever
	if isSuccess := uc.repo.UpdateBalance(model.CustomerAccount{
		AccountNumber: reciever.AccountNumber,
		Balance:       reciever.Balance}); isSuccess {
		// update sender
		uc.repo.UpdateBalance(model.CustomerAccount{
			AccountNumber: sender.AccountNumber,
			Balance:       sender.Balance})
	} else {
		return models.InternalServerError
	}

	return models.Success
}
