# Contributing to SochDB Go SDK

Thank you for your interest in contributing to the SochDB Go SDK! This guide provides all the information you need to build, test, and contribute to the project.

---

## Table of Contents

- [Development Setup](#development-setup)
- [Building from Source](#building-from-source)
- [Running Tests](#running-tests)
- [Server Setup for Development](#server-setup-for-development)
- [Code Style](#code-style)
- [Pull Request Process](#pull-request-process)
- [Architecture Overview](#architecture-overview)
- [Migration Guide](#migration-guide)

---

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Rust toolchain (for building server)
- Protocol Buffers compiler (protoc)
- Git

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/sochdb/sochdb-go.git
cd sochdb-go

# Download dependencies
go mod download

# Build
go build ./...

# Run tests
go test ./...
```

---

## Building from Source

### Go SDK Only

```bash
cd sochdb-go
go mod download
go build ./...
```

### With Protocol Buffers

If you need to regenerate gRPC stubs:

```bash
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate from proto files
cd sochdb/proto
protoc --go_out=. --go-grpc_out=. *.proto
```

---

## Running Tests

### Unit Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with verbose output
go test -v ./...
```

### Integration Tests

```bash
# Start SochDB server first
cd sochdb
cargo run -p sochdb-grpc

# In another terminal, run integration tests
cd sochdb-go
go test -tags=integration ./...
```

### Run Examples

```bash
# Test example
cd example
go run main.go
```

---

## Server Setup for Development

### Starting the Server

```bash
# Development mode
cd sochdb
cargo run -p sochdb-grpc

# Production mode (optimized)
cargo build --release -p sochdb-grpc
./target/release/sochdb-grpc --host 0.0.0.0 --port 50051
```

### Server Configuration

The server runs all business logic including:
- ✅ HNSW vector indexing (15x faster than ChromaDB)
- ✅ SQL query parsing and execution
- ✅ Graph traversal algorithms
- ✅ Policy evaluation
- ✅ Multi-tenant namespace isolation
- ✅ Collection management

### Configuration File

Create `sochdb-server-config.toml`:

```toml
[server]
host = "0.0.0.0"
port = 50051

[storage]
data_dir = "./data"

[logging]
level = "info"
```

---

## Code Style

### Go

We follow standard Go conventions:

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run

# Vet
go vet ./...
```

### Commit Messages

Follow conventional commits:

```
feat: Add temporal graph support
fix: Handle connection timeout
docs: Update API reference
test: Add integration tests for graphs
```

### Code Review Checklist

- [ ] All tests pass
- [ ] Code follows Go style guidelines
- [ ] Documentation updated (godoc)
- [ ] Examples added/updated if needed
- [ ] No breaking changes (or documented in CHANGELOG)

---

## Pull Request Process

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/sochdb-go.git
   cd sochdb-go
   ```

2. **Create Feature Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Changes**
   - Write code
   - Add tests
   - Update documentation

4. **Test Locally**
   ```bash
   go test ./...
   go fmt ./...
   golangci-lint run
   ```

5. **Commit and Push**
   ```bash
   git add .
   git commit -m "feat: Your feature description"
   git push origin feature/your-feature-name
   ```

6. **Create Pull Request**
   - Go to GitHub
   - Create PR from your branch
   - Fill out PR template
   - Wait for review

---

## Architecture Overview

### Thin Client Architecture

```
┌────────────────────────────────────────────────┐
│         Rust Server (sochdb-grpc)              │
├────────────────────────────────────────────────┤
│  • All business logic (Graph, Policy, Search)  │
│  • Vector operations (HNSW)                    │
│  • SQL parsing & execution                     │
│  • Collections & Namespaces                    │
│  • Single source of truth                      │
└────────────────────────────────────────────────┘
                       │ gRPC/IPC
                       ▼
            ┌─────────────────────┐
            │     Go SDK          │
            │   (~1,144 LOC)      │
            ├─────────────────────┤
            │ • Transport layer   │
            │ • Type definitions  │
            │ • Zero logic        │
            └─────────────────────┘
```

### Key Components

**client.go**
- Base client interface
- Connection management

**grpc_client.go**
- gRPC client implementation
- All server operations
- Error handling

**errors.go**
- SochDBError type
- Error constructors
- Error messages

**format.go**
- WireFormat enum
- ContextFormat enum
- FormatCapabilities utilities

**sochdb.go**
- Package documentation
- Version constants

### Comparison with Old Architecture

| Feature | Old (Fat Client) | New (Thin Client) |
|---------|------------------|-------------------|
| SDK Size | 4,302 LOC | 1,144 LOC (-73%) |
| Business Logic | In SDK (Go) | In Server (Rust) |
| Bug Fixes | Per language | Once in server |
| Semantic Drift | High risk | Zero risk |
| Performance | FFI overhead | Network call |
| Maintenance | 3x effort | 1x effort |

---

## Migration Guide

### From v0.3.3 to v0.3.4

**Key Changes:**
- Removed embedded `Database` type
- All operations now go through `GrpcClient`
- Server must be running for all operations
- FFI bindings removed

**Old Code:**
```go
import "github.com/sochdb/sochdb-go"

db := sochdb.Open("./data")
defer db.Close()

tx := db.Begin()
tx.Put([]byte("key"), []byte("value"))
tx.Commit()
```

**New Code:**
```go
import "github.com/sochdb/sochdb-go"

// Start server first: cargo run -p sochdb-grpc
client := sochdb.NewGrpcClient("localhost:50051")
defer client.Close()

err := client.PutKv("key", []byte("value"))
if err != nil {
    log.Fatal(err)
}
```

**Migration Checklist:**
- [ ] Start SochDB server (cargo run -p sochdb-grpc)
- [ ] Replace `sochdb.Open()` with `sochdb.NewGrpcClient()`
- [ ] Remove transaction Begin/Commit (server manages)
- [ ] Add error handling for all operations
- [ ] Update connection strings to point to server

---

## Release Process

### Version Bumping

```bash
# Update version in sochdb.go
vim sochdb.go

# Update go.mod if needed
vim go.mod

# Update CHANGELOG.md
vim CHANGELOG.md
```

### Tagging

```bash
# Create tag
git tag v0.3.4

# Push tag
git push origin v0.3.4
```

### Publishing

Go modules are automatically versioned via Git tags. Users import via:

```go
import "github.com/sochdb/sochdb-go@v0.3.4"
```

---

## Testing Checklist

Before submitting a PR, ensure:

- [ ] All unit tests pass: `go test ./...`
- [ ] Integration tests pass (with server): `go test -tags=integration ./...`
- [ ] Example runs: `cd example && go run main.go`
- [ ] Code formatted: `go fmt ./...`
- [ ] Linting passes: `golangci-lint run`
- [ ] Vet passes: `go vet ./...`
- [ ] Documentation updated (godoc comments)
- [ ] CHANGELOG.md updated

---

## Performance Testing

### Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Run with CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./...

# Analyze profile
go tool pprof cpu.prof
```

### Load Testing

```bash
# Start server
cd sochdb
cargo run -p sochdb-grpc --release

# Run load test
cd sochdb-go/tests
go test -bench=BenchmarkBatchInsert -benchtime=10s
```

---

## Getting Help

- **Main Repo**: https://github.com/sochdb/sochdb
- **Go SDK Issues**: https://github.com/sochdb/sochdb-go/issues
- **Discussions**: https://github.com/sochdb/sochdb/discussions
- **Contributing Guide**: See main repo [CONTRIBUTING.md](https://github.com/sochdb/sochdb/blob/main/CONTRIBUTING.md)

---

## License

By contributing to SochDB Go SDK, you agree that your contributions will be licensed under the Apache License 2.0.
