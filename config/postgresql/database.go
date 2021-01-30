package postgresql

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
)

// CreateDBConnection function for creating database connection
func CreateDBConnection(dialect, descriptor, maxConnPool string) *gorm.DB {
	fmt.Println(descriptor)
	db, err := gorm.Open(dialect, descriptor)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Successfully Connected!")
	}

	db.DB().SetMaxIdleConns(0)
	if maxConn, _ := strconv.Atoi(maxConnPool); maxConn == 0 {
		db.DB().SetMaxOpenConns(10)
	} else {
		db.DB().SetMaxOpenConns(maxConn)
	}

	db.LogMode(true)

	return db
}

// CloseDb function for closing database connection
func CloseDb(db *gorm.DB) {
	if db != nil {
		db.Close()
		db = nil
	}
}
