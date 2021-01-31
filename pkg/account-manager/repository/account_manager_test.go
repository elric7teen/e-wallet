package repository

import (
	"fmt"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
	testdata "linkaja.com/e-wallet/pkg/account-manager/test-data"
)

type Suite struct {
	suite.Suite
	DB  *gorm.DB
	e   *echo.Echo
	m   sync.Mutex
	tdb testdata.TestDB

	rp AccountManagerRepo
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

	s.rp = NewAccountManagerRepo(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	// require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetAccountInfo() {
	s.tdb.ResetDB()
	s.tdb.PopulateDB()
	result := s.rp.GetAccountInfo(1)
	require.NoError(s.T(), result.Error)

	result = s.rp.GetAccountInfo(0)
	require.Error(s.T(), result.Error)
}

func (s *Suite) TestUpdateBalance() {
	s.tdb.ResetDB()
	s.tdb.PopulateDB()

	customerAccount := model.CustomerAccount{}
	custAccA := customerAccount.NewCustomerAccount(1, 1, 1000)
	isSuccess := s.rp.UpdateBalance(*custAccA)
	require.Equal(s.T(), true, isSuccess)

	s.DB.DropTable(model.CustomerAccount{})
	isSuccess = s.rp.UpdateBalance(model.CustomerAccount{})
	require.Equal(s.T(), false, isSuccess)
}

func (s *Suite) TestCheckUserAndPassword() {
	s.tdb.ResetDB()
	s.tdb.PopulateDB()

	customer := &model.Customer{}
	customer = customer.NewCustomer(2, "KEANU", "abcdef")
	result := s.rp.CheckUserAndPassword(customer)
	require.NoError(s.T(), result.Error)

	customer = customer.NewCustomer(0, "AA", "BB")
	result = s.rp.CheckUserAndPassword(customer)
	require.Error(s.T(), result.Error)
}
