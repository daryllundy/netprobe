# Contributing to NetProbe

Thank you for your interest in contributing to NetProbe! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites
- Go 1.21 or later
- Git

### Getting Started
1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/netprobe.git`
3. Create a feature branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Run tests: `make test`
6. Run linter: `make lint`
7. Commit your changes: `git commit -am 'Add some feature'`
8. Push to the branch: `git push origin feature/your-feature-name`
9. Submit a pull request

## Development Guidelines

### Code Style
- Follow standard Go formatting (`go fmt`)
- Use `golangci-lint` for code quality checks
- Write comprehensive tests for new features
- Add documentation for exported functions and types

### Commit Messages
- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove, etc.)
- Keep the first line under 50 characters
- Add detailed description if needed

### Testing
- Write unit tests for all new functionality
- Ensure all tests pass before submitting PR
- Aim for good test coverage
- Use table-driven tests where appropriate

### Pull Request Process
1. Ensure your code follows the guidelines above
2. Update documentation if needed
3. Add tests for new functionality
4. Ensure CI passes
5. Request review from maintainers

## Project Structure
```
netprobe/
├── cmd/           # Main applications
├── internal/      # Private application code
├── pkg/           # Public library code
├── docs/          # Documentation
├── scripts/       # Build and utility scripts
├── test/          # Additional test files
└── templates/     # Template files
```

## Security Considerations
- Be mindful of security implications in network-related code
- Follow secure coding practices
- Report security issues privately to maintainers

## License
By contributing to this project, you agree that your contributions will be licensed under the same license as the project (MIT License).
