package monitor

import (
	"context"

	"github.com/daryllundy/netprobe/internal/logger"
)

type Monitor struct {
	logger logger.Logger
}

func New(log logger.Logger) *Monitor {
	return &Monitor{logger: log}
}

func (m *Monitor) Start(ctx context.Context) error {
	m.logger.Info("Starting connection monitor")
	// Implementation here
	return nil
}
