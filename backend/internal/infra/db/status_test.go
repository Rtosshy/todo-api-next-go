package db_test

import (
	"context"
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
	suite.sr = db.NewStatusRepository(db.NewBaseRepository(suite.DB))
}

func (suite *StatusRepositorySuite) TestStatus() {
	ctx := context.Background()
	cases := []struct {
		name string
		id   domain.StatusID
	}{
		{"todo", 1},
		{"inProgress", 2},
		{"done", 3},
		{"archive", 4},
		{"pending", 5},
	}

	for _, c := range cases {
		paramStatus, err := domain.NewStatus(c.name)
		suite.Assert().Nil(err)
		expectedName, err := domain.NewStatusName(c.name)
		suite.Assert().Nil(err)

		status, err := suite.sr.GetOrCreate(ctx, &paramStatus)
		suite.Assert().Nil(err)
		suite.Assert().Equal(c.id, status.ID)
		suite.Assert().True(expectedName.Equals(status.Name))

		status, err = suite.sr.GetOrCreate(ctx, &paramStatus)
		suite.Assert().Nil(err)
		suite.Assert().Equal(c.id, status.ID)
		suite.Assert().True(expectedName.Equals(status.Name))
	}
}
