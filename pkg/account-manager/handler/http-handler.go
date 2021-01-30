package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	models "linkaja.com/e-wallet/lib/base_models"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
	"linkaja.com/e-wallet/pkg/account-manager/model/dto"
	"linkaja.com/e-wallet/pkg/account-manager/usecase"
)

// HTTPHandler :S
type HTTPHandler struct {
	accountManagerUC usecase.AccountManagerUsecase
}

// NewAccountManagerHandler : account manager handler builder
func NewAccountManagerHandler(accountManagerUC usecase.AccountManagerUsecase) *HTTPHandler {
	return &HTTPHandler{
		accountManagerUC: accountManagerUC,
	}
}

// Mount :
func (h *HTTPHandler) Mount(group *echo.Group) {
	group.GET("/account/:account_number", h.AccountInfo)
	group.POST("/account/:account_number/transfer", h.Transfer)
	group.GET("/account/login", h.Login)
}

// Login :
func (h *HTTPHandler) Login(c echo.Context) error {
	var cust model.Customer
	if err := c.Bind(cust); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return err
	}
	return nil
}

// AccountInfo :
func (h *HTTPHandler) AccountInfo(c echo.Context) error {
	accNmbr, err := strconv.Atoi(c.Param("account_number"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseBadReq("invalid account number"))
	}
	result := h.accountManagerUC.ViewAccountInfo(accNmbr)
	if result.Error != nil {
		switch result.Error.Error() {
		case "record not found":
			return c.JSON(http.StatusNotFound, models.ResponseNotFound("account not found"))
		default:
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, result.Data)
}

// Transfer :
func (h *HTTPHandler) Transfer(c echo.Context) error {
	fromAccNmbr, err := strconv.Atoi(c.Param("account_number"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseBadReq("invalid account number"))
	}
	req := new(dto.ReqTransferParam)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseBadReq("invalid request body"))
	}

	if fromAccNmbr == req.ToAccountNmbr {
		return c.JSON(http.StatusBadRequest, models.ResponseBadReq("invalid request, can't transfer to same account"))
	}

	param := dto.Param{
		FromAccountNmbr: fromAccNmbr,
		ToAccountNmbr:   req.ToAccountNmbr,
		Amount:          req.Amount,
	}
	switch h.accountManagerUC.TransferCredit(param) {
	case models.BadRequest:
		return c.JSON(http.StatusBadRequest, models.ResponseBadReq("invalid request"))
	case models.InsufBalance:
		return c.JSON(502, models.ResponseInsufBalance("insufficient balance"))
	case models.NotFound:
		return c.JSON(http.StatusNotFound, models.ResponseNotFound("account not found"))
	case models.Success:
		return c.NoContent(201)
	default:
		return c.NoContent(http.StatusInternalServerError)
	}
}
