package domain_test

import (
	"backend/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatusName(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"todo", false},
		{"inProgress", false},
		{"done", false},
		{"archive", false},
		{"pending", false},
		{"unknown", true},
		{"", true},
		{"TODO", true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			n, err := domain.NewStatusName(tc.input)
			if tc.wantErr {
				assert.ErrorIs(t, err, domain.ErrInvalidStatusName)
				assert.True(t, n.IsZero())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.input, n.String())
			}
		})
	}
}

func TestStatusName_Constants(t *testing.T) {
	assert.Equal(t, "todo", domain.StatusNameTodo.String())
	assert.Equal(t, "inProgress", domain.StatusNameInProgress.String())
	assert.Equal(t, "done", domain.StatusNameDone.String())
	assert.Equal(t, "archive", domain.StatusNameArchive.String())
	assert.Equal(t, "pending", domain.StatusNamePending.String())
}

func TestStatusName_Equals(t *testing.T) {
	assert.True(t, domain.StatusNameTodo.Equals(domain.StatusNameTodo))
	assert.False(t, domain.StatusNameTodo.Equals(domain.StatusNameDone))
}

func TestNewStatus(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		s, err := domain.NewStatus("todo")
		assert.NoError(t, err)
		assert.Equal(t, "todo", s.Name.String())
	})

	t.Run("invalid name", func(t *testing.T) {
		_, err := domain.NewStatus("invalid")
		assert.ErrorIs(t, err, domain.ErrInvalidStatusName)
	})
}
