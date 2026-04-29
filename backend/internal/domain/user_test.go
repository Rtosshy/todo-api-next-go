package domain_test

import (
	"testing"

	"backend/internal/domain"
	"backend/pkg"

	"github.com/stretchr/testify/assert"
)

func validEmail(t *testing.T) domain.Email {
	t.Helper()
	e, err := domain.NewEmail("test@test.com")
	assert.NoError(t, err)
	return e
}

func validHashedPassword(t *testing.T) domain.HashedPassword {
	t.Helper()
	plain, err := domain.NewPlainPassword("password123")
	assert.NoError(t, err)
	hashed, err := plain.Hash()
	assert.NoError(t, err)
	return hashed
}

func TestNewUser(t *testing.T) {
	t.Run("valid email and password creates user", func(t *testing.T) {
		email := validEmail(t)
		hashed := validHashedPassword(t)

		user, err := domain.NewUser(email, hashed)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.True(t, user.Email().Equals(email))
		assert.Equal(t, hashed.String(), user.Password().String())
	})

	t.Run("zero email rejected", func(t *testing.T) {
		var zeroEmail domain.Email
		hashed := validHashedPassword(t)
		user, err := domain.NewUser(zeroEmail, hashed)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		assert.Nil(t, user)
	})

	t.Run("zero password rejected", func(t *testing.T) {
		email := validEmail(t)
		var zeroHash domain.HashedPassword
		user, err := domain.NewUser(email, zeroHash)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		assert.Nil(t, user)
	})
}

func TestReconstructUser(t *testing.T) {
	now := pkg.Str2time("2025-01-01")
	email := validEmail(t)
	hashed := validHashedPassword(t)

	user := domain.ReconstructUser(domain.UserID(42), email, hashed, now)
	assert.Equal(t, domain.UserID(42), user.ID())
	assert.True(t, user.Email().Equals(email))
	assert.Equal(t, hashed.String(), user.Password().String())
	assert.Equal(t, now, user.CreatedAt())
}
