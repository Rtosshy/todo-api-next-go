package domain_test

import (
	"backend/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlainPassword(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid 8 chars", "12345678", false},
		{"valid 20 chars", "abcdefghij1234567890", false},
		{"too short 7 chars", "1234567", true},
		{"empty", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := domain.NewPlainPassword(tc.input)
			if tc.wantErr {
				assert.ErrorIs(t, err, domain.ErrPasswordTooShort)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlainPassword_HashAndMatches(t *testing.T) {
	plain, err := domain.NewPlainPassword("password123")
	assert.NoError(t, err)

	hashed, err := plain.Hash()
	assert.NoError(t, err)
	assert.False(t, hashed.IsZero())
	assert.NotEqual(t, "password123", hashed.String())

	assert.True(t, plain.Matches(hashed))

	wrong, _ := domain.NewPlainPassword("different")
	assert.False(t, wrong.Matches(hashed))
}

func TestNewHashedPassword(t *testing.T) {
	t.Run("non-empty value succeeds", func(t *testing.T) {
		h, err := domain.NewHashedPassword("$2a$10$abcdef")
		assert.NoError(t, err)
		assert.False(t, h.IsZero())
		assert.Equal(t, "$2a$10$abcdef", h.String())
	})

	t.Run("empty value fails", func(t *testing.T) {
		_, err := domain.NewHashedPassword("")
		assert.ErrorIs(t, err, domain.ErrInvalidHashedPassword)
	})
}

func TestHashedPassword_ZeroValue(t *testing.T) {
	var h domain.HashedPassword
	assert.True(t, h.IsZero())
	assert.Equal(t, "", h.String())
}
