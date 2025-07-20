package discovery

import (
	"net"
	"github.com/yourusername/netprobe/internal/logger"
)

type Discovery struct {
	logger logger.Logger
}

func New(log logger.Logger) *Discovery {
	return &Discovery{logger: log}
}

func (d *Discovery) DiscoverInterfaces() ([]net.Interface, error) {
	return net.Interfaces()
}
