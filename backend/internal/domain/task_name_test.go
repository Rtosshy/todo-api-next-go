package domain_test

import (
	"strings"
	"testing"

	"backend/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewTaskName(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"single char", "x", false},
		{"normal", "Buy milk", false},
		{"japanese", "牛乳を買う", false},
		{"max length", strings.Repeat("a", 255), false},
		{"empty", "", true},
		{"too long", strings.Repeat("a", 256), true},
		{"too long japanese", strings.Repeat("あ", 256), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := domain.NewTaskName(tc.input)
			if tc.wantErr {
				assert.ErrorIs(t, err, domain.ErrInvalidTaskName)
				assert.True(t, n.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, n.String())
			}
		})
	}
}

func TestTaskName_Equals(t *testing.T) {
	a, _ := domain.NewTaskName("foo")
	b, _ := domain.NewTaskName("foo")
	c, _ := domain.NewTaskName("bar")
	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}
