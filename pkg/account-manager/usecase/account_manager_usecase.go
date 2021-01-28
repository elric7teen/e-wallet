package usecase

import (
	models "linkaja.com/e-wallet/lib/base_models"
	"linkaja.com/e-wallet/pkg/account-manager/model/dto"
)

// AccountManagerUsecase : Account Managager usecase interface
type AccountManagerUsecase interface {
	ViewAccountInfo(accNmbr int) models.Result
	TransferCredit(param dto.Param) (responseCode string)
}
