# TPSG Project Summary

## Project Overview

TPSG is a Go-based application currently in early development. The project follows a modular structure with separate concerns for logging, type definitions, configuration, and global storage management.

## Project Structure

```
tpsg/
├── ai/                          # AI assistant documentation
│   ├── RULES.md                 # Coding rules (4-space indentation, no commits by AI)
│   ├── SPECS.md                 # Project specifications
│   ├── PROMPTS.md               # AI prompts (not for AI to read)
│   └── SUMMARY.md               # This file - project summary
├── bins/                        # Compiled binaries (tpsg_debug, tpsg_release)
├── other/                       # Temporary reference files (not in repo)
├── tpsg/                        # Go module source code
│   ├── go.mod                   # Go module file (module: tpsg, Go 1.24.1)
│   ├── main.go                  # Application entry point
│   ├── logging.go               # Logging functionality
│   ├── types.go                 # Type definitions
│   ├── gkvs.go                  # Global Key-Value Storage implementation
│   └── config.go                # Configuration constants
├── build_debug.sh               # Build debug binary
├── build_release.sh             # Build release binary (optimized with -ldflags="-s -w")
├── run_console_debug.sh         # Run debug binary
└── run_console_release.sh       # Run release binary
```

## Implemented Components

### 1. Logging System (tpsg/logging.go)

Three logging functions with timestamp formatting `"Type | YYYY.MM.DD HH:mm:ss.milliseconds | message"`:

- **LogInfo(message string)**: For user debug prints - AI should NOT use unless explicitly instructed
- **LogEvent(message string)**: For logging application events - AI should use for event logging
- **LogError(message string)**: For logging errors

Usage: `LogEvent(fmt.Sprintf("formatted message: %s", value))`

### 2. Type System (tpsg/types.go)

**Status Type:**
```go
type Status bool
const (
    Success Status = true
    Failure Status = false
)
```

**TUserCreds Type:**
```go
type TUserCreds struct {
    Username string
    Password string
}
```

**GKVSTypes - Tagged Union for GKVS Values:**

A struct-based tagged union with a `Type` field indicating which value is active. All fields use direct values (not pointers) for clean API usage.

Supported types:
- Int8, UInt8, Int16, UInt16, Int32, UInt32, Int64, UInt64
- Float32, Float64
- String
- TUserCreds
- None (empty value)

Helper constructors: `NewGKVSString(value)`, `NewGKVSInt32(value)`, `NewGKVSNone()`, etc.

**Key Design Decision:** Uses direct values instead of pointers to avoid `&` and `*` operations in user code.

### 3. Global Key-Value Storage - GKVS (tpsg/gkvs.go)

Thread-safe key-value storage using `sync.RWMutex` for concurrent access.

**Structure:**
```go
type GKVS struct {
    storage map[string]GKVSTypes
    mutex   sync.RWMutex
}
```

**API Methods:**
- `Set(key string, value GKVSTypes) GKVSTypes` - Create/update key-value, returns value
- `Get(key string) GKVSTypes` - Retrieve value, returns None if not found
- `Delete(key string) GKVSTypes` - Remove and return value, returns None if not found

**Thread Safety:**
- Uses `RLock()` for `Get()` - allows concurrent reads
- Uses `Lock()` for `Set()` and `Delete()` - exclusive write access
- Each instance has its own mutex for independent operation
- Deadlock prevention: single lock acquisition per method, deferred unlocks

**Global Instance:**
- `TConfig *GKVS` - Global configuration storage, accessible throughout the application

**Usage Example:**
```go
TConfig.Set("key", NewGKVSString("value"))
result := TConfig.Get("key").String  // Direct field access, no dereferencing
```

### 4. Configuration (tpsg/config.go)

Hard-coded constants:
```go
const CONFIGS_FOLDER = "tpsg_configs"
const CONFIG_FILE = "config.toml"
const USERS_CONFIG_FILE = "users.json"
```

### 5. Main Application (tpsg/main.go)

Currently demonstrates GKVS usage by:
1. Constructing configuration paths from HOME environment variable
2. Storing paths in global `TConfig` GKVS instance
3. Retrieving and logging the paths using `LogInfo`

## Build System

- **Debug build**: Standard Go build without optimizations
- **Release build**: Optimized with stripped symbols (`-ldflags="-s -w"`)
- All build scripts are executable bash scripts in project root

## Key Design Principles

1. **Thread Safety**: GKVS uses RWMutex for safe concurrent access across goroutines
2. **Clean API**: Direct value types in GKVSTypes avoid pointer complexity
3. **Separation of Concerns**: Distinct files for logging, types, storage, and configuration
4. **Standardized Logging**: Consistent timestamp format across all log functions
5. **Global Accessibility**: TConfig is globally available for application-wide configuration

## Current Status

The project has foundational infrastructure in place:
- ✅ Logging system
- ✅ Type definitions and tagged union system
- ✅ Thread-safe global key-value storage
- ✅ Build and run scripts
- ✅ Basic configuration constants

The application currently serves as a proof-of-concept for the GKVS system, with the main function demonstrating storage and retrieval of configuration paths.
