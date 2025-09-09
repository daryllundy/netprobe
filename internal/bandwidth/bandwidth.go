package bandwidth

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/daryllundy/netprobe/internal/logger"
)

type Tester struct {
	options Options
}

type Options struct {
	Timeout      time.Duration
	BufferSize   int
	TestDuration time.Duration
	Protocol     string // "tcp" or "udp"
	Logger       logger.Logger
}

type Result struct {
	UploadSpeed   float64 // Mbps
	DownloadSpeed float64 // Mbps
	Latency       time.Duration
	Error         error
}

func New(opts Options) *Tester {
	if opts.BufferSize <= 0 {
		opts.BufferSize = 64 * 1024 // 64KB
	}
	if opts.TestDuration <= 0 {
		opts.TestDuration = 10 * time.Second
	}
	if opts.Protocol == "" {
		opts.Protocol = "tcp"
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 30 * time.Second
	}
	return &Tester{options: opts}
}

func (t *Tester) Test(ctx context.Context, host string, port int) (*Result, error) {
	result := &Result{}

	// Test latency first
	latency, err := t.measureLatency(host, port)
	if err != nil {
		result.Error = fmt.Errorf("latency test failed: %w", err)
		return result, result.Error
	}
	result.Latency = latency

	// Test download speed
	downloadSpeed, err := t.testDownload(ctx, host, port)
	if err != nil {
		result.Error = fmt.Errorf("download test failed: %w", err)
		return result, result.Error
	}
	result.DownloadSpeed = downloadSpeed

	// Test upload speed
	uploadSpeed, err := t.testUpload(ctx, host, port)
	if err != nil {
		result.Error = fmt.Errorf("upload test failed: %w", err)
		return result, result.Error
	}
	result.UploadSpeed = uploadSpeed

	return result, nil
}

func (t *Tester) measureLatency(host string, port int) (time.Duration, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	start := time.Now()
	conn, err := net.DialTimeout(t.options.Protocol, address, t.options.Timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	return time.Since(start), nil
}

func (t *Tester) testDownload(ctx context.Context, host string, port int) (float64, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout(t.options.Protocol, address, t.options.Timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	buffer := make([]byte, t.options.BufferSize)
	totalBytes := 0
	start := time.Now()

	// Set a timeout for individual read operations
	readWriteTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, readWriteTimeout)
	defer cancel()

	for time.Since(start) < t.options.TestDuration {
		select {
		case <-ctx.Done():
			// If the main context is cancelled, or read/write times out
			if ctx.Err() == context.DeadlineExceeded {
				return 0, fmt.Errorf("download read operation timed out after %v", readWriteTimeout)
			}
			return 0, ctx.Err()
		default:
		}

		// Perform a non-blocking read or use a deadline for the read call itself.
		// Using SetReadDeadline is a more direct approach for this.
		deadline := time.Now().Add(readWriteTimeout)
		if err := conn.SetReadDeadline(deadline); err != nil {
			return 0, fmt.Errorf("failed to set read deadline: %w", err)
		}

		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return 0, fmt.Errorf("download read operation timed out after %v", readWriteTimeout)
			}
			if err == io.EOF {
				break // End of stream
			}
			return 0, err // Other read errors
		}
		totalBytes += n
	}

	duration := time.Since(start).Seconds()
	if duration == 0 { // Avoid division by zero if test was extremely short or no data transferred
		duration = 1 // Treat as 1 second for Mbps calculation, or handle as 0 Mbps
	}
	mbps := float64(totalBytes*8) / (1000000 * duration)

	return mbps, nil
}

func (t *Tester) testUpload(ctx context.Context, host string, port int) (float64, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout(t.options.Protocol, address, t.options.Timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	buffer := make([]byte, t.options.BufferSize)
	// Fill buffer with test data
	for i := range buffer {
		buffer[i] = byte(i % 256)
	}

	totalBytes := 0
	start := time.Now()

	// Set a timeout for individual write operations
	readWriteTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, readWriteTimeout)
	defer cancel()

	for time.Since(start) < t.options.TestDuration {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				return 0, fmt.Errorf("upload write operation timed out after %v", readWriteTimeout)
			}
			return 0, ctx.Err()
		default:
		}

		// Perform a non-blocking write or use a deadline for the write call itself.
		deadline := time.Now().Add(readWriteTimeout)
		if err := conn.SetWriteDeadline(deadline); err != nil {
			return 0, fmt.Errorf("failed to set write deadline: %w", err)
		}

		n, err := conn.Write(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return 0, fmt.Errorf("upload write operation timed out after %v", readWriteTimeout)
			}
			return 0, err // Other write errors
		}
		totalBytes += n
	}

	duration := time.Since(start).Seconds()
	if duration == 0 { // Avoid division by zero
		duration = 1
	}
	mbps := float64(totalBytes*8) / (1000000 * duration)

	return mbps, nil
}
