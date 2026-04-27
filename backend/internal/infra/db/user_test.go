package db_test

import (
	"backend/internal/domain"
	"backend/internal/domain/repo"
	"backend/internal/infra/db/postgres"
	"backend/internal/testutil"
	"backend/pkg"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type UserRepositorySuite struct {
	testutil.DBSQLiteSuite
	ur repo.UserRepo
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.ur = postgres.NewUserRepository(suite.DB)
}

func (suite *UserRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := testutil.MockDB()
	suite.ur = postgres.NewUserRepository(mockGormDB)
	return mock
}

func (suite *UserRepositorySuite) AfterTest(suiteName, testName string) {
	suite.ur = postgres.NewUserRepository(suite.DB)
}

func (suite *UserRepositorySuite) TestUserRepositoryCRUD() {
	now := pkg.Str2time("2025-01-01")
	user := &domain.User{
		Email:     "test@test.com",
		CreatedAt: now,
	}
	user, err := suite.ur.Create(user)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(user.ID)
	suite.Assert().Equal("test@test.com", user.Email)
	suite.Assert().Equal(now, user.CreatedAt)

	getUser, err := suite.ur.Get(user.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test@test.com", getUser.Email)
	suite.Assert().Equal(now, getUser.CreatedAt)

	getUser, err = suite.ur.GetByEmail(user.Email)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test@test.com", getUser.Email)
	suite.Assert().Equal(now, getUser.CreatedAt)

	getUser.Email = "updated@updated.com"
	updatedUser, err := suite.ur.Save(getUser)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated@updated.com", updatedUser.Email)

	err = suite.ur.Delete(updatedUser.ID)
	suite.Assert().Nil(err)
	deletedUser, err := suite.ur.Get(updatedUser.ID)
	suite.Assert().Nil(deletedUser)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *UserRepositorySuite) TestUserCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","created_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs("test@test.com", "", sqlmock.AnyArg()).
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	user := &domain.User{
		Email: "test@test.com",
	}

	createdUser, err := suite.ur.Create(user)
	suite.Assert().Nil(createdUser)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *UserRepositorySuite) TestUserGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(errors.New("get error"))

	user, err := suite.ur.Get(1)
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *UserRepositorySuite) TestUserGetByEmail() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs("test@test.com", 1).WillReturnError(errors.New("get by email error"))

	user, err := suite.ur.GetByEmail("test@test.com")
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get by email error", err.Error())
}

func (suite *UserRepositorySuite) TestUserSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(errors.New("save error"))

	user := &domain.User{
		ID:    1,
		Email: "test@test.com",
	}

	user, err := suite.ur.Save(user)
	suite.Assert().Nil(user)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("save error", err.Error())
}

func (suite *UserRepositorySuite) TestUserDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).WithArgs(1).WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	err := suite.ur.Delete(1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
