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

type TaskRepositorySuite struct {
	testutil.DBSQLiteSuite
	tr repository.TaskRepository
	ur repository.UserRepository
	sr repository.StatusRepository
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}

func (suite *TaskRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.tr = db.NewTaskRepository(db.NewBaseRepository(suite.DB))
	suite.ur = db.NewUserRepository(db.NewBaseRepository(suite.DB))
	suite.sr = db.NewStatusRepository(db.NewBaseRepository(suite.DB))
}

func (suite *TaskRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := testutil.MockDB()
	suite.tr = db.NewTaskRepository(db.NewBaseRepository(mockGormDB))
	suite.ur = db.NewUserRepository(db.NewBaseRepository(mockGormDB))
	suite.sr = db.NewStatusRepository(db.NewBaseRepository(mockGormDB))
	return mock
}

func (suite *TaskRepositorySuite) AfterTest(suiteName, testName string) {
	suite.tr = db.NewTaskRepository(db.NewBaseRepository(suite.DB))
	suite.ur = db.NewUserRepository(db.NewBaseRepository(suite.DB))
	suite.sr = db.NewStatusRepository(db.NewBaseRepository(suite.DB))
}

func (suite *TaskRepositorySuite) buildUser(emailStr string) *domain.User {
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

func (suite *TaskRepositorySuite) TestTaskRepositoryCRUD() {
	ctx := context.Background()

	user, err := suite.ur.Create(ctx, suite.buildUser("test@test.com"))
	suite.Assert().Nil(err)

	statusParam, err := domain.NewStatus("todo")
	suite.Assert().Nil(err)
	resolvedStatus, err := suite.sr.GetOrCreate(ctx, &statusParam)
	suite.Assert().Nil(err)

	taskName, err := domain.NewTaskName("test")
	suite.Assert().Nil(err)
	task, err := domain.NewTask(taskName, *resolvedStatus, user.ID(), nil)
	suite.Assert().Nil(err)

	created, err := suite.tr.Create(ctx, task)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(created.ID())
	suite.Assert().Equal("test", created.Name().String())
	suite.Assert().NotZero(created.Status().ID)
	suite.Assert().Equal("todo", created.Status().Name.String())
	suite.Assert().Equal(user.ID(), created.UserID())

	got, err := suite.tr.Get(ctx, created.ID(), user.ID())
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", got.Name().String())
	suite.Assert().NotZero(got.Status().ID)
	suite.Assert().Equal("todo", got.Status().Name.String())

	updatedName, err := domain.NewTaskName("updated")
	suite.Assert().Nil(err)
	updatedTask := domain.ReconstructTask(got.ID(), updatedName, got.Status(), got.UserID(), got.Deadline())
	updated, err := suite.tr.Save(ctx, updatedTask)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updated.Name().String())
	suite.Assert().NotZero(updated.Status().ID)
	suite.Assert().Equal("todo", updated.Status().Name.String())

	err = suite.tr.Delete(ctx, updated.ID(), updated.UserID())
	suite.Assert().Nil(err)
	deleted, err := suite.tr.Get(ctx, updated.ID(), updated.UserID())
	suite.Assert().Nil(deleted)
	suite.Assert().True(strings.Contains(err.Error(), "record not found"))
}

func (suite *TaskRepositorySuite) TestTaskCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tasks"`)).
		WillReturnError(errors.New("create error"))
	mockDB.ExpectRollback()

	taskName, err := domain.NewTaskName("test")
	suite.Assert().Nil(err)
	statusName, err := domain.NewStatusName("todo")
	suite.Assert().Nil(err)
	task, err := domain.NewTask(taskName, domain.Status{ID: 1, Name: statusName}, 1, nil)
	suite.Assert().Nil(err)

	created, err := suite.tr.Create(context.Background(), task)
	suite.Assert().Nil(created)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE id = $1 AND user_id = $2 ORDER BY "tasks"."id" LIMIT $3`)).
		WithArgs(1, 1, 1).
		WillReturnError(errors.New("get error"))

	task, err := suite.tr.Get(context.Background(), 1, 1)
	suite.Assert().Nil(task)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE id = $1 AND user_id = $2`)).
		WithArgs(1, 1).
		WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()

	err := suite.tr.Delete(context.Background(), 1, 1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
