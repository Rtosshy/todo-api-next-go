package db_test

import (
	"backend/internal/domain"
	"backend/internal/domain/repository"
	"backend/internal/infra/db"
	"backend/internal/testutil"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type TaskRepositorySuite struct {
	testutil.DBSQLiteSuite
	tr repository.TaskRepository
	ur repository.UserRepository
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositorySuite))
}

func (suite *TaskRepositorySuite) SetupSuite() {
	suite.DBSQLiteSuite.SetupSuite()
	suite.tr = db.NewTaskRepository(suite.DB)
	suite.ur = db.NewUserRepository(suite.DB)
}

func (suite *TaskRepositorySuite) MockDB() sqlmock.Sqlmock {
	mock, mockGormDB := testutil.MockDB()
	suite.tr = db.NewTaskRepository(mockGormDB)
	suite.ur = db.NewUserRepository(mockGormDB)
	return mock
}

func (suite *TaskRepositorySuite) AfterTest(suiteName, testName string) {
	suite.tr = db.NewTaskRepository(suite.DB)
	suite.ur = db.NewUserRepository(suite.DB)
}

func (suite *TaskRepositorySuite) TestTaskRepositoryCRUD() {
	user := &domain.User{
		ID:    1,
		Email: "test@test.com",
	}

	user, _ = suite.ur.Create(user)

	task := &domain.Task{
		Name:   "test",
		Status: domain.Status{Name: domain.StatusName("todo")},
		User:   *user,
	}

	// test create
	task, err := suite.tr.Create(task)
	suite.Assert().Nil(err)
	suite.Assert().NotZero(task.ID)
	suite.Assert().Equal("test", task.Name)
	suite.Assert().NotZero(task.Status.ID)
	suite.Assert().Equal(domain.StatusName("todo"), task.Status.Name)
	suite.Assert().NotZero(task.User.ID)
	suite.Assert().Equal("test@test.com", task.User.Email)

	// test get
	getTask, err := suite.tr.Get(task.ID, user.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal("test", getTask.Name)
	suite.Assert().NotZero(getTask.Status.ID)
	suite.Assert().Equal(domain.StatusName("todo"), getTask.Status.Name)

	// test get all
	// getTasks, err := suite.tr.GetAll()

	// test save
	getTask.Name = "updated"
	updatedTask, err := suite.tr.Save(getTask)
	suite.Assert().Nil(err)
	suite.Assert().Equal("updated", updatedTask.Name)
	suite.Assert().NotZero(updatedTask.Status.ID)
	suite.Assert().Equal(domain.StatusName("todo"), updatedTask.Status.Name)

	// test delete
	err = suite.tr.Delete(updatedTask.ID, updatedTask.UserID)
	suite.Assert().Nil(err)
	deletedTask, err := suite.tr.Get(updatedTask.ID, updatedTask.UserID)
	suite.Assert().Nil(deletedTask)
	suite.Assert().True(strings.Contains("record not found", err.Error()))
}

func (suite *TaskRepositorySuite) TestTaskCreateFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "statuses" WHERE "statuses"."name" = $1 ORDER BY "statuses"."id" LIMIT $2`)).WithArgs("todo", 1).WillReturnError(errors.New("create error"))

	task := &domain.Task{
		Name:   "test",
		Status: domain.Status{Name: domain.StatusName("todo")},
	}

	createdTask, err := suite.tr.Create(task)
	suite.Assert().Nil(createdTask)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("create error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskGetFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE id = $1 AND user_id = $2 ORDER BY "tasks"."id" LIMIT $3`)).WithArgs(1, 1, 1).WillReturnError(errors.New("get error"))

	task, err := suite.tr.Get(1, 1)
	suite.Assert().Nil(task)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("get error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskSaveFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks" WHERE id = $1 AND user_id = $2 ORDER BY "tasks"."id" LIMIT $3`)).WithArgs(1, 1, 1).WillReturnError(errors.New("save error"))

	task := &domain.Task{
		ID:     1,
		Name:   "test",
		Status: domain.Status{Name: domain.StatusName("todo")},
		UserID: 1,
	}

	task, err := suite.tr.Save(task)
	suite.Assert().Nil(task)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("save error", err.Error())
}

func (suite *TaskRepositorySuite) TestTaskDeleteFailure() {
	mockDB := suite.MockDB()
	mockDB.ExpectBegin()
	mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tasks" WHERE id = $1 AND user_id = $2`)).WithArgs(1, 1).WillReturnError(errors.New("delete error"))
	mockDB.ExpectRollback()
	mockDB.ExpectCommit()

	err := suite.tr.Delete(1, 1)
	suite.Assert().NotNil(err)
	suite.Assert().Equal("delete error", err.Error())
}
