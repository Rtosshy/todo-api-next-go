package handler

import (
	"backend/internal/domain"
	"backend/internal/usecase"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserUseCase struct {
	mock.Mock
}

func NewMockUserUseCase() usecase.UserUsecase {
	return &MockUserUseCase{}
}

func (m *MockUserUseCase) SignUp(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) Login(user *domain.User) (string, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *MockUserUseCase) Save(user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) Delete(userID domain.UserID) error {
	args := m.Called(userID)
	return args.Error(0)
}

type UserHandlerSuite struct {
	suite.Suite
	uh UserHandler
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerSuite))
}

func (suite *UserHandlerSuite) TestSignUp() {}

func (suite *UserHandlerSuite) TestLogin() {}

func (suite *UserHandlerSuite) TestLogout() {}

func (suite *UserHandlerSuite) TestCsrfToken() {}
