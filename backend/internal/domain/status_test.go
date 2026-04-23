package domain_test

import (
	"backend/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	status := domain.Status{
		ID:   1,
		Name: "todo",
	}
	assert.Equal(t, domain.StatusID(1), status.ID)
	assert.Equal(t, domain.StatusName("todo"), status.Name)
}
