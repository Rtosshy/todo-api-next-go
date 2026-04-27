package domain_test

import (
	"backend/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDeadline(t *testing.T) {
	t.Run("valid future time", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		d, err := domain.NewDeadline(future)
		assert.NoError(t, err)
		assert.True(t, d.Time().Equal(future))
	})

	t.Run("valid past time", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		d, err := domain.NewDeadline(past)
		assert.NoError(t, err)
		assert.True(t, d.Time().Equal(past))
	})

	t.Run("zero value rejected", func(t *testing.T) {
		_, err := domain.NewDeadline(time.Time{})
		assert.ErrorIs(t, err, domain.ErrInvalidDeadline)
	})
}

func TestDeadline_IsFuture(t *testing.T) {
	future, _ := domain.NewDeadline(time.Now().Add(time.Hour))
	past, _ := domain.NewDeadline(time.Now().Add(-time.Hour))

	assert.True(t, future.IsFuture())
	assert.False(t, past.IsFuture())
}

func TestDeadline_Equals(t *testing.T) {
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	a, _ := domain.NewDeadline(t1)
	b, _ := domain.NewDeadline(t1)
	c, _ := domain.NewDeadline(t1.Add(time.Hour))

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}
