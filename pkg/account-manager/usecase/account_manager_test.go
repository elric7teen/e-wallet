package usecase

import (
	"database/sql"
	"regexp"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	models "linkaja.com/e-wallet/lib/base_models"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
	"linkaja.com/e-wallet/pkg/account-manager/model/dto"
	"linkaja.com/e-wallet/pkg/account-manager/repository"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	e    *echo.Echo
	m    sync.Mutex

	rp repository.AccountManagerRepo
	uc AccountManagerUsecase
}

var (
	custAccFiled = []string{"account_number", "customer_number", "account_balance"}
	accInfoFiled = []string{"account_number", "customer_name", "account_balance"}
	custAcc      = model.CustomerAccount{}
)

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.rp = repository.NewAccountManagerRepo(s.DB)
	s.uc = NewAccountManagerUsecase(s.rp)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestViewAccountInfo() {
	// success view account
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 1000))
	result := s.uc.ViewAccountInfo(1)
	require.NoError(s.T(), result.Error)

	// account not found
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(0).
		WillReturnRows(sqlmock.NewRows(nil))
	result = s.uc.ViewAccountInfo(0)
	require.Error(s.T(), result.Error)
}

func (s *Suite) TestTransferCredit() {
	// transfer success
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 1000))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(2, "Bob", 2000))

	query := regexp.QuoteMeta(`UPDATE "customer_accounts" SET "account_balance" = $1 WHERE (account_number = $2)`)

	custAccB := custAcc.NewCustomerAccount(2, 2, 3000)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WithArgs(custAccB.Balance, custAccB.AccountNumber).WillReturnResult(sqlmock.NewResult(0, 0))
	s.mock.ExpectCommit()

	custAccA := custAcc.NewCustomerAccount(1, 1, 0)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WithArgs(custAccA.Balance, custAccA.AccountNumber).WillReturnResult(sqlmock.NewResult(0, 0))
	s.mock.ExpectCommit()

	responseCode := s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.Success, responseCode)

	// not found layer 1
	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.NotFound, responseCode)

	// not found layer 2
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 1000))

	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.NotFound, responseCode)

	// insufficient balance
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 500))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(2, "Bob", 2000))

	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.InsufBalance, responseCode)

	// general error
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 1000))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(2, "Bob", 2000))

	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.InternalServerError, responseCode)
}
