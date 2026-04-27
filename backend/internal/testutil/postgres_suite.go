package testutil

import (
	"backend/internal/infra/db"
	"backend/internal/infra/db/dao"
	"backend/pkg"
	"context"
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

type DBPostgresSuite struct {
	suite.Suite
	postgresContainer testcontainers.Container
	ctx               context.Context
	DB                *gorm.DB
}

func (suite *DBPostgresSuite) SetupTestContainers() (err error) {
	configs := database.NewConfigPostgres()
	pkg.WaitForPort(configs.Database, configs.Port, 10*time.Second)
	suite.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image: "postgres:latest",
		Env: map[string]string{
			"POSTGRES_DB":       configs.Database,
			"POSTGRES_USER":     configs.User,
			"POSTGRES_PASSWORD": configs.Password,
		},
		ExposedPorts: []string{fmt.Sprintf("%s:5432/tcp", configs.Port)},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	suite.postgresContainer, err = testcontainers.GenericContainer(
		suite.ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	suite.Assert().Nil(err)
	return nil
}

func (suite *DBPostgresSuite) SetupSuite() {
	err := suite.SetupTestContainers()
	suite.Assert().Nil(err)

	db, err := database.NewDatabaseSQLFactory(database.InstancePostgres)
	suite.Assert().Nil(err)
	suite.DB = db
	for _, model := range dao.NewDAOs() {
		err = suite.DB.AutoMigrate(model)
		suite.Assert().Nil(err)
	}
}

func (suite *DBPostgresSuite) TearDownSuite() {
	if suite.postgresContainer == nil {
		return
	}
	err := suite.postgresContainer.Terminate(suite.ctx)
	suite.Assert().Nil(err)
}
