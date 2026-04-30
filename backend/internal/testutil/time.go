package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func MustDate(t *testing.T, s string) time.Time {
	t.Helper()
	parsed, err := time.Parse("2006-01-02", s)
	require.NoError(t, err)
	return parsed
}
