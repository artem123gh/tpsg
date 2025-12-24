# TPSG Project Summary

## Project Overview

TPSG is a Go-based application that provides TCP and WebSocket server functionality with user authentication and configuration management. The project follows a modular structure with separate concerns for logging, type definitions, configuration management, global storage management, and network communication. The application uses TOML for application configuration and JSON for user credentials, providing a thread-safe global key-value storage system for runtime data.

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
│   ├── main.go                  # Application entry point with test mode support
│   ├── logging.go               # Logging functionality
│   ├── types.go                 # Type definitions
│   ├── gkvs.go                  # Global Key-Value Storage implementation
│   ├── config.go                # Configuration management and constants
│   ├── server_tcp.go            # TCP server implementation
│   ├── server_ws.go             # WebSocket server implementation
│   └── test_gkvs.go             # GKVS concurrent access test
├── build_debug.sh               # Build debug binary
├── build_release.sh             # Build release binary (optimized with -ldflags="-s -w")
├── run_console_debug.sh         # Run debug binary
├── run_console_release.sh       # Run release binary
└── run_test_gkvs.sh             # Run GKVS test
```

## External Configuration Files

The application uses external configuration files stored in `~/tpsg_configs/`:

- **config.toml** - Application settings (TCP and WS ports)
- **users.json** - User credentials (username/password pairs)

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

Represents user credentials with both username and password fields.

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

Helper constructors: `NewGKVSString(value)`, `NewGKVSInt32(value)`, `NewGKVSTConfigTOML(value)`, `NewGKVSTUserCreds(value)`, `NewGKVSNone()`, etc.

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

**Global Instances:**
- `TConfig *GKVS` - Global configuration storage, accessible throughout the application
- `TUsers *GKVS` - Global user credentials storage, accessible throughout the application

**Usage Example:**
```go
TConfig.Set("key", NewGKVSString("value"))
result := TConfig.Get("key").String  // Direct field access, no dereferencing
config := TConfig.Get("config").TConfigTOML  // Retrieve config object
userCreds := TUsers.Get("username1").TUserCreds  // Retrieve user credentials
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
- `TUsers *GKVS` - Global GKVS instance for user credentials storage

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

**ReadUsersConfig Function:**
```go
func ReadUsersConfig(usersConfigPath string) error
```

Reads and parses the JSON user credentials file:
- Takes full path to users.json file
- Reads file contents using `os.ReadFile`
- Parses JSON using `encoding/json` package
- Creates `TUserCreds` objects for each user
- Stores each user in the global `TUsers` GKVS instance with username as key
- Returns error if reading or parsing fails

**Users Configuration File Structure:**

The external `users.json` file (located in `~/tpsg_configs/users.json`) contains:
```json
{
    "username1": {
        "password": "password1"
    },
    "username2": {
        "password": "password2"
    },
    "username3": {
        "password": "password3"
    }
}
```

### 5. TCP Server (tpsg/server_tcp.go)

The TCP server provides network communication functionality with concurrent connection handling.

**RunTCPServer Function:**
```go
func RunTCPServer(port uint16)
```

Starts the TCP server in a background goroutine:
- Creates TCP listener on specified port
- Accepts incoming connections in a loop
- Spawns a new goroutine for each connection using `HandleTCPConnection`
- Logs server startup, connection acceptance, and errors
- Non-blocking - runs in background goroutine

**HandleTCPConnection Function:**
```go
func HandleTCPConnection(conn net.Conn)
```

Processes individual TCP connections:
- Runs in its own goroutine for each client
- Uses `bufio.Reader` to read requests line by line (newline-terminated)
- Processes requests synchronously within the connection
- Calls `ProcessTCPRequest` for each received request
- Sends responses back to the client
- Handles connection closure and errors gracefully
- Logs connection events and errors

**ProcessTCPRequest Function:**
```go
func ProcessTCPRequest(request string) string
```

Placeholder request processor (currently implements echo server):
- Receives request as string
- Logs the received request
- Returns echo response in format: `"Echo: <request>"`
- Will be replaced with actual protocol implementation later

**Connection Protocol:**
- Requests are newline-terminated text strings
- Each request is processed synchronously
- Response is sent immediately after processing
- Connection remains open for multiple requests
- Connection closes on client disconnect or error

### 6. WebSocket Server (tpsg/server_ws.go)

The WebSocket server provides bidirectional real-time communication functionality with concurrent connection handling and asynchronous request processing.

**RunWSServer Function:**
```go
func RunWSServer(port uint16)
```

Starts the WebSocket server in a background goroutine:
- Creates TCP listener on specified port (explicit `net.Listen`)
- Logs server startup after listener is successfully established
- Sets up HTTP upgrade handler for WebSocket connections
- Uses `http.Serve` with the established listener
- Spawns a new goroutine for each WebSocket connection using `HandleWSConnection`
- Non-blocking - runs in background goroutine

**HandleWSConnection Function:**
```go
func HandleWSConnection(conn *websocket.Conn)
```

Processes individual WebSocket connections:
- Runs in its own goroutine for each client
- Reads messages from the WebSocket connection in a loop
- Processes requests asynchronously - spawns a new goroutine for each incoming message
- Each goroutine calls `ProcessWSRequest` and sends the response back
- Handles connection closure and errors gracefully
- Logs connection events, message reception, and errors

**ProcessWSRequest Function:**
```go
func ProcessWSRequest(request string) string
```

Placeholder request processor (currently implements echo server):
- Receives request as string
- Logs the received request
- Returns echo response in format: `"Echo: <request>"`
- Will be replaced with actual protocol implementation later

**WebSocket Upgrader:**
- Uses `github.com/gorilla/websocket` library
- `CheckOrigin` returns true to accept connections from any origin
- Can be configured for production security requirements

**Connection Protocol:**
- Messages can be sent/received at any time (bidirectional)
- Each message is processed asynchronously in its own goroutine
- Response is sent after processing completes
- Connection remains open for multiple messages
- Connection closes on client disconnect or error

**Key Difference from TCP Server:**
- TCP: Requests processed synchronously within each connection
- WebSocket: Requests processed asynchronously (each in separate goroutine)

### 7. Main Application (tpsg/main.go)

The application entry point supports two modes: normal server mode and test mode.

**Test Mode:**
- Activated by passing `test-gkvs` as the first command-line argument
- Calls `TestGKVS()` function and exits without loading configuration or starting the server
- Allows testing specific features independently of the full application

**Normal Server Mode - Initialization Sequence:**
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

**User Credentials Loading:**
7. Calls `ReadUsersConfig(users_config_fullpath)` to read and parse users.json
8. Handles errors by logging with `LogError`
9. On success, each user is stored in TUsers GKVS with username as key
10. Logs success event with `LogEvent`

**Demonstration of GKVS Retrieval:**
11. Retrieves all stored paths from TConfig
12. Logs path values using `LogInfo`
13. Retrieves the config object from TConfig
14. Logs TCP and WS port values from the retrieved config

**Server Startup:**
15. Retrieves the config from TConfig
16. Calls `RunTCPServer(config_r.TCP)` with the TCP port from configuration
17. Calls `RunWSServer(config_r.WS)` with the WebSocket port from configuration
18. Uses `select {}` to keep the program running indefinitely

This workflow demonstrates:
- Reading external TOML configuration
- Reading external JSON user credentials
- Storing structured data in GKVS
- Retrieving and using stored configuration values
- Starting both TCP and WebSocket servers with configured ports
- Running the servers indefinitely

### 8. Testing System (tpsg/test_gkvs.go)

The project includes a testing framework for isolated feature testing without running the full server.

**TestGKVS Function:**
```go
func TestGKVS()
```

Demonstrates and tests concurrent access to the GKVS system:
- Spawns 3 goroutines that interact with `TConfig` at different times
- **Goroutine 1**: Immediately sets `test1 = 1` (uint16)
- **Goroutine 2**: Waits 5 seconds, retrieves and prints `test1`, then sets `test1 = 2`
- **Goroutine 3**: Waits 10 seconds, retrieves and prints `test1`
- Demonstrates thread-safe concurrent read/write operations
- Validates mutex-based synchronization prevents race conditions
- Uses direct field access pattern: `TConfig.Get("test1").UInt16`
- Runs for ~12 seconds to allow all goroutines to complete

**Test Execution:**
- Tests are invoked via command-line argument to `main()`
- Run with: `./run_test_gkvs.sh` or `go run . test-gkvs` (from tpsg/ directory)
- Test mode bypasses configuration loading and server startup
- Additional test functions can be added and invoked with different arguments

**Design Pattern:**
- All source files use `package main` (not an importable library)
- Test functions are part of the main package
- Command-line arguments switch between server mode and test mode
- Allows testing specific features in isolation without external dependencies

## Build System

- **Debug build**: Standard Go build without optimizations
- **Release build**: Optimized with stripped symbols (`-ldflags="-s -w"`)
- All build and run scripts are executable bash scripts in project root
- **Test script**: `run_test_gkvs.sh` - Runs GKVS concurrent access test

## Dependencies

External packages used:
- `github.com/BurntSushi/toml` v1.6.0 - TOML parsing
- `github.com/gorilla/websocket` v1.5.3 - WebSocket protocol implementation

## Key Design Principles

1. **Thread Safety**: GKVS uses RWMutex for safe concurrent access across goroutines
2. **Clean API**: Direct value types in GKVSTypes avoid pointer complexity
3. **Separation of Concerns**: Distinct files for logging, types, storage, configuration, servers, and testing functionality
4. **Standardized Logging**: Consistent timestamp format across all log functions - logs appear after critical resources are established
5. **Global Accessibility**: TConfig and TUsers are globally available for application-wide access
6. **External Configuration**: TOML-based config for settings, JSON-based config for user credentials
7. **Error Handling**: Functions return errors for caller to handle appropriately
8. **Concurrent Connection Handling**: Each TCP and WebSocket connection runs in its own goroutine
9. **Non-blocking Servers**: Both TCP and WebSocket servers run in background goroutines
10. **Asynchronous Request Processing**: WebSocket requests processed asynchronously (each in separate goroutine)
11. **Testability**: Command-line argument-based test mode for isolated feature testing without dependencies

## Current Status

The project has foundational infrastructure in place:
- ✅ Logging system
- ✅ Type definitions and tagged union system
- ✅ Thread-safe global key-value storage
- ✅ TOML configuration reading and parsing
- ✅ JSON user credentials reading and parsing
- ✅ Configuration storage in global GKVS
- ✅ User credentials storage in global GKVS
- ✅ Build and run scripts
- ✅ Path management for external config files
- ✅ TCP server with concurrent connection handling
- ✅ TCP request/response protocol (currently echo placeholder)
- ✅ WebSocket server with concurrent connection handling
- ✅ WebSocket asynchronous request processing
- ✅ WebSocket request/response protocol (currently echo placeholder)
- ✅ Testing framework with command-line test mode
- ✅ GKVS concurrent access test

The application is functional and can:
- Load configuration from external TOML file
- Load user credentials from external JSON file
- Start TCP server on the configured port
- Start WebSocket server on the configured port
- Accept multiple concurrent TCP connections
- Accept multiple concurrent WebSocket connections
- Process TCP requests synchronously and send responses (currently echo mode)
- Process WebSocket messages asynchronously and send responses (currently echo mode)
- Log all events and errors with timestamps
- Run indefinitely serving both TCP and WebSocket clients
- Run isolated feature tests via command-line arguments

**Next development steps** (as indicated in SPECS.md):
- Replace `ProcessTCPRequest` placeholder with actual protocol implementation
- Replace `ProcessWSRequest` placeholder with actual protocol implementation
- Define and implement request/response protocol specifications
