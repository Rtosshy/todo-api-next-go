package db_test

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"

	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db"
	"backend/internal/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type UserRepositorySuite struct {
	testutil.DBSQLiteSuite
	ur repository.UserRepository
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.ur = db.NewUserRepository(db.NewBaseRepository(suite.DB))
}

func (suite *UserRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := testutil.MockDB()
	suite.ur = db.NewUserRepository(db.NewBaseRepository(mockGormDB))
	return mock
}

func (suite *UserRepositorySuite) AfterTest(suiteName, testName string) {
	suite.ur = db.NewUserRepository(db.NewBaseRepository(suite.DB))
}

func (suite *UserRepositorySuite) buildUser(emailStr string) *domain.User {
	email, err := domain.NewEmail(emailStr)
	suite.Require().Nil(err)
	plain, err := domain.NewPlainPassword("password123")
	suite.Require().Nil(err)
	hashed, err := plain.Hash()
	suite.Require().Nil(err)
	user, err := domain.NewUser(email, hashed)
	suite.Require().Nil(err)
	return user
}

func (suite *UserRepositorySuite) TestUserRepositoryCRUD() {
	ctx := context.Background()
	user := suite.buildUser("test@test.com")

	created, err := suite.ur.Create(ctx, user)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(created.ID())
	suite.Assert().Equal("test@test.com", created.Email().String())

	got, err := suite.ur.Get(ctx, created.ID())
	suite.Assert().Nil(err)
	suite.Assert().Equal("test@test.com", got.Email().String())

	got, err = suite.ur.GetByEmail(ctx, created.Email().String())
	suite.Assert().Nil(err)
	suite.Assert().Equal("test@test.com", got.Email().String())

	updatedEmail, err := domain.NewEmail("updated@updated.com")
	suite.Assert().Nil(err)
	reUser := domain.ReconstructUser(got.ID(), updatedEmail, got.Password(), got.CreatedAt())
	updated, err := suite.ur.Save(ctx, reUser)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated@updated.com", updated.Email().String())

	err = suite.ur.Delete(ctx, updated.ID())
	suite.Assert().Nil(err)
	deleted, err := suite.ur.Get(ctx, updated.ID())
	suite.Assert().Nil(deleted)
	suite.Assert().True(strings.Contains(err.Error(), "record not found"))
}

func (suite *UserRepositorySuite) TestUserCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	user := suite.buildUser("test@test.com")
	created, err := suite.ur.Create(context.Background(), user)
	suite.Assert().Nil(created)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *UserRepositorySuite) TestUserGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnError(errors.New("get error"))

	user, err := suite.ur.Get(context.Background(), 1)
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *UserRepositorySuite) TestUserGetByEmail() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("test@test.com", 1).
		WillReturnError(errors.New("get by email error"))

	user, err := suite.ur.GetByEmail(context.Background(), "test@test.com")
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get by email error", err.Error())
}

func (suite *UserRepositorySuite) TestUserDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.ur.Delete(context.Background(), 1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
