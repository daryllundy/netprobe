package discovery

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/yourusername/netprobe/internal/logger"
)

type Discovery struct {
	logger logger.Logger
	opts   Options
}

type Options struct {
	Timeout    time.Duration
	Concurrent int
	Logger     logger.Logger
}

type HostResult struct {
	IP       string
	Hostname string
	MAC      string
}

func New(opts Options) *Discovery {
	if opts.Timeout <= 0 {
		opts.Timeout = 2 * time.Second
	}
	if opts.Concurrent <= 0 {
		opts.Concurrent = 10
	}
	return &Discovery{logger: opts.Logger, opts: opts}
}

func (d *Discovery) DiscoverInterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		d.logger.Error("Failed to discover interfaces", "error", err)
		return nil, err
	}
	d.logger.Info("Discovered network interfaces", "count", len(interfaces))
	return interfaces, nil
}

func (d *Discovery) DiscoverHosts(ctx context.Context, subnet string) ([]HostResult, error) {
	// Parse subnet (e.g., "192.168.1.0/24")
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet: %v", err)
	}

	var hosts []HostResult
	sem := make(chan struct{}, d.opts.Concurrent)

	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		select {
		case <-ctx.Done():
			return hosts, ctx.Err()
		default:
		}

		sem <- struct{}{}
		go func(ip net.IP) {
			defer func() { <-sem }()
			if host := d.pingHost(ip.String()); host != nil {
				hosts = append(hosts, *host)
			}
		}(dupIP(ip))
	}

	// Wait for all goroutines
	for i := 0; i < cap(sem); i++ {
		sem <- struct{}{}
	}

	d.logger.Info("Discovered hosts", "subnet", subnet, "count", len(hosts))
	return hosts, nil
}

func (d *Discovery) pingHost(ip string) *HostResult {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	if err != nil {
		return nil
	}

	// Get hostname
	names, err := net.LookupAddr(ip)
	hostname := ""
	if err == nil && len(names) > 0 {
		hostname = strings.TrimSuffix(names[0], ".")
	}

	return &HostResult{
		IP:       ip,
		Hostname: hostname,
	}
}

// Helper functions for IP manipulation
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func dupIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}
