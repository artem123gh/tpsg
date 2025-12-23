# TPSG Project Summary

## Project Overview

TPSG is a Go-based application currently in early development. The project follows a modular structure with separate concerns for logging, type definitions, configuration management, and global storage management. The application uses TOML for external configuration and provides a thread-safe global key-value storage system for runtime data.

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
│   ├── go.sum                   # Go dependencies checksums
│   ├── main.go                  # Application entry point
│   ├── logging.go               # Logging functionality
│   ├── types.go                 # Type definitions
│   ├── gkvs.go                  # Global Key-Value Storage implementation
│   └── config.go                # Configuration management and constants
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

**TConfigTOML Type:**
```go
type TConfigTOML struct {
    TCP uint16 `toml:"TCP"`
    WS  uint16 `toml:"WS"`
}
```

Represents the structure of the external TOML configuration file. Contains server port settings:
- **TCP**: TCP listening port for server
- **WS**: WebSocket listening port for server

**GKVSTypes - Tagged Union for GKVS Values:**

A struct-based tagged union with a `Type` field indicating which value is active. All fields use direct values (not pointers) for clean API usage.

Supported types:
- Int8, UInt8, Int16, UInt16, Int32, UInt32, Int64, UInt64
- Float32, Float64
- String
- TUserCreds
- TConfigTOML
- None (empty value)

Helper constructors: `NewGKVSString(value)`, `NewGKVSInt32(value)`, `NewGKVSTConfigTOML(value)`, `NewGKVSNone()`, etc.

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
config := TConfig.Get("config").TConfigTOML  // Retrieve config object
```

### 4. Configuration Management (tpsg/config.go)

**Hard-coded Constants:**
```go
const CONFIGS_FOLDER = "tpsg_configs"
const CONFIG_FILE = "config.toml"
const USERS_CONFIG_FILE = "users.json"
```

**Global Storage:**
- `TConfig *GKVS` - Global GKVS instance for application-wide configuration

**ReadConfig Function:**
```go
func ReadConfig(configPath string) (TConfigTOML, error)
```

Reads and parses the TOML configuration file:
- Takes full path to config.toml file
- Reads file contents using `os.ReadFile`
- Parses TOML using `github.com/BurntSushi/toml` package
- Returns `TConfigTOML` struct and error
- Does NOT store in TConfig - caller is responsible for storage

**Configuration File Structure:**

The external `config.toml` file (located in `~/tpsg_configs/config.toml`) contains:
```toml
TCP = 8080
WS = 8081
```

### 5. Main Application (tpsg/main.go)

The application demonstrates the complete configuration and storage workflow:

**Initialization Sequence:**
1. Constructs configuration paths from HOME environment variable
2. Stores all paths in global `TConfig` GKVS instance:
   - `user_folder` - User's home directory
   - `configs_folder_path` - Path to config folder
   - `config_fullpath` - Full path to config.toml
   - `users_config_fullpath` - Full path to users.json

**TOML Configuration Loading:**
3. Calls `ReadConfig(config_fullpath)` to read and parse config.toml
4. Handles errors by logging with `LogError`
5. On success, stores parsed config in TConfig under key "config"
6. Logs success event with `LogEvent`

**Demonstration of GKVS Retrieval:**
7. Retrieves all stored paths from TConfig
8. Logs path values using `LogInfo`
9. Retrieves the config object from TConfig
10. Logs TCP and WS port values from the retrieved config

This workflow demonstrates:
- Reading external TOML configuration
- Storing structured data in GKVS
- Retrieving and using stored configuration values

## Build System

- **Debug build**: Standard Go build without optimizations
- **Release build**: Optimized with stripped symbols (`-ldflags="-s -w"`)
- All build scripts are executable bash scripts in project root

## Dependencies

External packages used:
- `github.com/BurntSushi/toml` v1.6.0 - TOML parsing

## Key Design Principles

1. **Thread Safety**: GKVS uses RWMutex for safe concurrent access across goroutines
2. **Clean API**: Direct value types in GKVSTypes avoid pointer complexity
3. **Separation of Concerns**: Distinct files for logging, types, storage, and configuration
4. **Standardized Logging**: Consistent timestamp format across all log functions
5. **Global Accessibility**: TConfig is globally available for application-wide configuration
6. **External Configuration**: TOML-based config files for runtime settings
7. **Error Handling**: Functions return errors for caller to handle appropriately

## Current Status

The project has foundational infrastructure in place:
- ✅ Logging system
- ✅ Type definitions and tagged union system
- ✅ Thread-safe global key-value storage
- ✅ TOML configuration reading and parsing
- ✅ Configuration storage in global GKVS
- ✅ Build and run scripts
- ✅ Path management for external config files

The application currently demonstrates the complete configuration workflow: constructing paths, reading external TOML config, storing configuration in the global GKVS, and retrieving values for use throughout the application.
