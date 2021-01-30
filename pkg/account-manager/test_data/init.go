package testdata

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	model "linkaja.com/e-wallet/pkg/account-manager/model/db"
)

// TestDB is an interface to help mock database using sqllite
type TestDB interface {
	ResetDB()
	PopulateDB()
}

type testDB struct {
	dbUser *gorm.DB
}

// NewTestDB is a TestDB factory
func NewTestDB(dbUser *gorm.DB) TestDB {
	return &testDB{
		dbUser: dbUser,
	}
}

func (tdb *testDB) dropAllTable() {
	tdb.dbUser.DropTableIfExists(&model.Customer{})
	tdb.dbUser.DropTableIfExists(&model.CustomerAccount{})
}

func (tdb *testDB) createAllTable() {
	tdb.dbUser.AutoMigrate(&model.Customer{})
	tdb.dbUser.AutoMigrate(&model.CustomerAccount{})
}

// clear all table
func (tdb *testDB) ResetDB() {
	tdb.dropAllTable()
	tdb.createAllTable()
}

// populate database with mock data provided in db.sql file
func (tdb *testDB) PopulateDB() {
	if err := tdb.runSQLFile(getSQLFile()); err != nil {
		panic(fmt.Errorf("error while initializing test database: %s", err))
	}
}

func getSQLFile() string {
	return "/media/ghozi/6C908BD5908BA3E4/work/backup_laptop_2_aino/go/e-wallet/pkg/account-manager/test_data/db.sql"
}

func (tdb *testDB) runSQLFile(filepath string) error {
	src, err := os.Open(filepath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		sql := scanner.Text()
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}

		if result := tdb.dbUser.Exec(sql); result.Error != nil {
			fmt.Println(sql)
			return result.Error
		}
	}

	return nil
}
