package scanner

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/daryllundy/netprobe/internal/logger"
)

type Scanner struct {
	options Options
}

type Options struct {
	Timeout       time.Duration
	Concurrent    int
	DetectService bool
	RateLimit     int
	Logger        logger.Logger
}

type ScanResult struct {
	Host  string
	Ports []PortResult
}

type PortResult struct {
	Number  int
	Open    bool
	Service string
}

func New(opts Options) *Scanner {
	if opts.Concurrent <= 0 {
		opts.Concurrent = 100
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 3 * time.Second
	}
	return &Scanner{options: opts}
}

func (s *Scanner) Scan(ctx context.Context, hosts []string, ports []int) ([]ScanResult, error) {
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts provided")
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports provided")
	}

	results := make([]ScanResult, 0, len(hosts))

	for _, host := range hosts {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		hostResult := ScanResult{
			Host:  host,
			Ports: make([]PortResult, 0, len(ports)),
		}

		// Use a map to collect results by port number for ordering
		portResults := make(map[int]PortResult)
		var mapMu sync.Mutex

		// Create semaphore for concurrency control
		sem := make(chan struct{}, s.options.Concurrent)
		var wg sync.WaitGroup

		for _, port := range ports {
			wg.Add(1)
			sem <- struct{}{}

			go func(p int) {
				defer wg.Done()
				defer func() { <-sem }()

				result := s.scanPort(host, p)

				mapMu.Lock()
				portResults[p] = result
				mapMu.Unlock()
			}(port)
		}

		wg.Wait()

		// Sort ports and build the results slice in numeric order
		sortedPorts := make([]int, 0, len(ports))
		for port := range portResults {
			sortedPorts = append(sortedPorts, port)
		}
		sort.Ints(sortedPorts)

		for _, port := range sortedPorts {
			hostResult.Ports = append(hostResult.Ports, portResults[port])
		}

		results = append(results, hostResult)
	}

	return results, nil
}

func (s *Scanner) scanPort(host string, port int) PortResult {
	result := PortResult{
		Number: port,
		Open:   false,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, s.options.Timeout)
	if err != nil {
		return result
	}
	defer conn.Close()

	result.Open = true

	if s.options.DetectService {
		result.Service = s.detectService(conn)
	}

	return result
}

func (s *Scanner) detectService(conn net.Conn) string {
	// Simple service detection based on port numbers
	// In a real implementation, this would do banner grabbing
	localAddr := conn.LocalAddr().String()
	_, portStr, _ := net.SplitHostPort(localAddr)

	switch portStr {
	case "22":
		return "ssh"
	case "80":
		return "http"
	case "443":
		return "https"
	case "3306":
		return "mysql"
	case "5432":
		return "postgresql"
	default:
		return "unknown"
	}
}
