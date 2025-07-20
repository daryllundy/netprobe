# NetProbe

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org) [![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A federal-compliant network analysis tool for restricted and air-gapped environments.

## Features

- Port Scanning with service detection
- Network Discovery and monitoring
- DNS Analysis and testing
- Bandwidth and latency measurements
- NIST 800-53 compliant logging
- FIPS 140-2 encryption support
- Zero runtime dependencies

## Installation

```bash
# Clone and build
git clone https://github.com/daryllundy/netprobe.git
cd netprobe
make build

# Install
make install
```

## Quick Start

```bash
# Basic port scan
netprobe scan --host 192.168.1.1 --ports 1-1000

# Monitor connections
netprobe monitor --duration 60s

# Generate compliance report
netprobe report --format pdf --output audit.pdf
```

## Documentation

See the [docs](docs/) directory for detailed documentation.

## License

MIT License - see [LICENSE](LICENSE) file for details.
