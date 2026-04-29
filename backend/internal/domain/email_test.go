package domain_test

import (
	"backend/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmail(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.co.jp", false},
		{"empty", "", true},
		{"missing @", "userexample.com", true},
		{"missing local part", "@example.com", true},
		{"missing domain", "user@", true},
		{"spaces only", "   ", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			email, err := domain.NewEmail(tc.input)
			if tc.wantErr {
				assert.ErrorIs(t, err, domain.ErrInvalidEmail)
				assert.True(t, email.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, email.String())
				assert.False(t, email.IsZero())
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	a, _ := domain.NewEmail("user@example.com")
	b, _ := domain.NewEmail("user@example.com")
	c, _ := domain.NewEmail("other@example.com")

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}

func TestEmail_ZeroValue(t *testing.T) {
	var e domain.Email
	assert.True(t, e.IsZero())
	assert.Equal(t, "", e.String())
}
