package mock

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func MockDb(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, *sql.DB) {
	db, mock, err1 := sqlmock.New()
	if err1 != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err1)
	}
	d := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	gormDb, err1 := gorm.Open(d, &gorm.Config{})
	if err1 != nil {
		t.Fatalf("an error '%s' with using sqlmock with gorm", err1)
	}
	return mock, gormDb, db
}
