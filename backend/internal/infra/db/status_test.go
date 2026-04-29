package db_test

import (
	"testing"

	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db"
	"backend/internal/testutil"

	"github.com/stretchr/testify/suite"
)

type StatusRepositorySuite struct {
	testutil.DBSQLiteSuite
	sr repository.StatusRepository
}

func TestStatusRepositorySuite(t *testing.T) {
	suite.Run(t, new(StatusRepositorySuite))
}

func (suite *StatusRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.sr = db.NewStatusRepository(suite.DB)
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
