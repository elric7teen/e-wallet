package repository

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
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	e    *echo.Echo
	m    sync.Mutex

	rp AccountManagerRepo
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

	s.rp = NewAccountManagerRepo(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetAccountInfo() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 1000))
	result := s.rp.GetAccountInfo(1)
	require.NoError(s.T(), result.Error)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(0).
		WillReturnRows(sqlmock.NewRows(nil))
	result = s.rp.GetAccountInfo(0)
	require.Error(s.T(), result.Error)
}

func (s *Suite) TestUpdateBalance() {
	custAccA := custAcc.NewCustomerAccount(1, 1, 1000)

	query := regexp.QuoteMeta(`UPDATE "customer_accounts" SET "account_balance" = $1 WHERE (account_number = $2)`)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WithArgs(custAccA.Balance, custAccA.AccountNumber).WillReturnResult(sqlmock.NewResult(0, 0))
	s.mock.ExpectCommit()

	isSuccess := s.rp.UpdateBalance(*custAccA)
	require.Equal(s.T(), true, isSuccess)

	// s.mock.ExpectBegin()
	// s.mock.ExpectExec(query).WithArgs(custAccA.Balance, custAccA.AccountNumber).WillReturnResult(sqlmock.NewResult(0, 0))
	// s.mock.ExpectCommit()

	custAccA.AccountNumber = 2
	isSuccess = s.rp.UpdateBalance(*custAccA)
	require.Equal(s.T(), false, isSuccess)

}
