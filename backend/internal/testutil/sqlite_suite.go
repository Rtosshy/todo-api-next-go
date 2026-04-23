package testutil

import (
	"backend/internal/domain"
	"backend/internal/infra/db"
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DBSQLiteSuite struct {
	suite.Suite
	DB     *gorm.DB
	DBName string
}

func (suite *DBSQLiteSuite) SetupSuite() {
	suite.DBName = fmt.Sprintf("%s.unittest.sqlite", suite.T().Name())

	os.Setenv("DB_NAME", suite.DBName)
	db, err := database.NewDatabaseSQLFactory(database.InstanceSQLite)
	suite.Assert().Nil(err)
	suite.DB = db

	for _, model := range domain.NewDomains() {
		err := suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBSQLiteSuite) TearDownSuite() {
	err := os.Remove(suite.DBName)
	suite.Assert().Nil(err)
	os.Unsetenv(suite.DBName)
}
