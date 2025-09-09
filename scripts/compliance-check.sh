#!/bin/bash
set -euo pipefail

# NetProbe Compliance Check Script
# Performs basic security and compliance checks

echo "🔍 Running NetProbe Compliance Checks..."
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    local status=$1
    local message=$2
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $message"
    elif [ "$status" = "WARN" ]; then
        echo -e "${YELLOW}⚠${NC} $message"
    else
        echo -e "${RED}✗${NC} $message"
    fi
}

# Check if running as root (security concern)
if [ "$EUID" = 0 ]; then
    print_status "WARN" "Running as root - consider running as non-privileged user"
else
    print_status "PASS" "Running as non-privileged user"
fi

# Check for required files
echo ""
echo "📁 Checking required files..."

files=("go.mod" "go.sum" "README.md" "LICENSE")
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        print_status "PASS" "Found $file"
    else
        print_status "FAIL" "Missing $file"
    fi
done

# Check Go version
echo ""
echo "🐹 Checking Go environment..."
if command -v go &> /dev/null; then
    go_version=$(go version | grep -oE 'go[0-9]+\.[0-9]+')
    print_status "PASS" "Go installed: $go_version"

    # Check if Go version is 1.21+
    go_major=$(echo "$go_version" | sed 's/go//' | cut -d. -f1)
    go_minor=$(echo "$go_version" | sed 's/go//' | cut -d. -f2)

    if [ "$go_major" -gt 1 ] || ([ "$go_major" -eq 1 ] && [ "$go_minor" -ge 21 ]); then
        print_status "PASS" "Go version meets minimum requirement (1.21+)"
    else
        print_status "FAIL" "Go version too old (requires 1.21+)"
    fi
else
    print_status "FAIL" "Go not installed"
fi

# Check for security tools
echo ""
echo "🔒 Checking security tools..."
security_tools=("golangci-lint" "gosec")
for tool in "${security_tools[@]}"; do
    if command -v "$tool" &> /dev/null; then
        print_status "PASS" "$tool available"
    else
        print_status "WARN" "$tool not found (recommended for security checks)"
    fi
done

# Check for sensitive files that shouldn't be committed
echo ""
echo "🚫 Checking for sensitive files..."
sensitive_files=(".env" "config.yaml" "*.key" "*.pem" "*.crt")
found_sensitive=false
for pattern in "${sensitive_files[@]}"; do
    if compgen -G "$pattern" > /dev/null; then
        print_status "WARN" "Found sensitive file pattern: $pattern"
        found_sensitive=true
    fi
done
if [ "$found_sensitive" = false ]; then
    print_status "PASS" "No sensitive files found in repository"
fi

# Check .gitignore
echo ""
echo "📋 Checking .gitignore..."
if [ -f ".gitignore" ]; then
    print_status "PASS" ".gitignore exists"

    # Check for important patterns in .gitignore
    required_patterns=("*.exe" "*.dll" ".env" "config.yaml")
    for pattern in "${required_patterns[@]}"; do
        if grep -q "$pattern" .gitignore; then
            print_status "PASS" ".gitignore contains: $pattern"
        else
            print_status "WARN" ".gitignore missing: $pattern"
        fi
    done
else
    print_status "FAIL" ".gitignore not found"
fi

# Check for hardcoded secrets (basic pattern matching)
echo ""
echo "🔍 Checking for potential hardcoded secrets..."
secret_patterns=("password.*=" "secret.*=" "key.*=" "token.*=")
found_secrets=false

# Search Go files for potential secrets
while IFS= read -r -d '' file; do
    for pattern in "${secret_patterns[@]}"; do
        if grep -qi "$pattern" "$file"; then
            print_status "WARN" "Potential secret found in: $file"
            found_secrets=true
        fi
    done
done < <(find . -name "*.go" -type f -print0)

if [ "$found_secrets" = false ]; then
    print_status "PASS" "No obvious hardcoded secrets detected"
fi

# Final summary
echo ""
echo "========================================"
echo "🏁 Compliance check completed!"
echo ""
echo "Recommendations:"
echo "- Ensure all security tools are installed for comprehensive checks"
echo "- Review any WARN or FAIL items above"
echo "- Consider using gosec for static security analysis"
echo "- Keep dependencies updated with 'go mod tidy'"
echo ""
print_status "PASS" "Basic compliance checks completed"
