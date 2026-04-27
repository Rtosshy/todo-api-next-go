package db_test

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"
	"backend/internal/infra/db/postgres"
	"backend/internal/testutil"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatusRepositorySuite struct {
	testutil.DBSQLiteSuite
	sr repo.StatusRepo
}

func TestStatusRepositorySuite(t *testing.T) {
	suite.Run(t, new(StatusRepositorySuite))
}

func (suite *StatusRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.sr = postgres.NewStatusRepository(suite.DB)
}

func (suite *StatusRepositorySuite) TestStatus() {
	paramStatus, err := domain.NewStatus("todo")
	expectedID := domain.StatusID(1)
	expectedName := domain.StatusName("todo")
	suite.Assert().Nil(err)
	status, err := suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = domain.NewStatus("inProgress")
	expectedID = domain.StatusID(2)
	expectedName = domain.StatusName("inProgress")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = domain.NewStatus("done")
	expectedID = domain.StatusID(3)
	expectedName = domain.StatusName("done")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = domain.NewStatus("archive")
	expectedID = domain.StatusID(4)
	expectedName = domain.StatusName("archive")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = domain.NewStatus("pending")
	expectedID = domain.StatusID(5)
	expectedName = domain.StatusName("pending")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreate(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)
}
