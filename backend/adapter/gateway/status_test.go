package gateway_test

import (
	"backend/adapter/gateway"
	"backend/entity"
	"backend/pkg/tester"
	"backend/usecase"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatusRepositorySuite struct {
	tester.DBSQLiteSuite
	sr usecase.IStatusRepository
}

func TestStatusRepositorySuite(t *testing.T) {
	suite.Run(t, new(StatusRepositorySuite))
}

func (suite *StatusRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.sr = gateway.NewStatusRepository(suite.DB)
}

func (suite *StatusRepositorySuite) TestStatus() {
	paramStatus, err := entity.NewStatus("todo")
	expectedID := entity.StatusID(1)
	expectedName := entity.StatusName("todo")
	suite.Assert().Nil(err)
	status, err := suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = entity.NewStatus("inProgress")
	expectedID = entity.StatusID(2)
	expectedName = entity.StatusName("inProgress")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = entity.NewStatus("done")
	expectedID = entity.StatusID(3)
	expectedName = entity.StatusName("done")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = entity.NewStatus("archive")
	expectedID = entity.StatusID(4)
	expectedName = entity.StatusName("archive")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	paramStatus, err = entity.NewStatus("pending")
	expectedID = entity.StatusID(5)
	expectedName = entity.StatusName("pending")
	suite.Assert().Nil(err)
	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)

	status, err = suite.sr.GetOrCreateStatus(paramStatus)
	suite.Assert().Nil(err)
	suite.Assert().Equal(expectedID, status.ID)
	suite.Assert().Equal(expectedName, status.Name)
}
