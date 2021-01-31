package usecase

import (
	"database/sql"
	"fmt"
	"regexp"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	models "linkaja.com/e-wallet/lib/base_models"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
	"linkaja.com/e-wallet/pkg/account-manager/model/dto"
	"linkaja.com/e-wallet/pkg/account-manager/repository"
	testdata "linkaja.com/e-wallet/pkg/account-manager/test-data"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	e    *echo.Echo
	m    sync.Mutex
	tdb  testdata.TestDB

	rp repository.AccountManagerRepo
	uc AccountManagerUsecase
}

var (
	custAccFiled    = []string{"account_number", "customer_number", "account_balance"}
	accInfoFiled    = []string{"account_number", "customer_name", "account_balance"}
	customerAccount = model.CustomerAccount{}
	customer        model.Customer
)

func (s *Suite) useSQLite() {
	mockDBUser, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully Connected!")
	}
	mockDBUser.Exec("PRAGMA foreign_keys = ON") // SQLite defaults to `foreign_keys = off'`

	s.tdb = testdata.NewTestDB(mockDBUser)
	s.DB = mockDBUser

	s.rp = repository.NewAccountManagerRepo(s.DB)
	s.uc = NewAccountManagerUsecase(s.rp)
}

func (s *Suite) useSQLMock() {
	var (
		err error
		db  *sql.DB
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.rp = repository.NewAccountManagerRepo(s.DB)
	s.uc = NewAccountManagerUsecase(s.rp)
}

func (s *Suite) SetupSuite() {
	// create mock db user using sqllite3
	// file::memory:?cache=shared >>> prevent table not found
	mockDBUser, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully Connected!")
	}
	mockDBUser.Exec("PRAGMA foreign_keys = ON") // SQLite defaults to `foreign_keys = off'`

	s.tdb = testdata.NewTestDB(mockDBUser)
	s.DB = mockDBUser

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
	s.tdb.ResetDB()
	s.tdb.PopulateDB()
	result := s.uc.ViewAccountInfo(1)
	require.NoError(s.T(), result.Error)

	// failure on database acces
	s.DB.DropTable(model.CustomerAccount{})
	result = s.uc.ViewAccountInfo(0)
	require.Error(s.T(), result.Error)
}

func (s *Suite) TestTransferCredit() {
	s.tdb.ResetDB()
	s.tdb.PopulateDB()

	// transfer success
	responseCode := s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.Success, responseCode)

	s.tdb.ResetDB()
	// not found layer 1
	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.NotFound, responseCode)

	s.tdb.PopulateDB()
	s.DB.Where("account_number = ? ", 2).Delete(&customerAccount)
	// not found layer 2
	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.NotFound, responseCode)

	s.tdb.PopulateDB()
	// insufficient balance
	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          5000,
	})
	require.Equal(s.T(), models.InsufBalance, responseCode)

	// SPECIAL CASE : cant test using sqlite database mock
	// use sqlmock instead
	s.useSQLMock()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(1, "John", 1000))

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT ca.account_number as account_number, c.customer_name as customer_name, ca.account_balance as account_balance FROM customer_accounts ca join customers as c on ca.customer_number = c.customer_number WHERE (ca.account_number = $1)`)).
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows(accInfoFiled).AddRow(2, "Bob", 2000))

	// general error
	responseCode = s.uc.TransferCredit(dto.Param{
		ToAccountNmbr:   2,
		FromAccountNmbr: 1,
		Amount:          1000,
	})
	require.Equal(s.T(), models.InternalServerError, responseCode)

	// return mock db to sqlite
	s.useSQLite()
}

func (s *Suite) Test_IsUserExist() {
	s.tdb.ResetDB()
	s.tdb.PopulateDB()
	// test user exist
	custA := customer.NewCustomer(1, "UNCLE BOB", "abcdef")
	require.True(s.T(), s.uc.IsUserExist(custA))
	// test user not exist
	custB := customer.NewCustomer(2, "Jack", "abcdef")
	require.False(s.T(), s.uc.IsUserExist(custB))
	// test failure on database access
	s.DB.DropTableIfExists(model.Customer{})
	require.False(s.T(), s.uc.IsUserExist(custB))
}
