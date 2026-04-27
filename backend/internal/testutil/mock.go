package testutil

import (
	"backend/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MockDB() (mock sqlmock.Sqlmock, mockGormDB *gorm.DB) {
	mockDB, mock, err := sqlmock.New(
		sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	mockGormDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:        "mock_db",
		DriverName: "postgres",
		Conn:       mockDB,
	}), &gorm.Config{})
	if err != nil {
		logger.Fatal(err.Error())
	}
	return mock, mockGormDB
}
