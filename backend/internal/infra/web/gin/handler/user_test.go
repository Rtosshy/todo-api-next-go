package handler

import (
	"backend/internal/domain"
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserUseCase struct {
	mock.Mock
}

func NewMockUserUseCase() UserUsecase {
	return &MockUserUseCase{}
}

func (m *MockUserUseCase) SignUp(ctx context.Context, email domain.Email, password domain.PlainPassword) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) Login(ctx context.Context, email domain.Email, password domain.PlainPassword) (string, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *MockUserUseCase) Save(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUseCase) Delete(ctx context.Context, userID domain.UserID) error {
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
