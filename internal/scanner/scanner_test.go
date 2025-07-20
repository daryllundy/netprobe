package scanner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/netprobe/internal/logger"
)

func TestScanner_Scan(t *testing.T) {
	ctx := context.Background()
	log := logger.NewTestLogger(t)

	scanner := New(Options{
		Timeout:    time.Second,
		Concurrent: 10,
		Logger:     log,
	})

	// Test with localhost
	results, err := scanner.Scan(ctx, []string{"127.0.0.1"}, []int{22, 80, 443})

	require.NoError(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 1)
	assert.Equal(t, "127.0.0.1", results[0].Host)
}

func TestScanner_InvalidInput(t *testing.T) {
	ctx := context.Background()
	log := logger.NewTestLogger(t)

	scanner := New(Options{
		Timeout: time.Second,
		Logger:  log,
	})

	// Test empty hosts
	_, err := scanner.Scan(ctx, []string{}, []int{80})
	assert.Error(t, err)

	// Test empty ports
	_, err = scanner.Scan(ctx, []string{"127.0.0.1"}, []int{})
	assert.Error(t, err)
}
