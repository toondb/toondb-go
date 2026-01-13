# Quick Start Guide - SochDB Go SDK

## Installation (3 Steps)

### 1. Install pkg-config (if not installed)

**macOS:**
```bash
brew install pkg-config
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install pkg-config

# Fedora/RHEL
sudo yum install pkgconfig
```

### 2. Install the Native Library

```bash
# Clone and install
git clone https://github.com/sochdb/sochdb-go.git
cd sochdb-go

# Run installer (requires sudo for system installation)
SOCHDB_ROOT=/path/to/sochdb ./install-lib.sh
```

The installer will:
- Copy `libsochdb_storage.{dylib,so}` to `/usr/local/lib`
- Create pkg-config file for automatic detection
- Update library cache (Linux)

### 3. Install Go SDK

```bash
go get github.com/sochdb/sochdb-go@v0.4.0
```

## Verify Installation

```bash
pkg-config --libs libsochdb_storage
# Expected: -L/usr/local/lib -lsochdb_storage
```

## Your First Program

Create `main.go`:

```go
package main

import (
    "fmt"
    "log"
    "github.com/sochdb/sochdb-go/embedded"
)

func main() {
    // Open database in embedded mode
    db, err := embedded.Open("./mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Write data
    err = db.Put([]byte("greeting"), []byte("Hello, SochDB!"))
    if err != nil {
        log.Fatal(err)
    }

    // Read data
    value, err := db.Get([]byte("greeting"))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(value)) // Output: Hello, SochDB!
}
```

Run it:

```bash
go run main.go
```

That's it! No environment variables, no manual CGO_LDFLAGS - it just works! 🎉

## What Makes This Easy?

- **pkg-config integration**: Automatically finds the native library
- **One-time setup**: Install once, use everywhere
- **Standard Go workflow**: Works like any other Go package
- **No environment variables**: No need to set DYLD_LIBRARY_PATH or CGO_LDFLAGS

## Troubleshooting

### Issue: "pkg-config: executable file not found"

**Solution:** Install pkg-config (see step 1 above)

### Issue: "package 'libsochdb_storage' not found"

**Solution:** Run the install script:
```bash
SOCHDB_ROOT=/path/to/sochdb ./install-lib.sh
```

### Issue: Library not found at runtime

**Solution:** Verify installation:
```bash
ls -l /usr/local/lib/libsochdb_storage.*
```

## Development Mode (No Installation)

If you're developing and don't want to install system-wide:

```bash
# Set environment variables
export DYLD_LIBRARY_PATH=/path/to/sochdb/target/debug  # macOS
export LD_LIBRARY_PATH=/path/to/sochdb/target/debug    # Linux
export CGO_LDFLAGS="-L/path/to/sochdb/target/debug -lsochdb_storage"

# Run your code
go run main.go
```

## Next Steps

- See [example/main.go](example/main.go) for comprehensive examples
- Read the [full documentation](README.md) for advanced features
- Check [examples/embedded/main.go](examples/embedded/main.go) for FFI usage patterns
