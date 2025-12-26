# tpserde

Serialization and deserialization library for "theplatform" binary format.

## Overview

This package provides functionality to serialize and deserialize data between Go and "theplatform" (a Rust-based programming language) over TCP and WebSocket connections. It implements the theplatform binary format with optional LZ4 compression.

## Features

- **Complete Type System**: Supports all theplatform data types including scalars, vectors, and complex types (lists, dicts, tables, lambdas, patterns)
- **Special Numeric Values**: Full support for NULL values (`0N`) and infinities (`W`, `-W`) for all numeric types
- **Binary Serialization**: Converts Go types to theplatform binary format
- **Binary Deserialization**: Converts theplatform binary format to Go types
- **LZ4 Compression**: Automatic compression/decompression for large data (> 4096 bytes)
- **Full Compatibility**: Matches theplatform's binary format exactly (little-endian, same structure)

## Main API

### Types

- `TPTypes` - Main data type that holds theplatform values
- `TPBinary` - Binary data in theplatform format ([]byte)

### Functions

#### TPDataSer

Serializes TPTypes to theplatform binary format.

```go
func TPDataSer(data TPTypes, compress bool) (TPBinary, error)
```

Parameters:
- `data`: The TPTypes value to serialize
- `compress`: If true, compresses data with LZ4 when size exceeds 4096 bytes

Returns:
- `TPBinary`: The serialized binary data
- `error`: Error if serialization fails

#### TPDataDe

Deserializes theplatform binary data to TPTypes. Automatically detects and decompresses LZ4 data.

```go
func TPDataDe(data TPBinary) (TPTypes, error)
```

Parameters:
- `data`: Binary data in theplatform format

Returns:
- `TPTypes`: The deserialized value
- `error`: Error if deserialization fails

#### Handshake Functions

For IPC protocol compatibility with theplatform, perform a handshake before data exchange.

**ExchangeHandshake** - Client side (active side) handshake:

```go
func ExchangeHandshake(stream io.ReadWriter, features Features) (Handshake, error)
```

Parameters:
- `stream`: The network connection (e.g., net.Conn)
- `features`: Requested features (e.g., `NewFeatures().WithBuffered()`)

Returns:
- `Handshake`: The server's handshake response
- `error`: Error if handshake fails

**ResponseHandshake** - Server side (passive side) handshake:

```go
func ResponseHandshake(stream io.ReadWriter, features Features) (Handshake, error)
```

Parameters:
- `stream`: The network connection (e.g., net.Conn)
- `features`: Supported features (e.g., `NewFeatures().WithBuffered()`)

Returns:
- `Handshake`: The client's handshake request
- `error`: Error if handshake fails

**Handshake Structure:**
- `Version`: IPC protocol version (currently 0.1.0)
- `Features`: Feature flags (compressed, buffered, etc.)

## Usage Examples

### Scalar Types

```go
import "tpsg/tpserde"

// Create scalar values
intVal := tpserde.NewTPInt(42)
floatVal := tpserde.NewTPFloat(3.14)
strVal := tpserde.NewTPVecChar("Hello, World!")

// Serialize without compression
binary, err := tpserde.TPDataSer(intVal, false)
if err != nil {
    // handle error
}

// Deserialize
result, err := tpserde.TPDataDe(binary)
if err != nil {
    // handle error
}
```

### Vector Types

```go
// Create vector values
vecInt := tpserde.NewTPVecInt([]int32{1, 2, 3, 4, 5})
vecFloat := tpserde.NewTPVecFloat([]float64{1.1, 2.2, 3.3})
vecStr := tpserde.NewTPVecSymbol([]string{"foo", "bar", "baz"})

// Serialize with compression (if large enough)
binary, err := tpserde.TPDataSer(vecInt, true)
```

### Working with Special Values

```go
// Create vector with NULL and infinity values
vecLong := tpserde.NewTPVecLong([]int64{
    100,
    tpserde.NULL_LONG,      // NULL value (0N)
    tpserde.INF_LONG,        // Positive infinity (W)
    tpserde.NEG_INF_LONG,    // Negative infinity (-W)
    200,
})

// Or use constructor functions
data := tpserde.NewTPList([]tpserde.TPTypes{
    tpserde.NewTPInt(42),
    tpserde.NewTPNullLong(),    // NULL long
    tpserde.NewTPInfFloat(),    // Positive infinity float
})

// Serialize and deserialize
binary, _ := tpserde.TPDataSer(vecLong, false)
result, _ := tpserde.TPDataDe(binary)

// Check for special values
resultVec := result.Data.(tpserde.TPVecLong)
for i, v := range resultVec {
    if tpserde.IsNullLong(int64(v)) {
        fmt.Printf("Element %d is NULL\n", i)
    } else if tpserde.IsPosInfLong(int64(v)) {
        fmt.Printf("Element %d is positive infinity\n", i)
    } else if tpserde.IsNegInfLong(int64(v)) {
        fmt.Printf("Element %d is negative infinity\n", i)
    } else {
        fmt.Printf("Element %d = %d\n", i, v)
    }
}
```

### Complex Types

```go
// Create a list
list := tpserde.NewTPList([]tpserde.TPTypes{
    tpserde.NewTPInt(42),
    tpserde.NewTPVecChar("hello"),
    tpserde.NewTPFloat(3.14),
})

// Create a dict
keys := tpserde.NewTPVecSymbol([]string{"a", "b", "c"})
values := tpserde.NewTPVecInt([]int32{1, 2, 3})
dict := tpserde.NewTPDict(keys, values)

// Create a table
tableKeys := tpserde.NewTPVecSymbol([]string{"col1", "col2"})
tableValues := tpserde.NewTPList([]tpserde.TPTypes{
    tpserde.NewTPVecInt([]int32{1, 2, 3}),
    tpserde.NewTPVecFloat([]float64{1.1, 2.2, 3.3}),
})
table := tpserde.NewTPTable(tableKeys, tableValues)
```

### IPC Protocol with Handshake

When communicating with theplatform using the IPC protocol, you must perform a handshake before exchanging data:

```go
import (
    "net"
    "tpsg/tpserde"
)

// Client side (active side)
func clientExample() {
    // Connect to server
    conn, err := net.Dial("tcp", "127.0.0.1:8080")
    if err != nil {
        // handle error
    }
    defer conn.Close()

    // Perform handshake as active side (client)
    features := tpserde.NewFeatures().WithBuffered()
    handshake, err := tpserde.ExchangeHandshake(conn, features)
    if err != nil {
        // handle error
    }

    // Check server features
    if handshake.Features.IsBuffered() {
        // Server supports buffered mode
    }

    // Now exchange data
    data := tpserde.NewTPVecChar("Hello, theplatform!")
    binary, _ := tpserde.TPDataSer(data, true)
    conn.Write(binary)
}

// Server side (passive side)
func serverExample(conn net.Conn) {
    defer conn.Close()

    // Perform handshake as passive side (server)
    features := tpserde.NewFeatures().WithBuffered()
    handshake, err := tpserde.ResponseHandshake(conn, features)
    if err != nil {
        // handle error
    }

    // Check client features
    if handshake.Features.IsBuffered() {
        // Client supports buffered mode
    }

    // Now receive and process data
    // (implement your data reading logic here)
}
```

### With Network I/O (Simple Example)

```go
import (
    "net"
    "tpsg/tpserde"
)

// Sending data over TCP
func sendData(conn net.Conn, data tpserde.TPTypes) error {
    // Serialize with compression
    binary, err := tpserde.TPDataSer(data, true)
    if err != nil {
        return err
    }

    // Send over connection
    _, err = conn.Write(binary)
    return err
}

// Receiving data from TCP
func receiveData(conn net.Conn) (tpserde.TPTypes, error) {
    // Read binary data from connection
    // (You'll need to implement framing/buffering based on your protocol)
    buffer := make([]byte, 65536)
    n, err := conn.Read(buffer)
    if err != nil {
        return tpserde.TPTypes{}, err
    }

    // Deserialize
    return tpserde.TPDataDe(tpserde.TPBinary(buffer[:n]))
}
```

## Supported Types

### Scalar Types
- NIL, ANY
- SC_BOOL (bool)
- SC_BYTE (int8), SC_SHORT (int16), SC_INT (int32), SC_LONG (int64)
- SC_REAL (float32), SC_FLOAT (float64)
- SC_ENUM (uint32)
- SC_GUID (UUID)
- SC_SYMBOL (string)

### Temporal Types
- SC_MONTH, SC_DATE, SC_MINUTE, SC_SECOND, SC_TIME (int32)
- SC_TIMESTAMP, SC_DATETIME, SC_TIMESPAN (int64)

### Special Numeric Values

Theplatform defines special values for numeric types, which are fully supported:

#### NULL Values (equivalent to theplatform's `0N`)
- `NULL_BYTE`, `NULL_SHORT`, `NULL_INT`, `NULL_LONG` (minimum values)
- `NULL_REAL`, `NULL_FLOAT` (negative NaN)
- `NULL_ENUM`

#### Infinity Values (theplatform's `W` and `-W`)
- Positive infinity: `INF_BYTE`, `INF_SHORT`, `INF_INT`, `INF_LONG`, `INF_REAL`, `INF_FLOAT`
- Negative infinity: `NEG_INF_BYTE`, `NEG_INF_SHORT`, `NEG_INF_INT`, `NEG_INF_LONG`, `NEG_INF_REAL`, `NEG_INF_FLOAT`

Helper functions are provided to detect and create these special values:
```go
// Detection
IsNullLong(value)      // Check if value is NULL
IsPosInfLong(value)    // Check if value is positive infinity
IsNegInfLong(value)    // Check if value is negative infinity
IsInfLong(value)       // Check if value is any infinity

// Creation
NewTPNullLong()        // Create NULL long (0N)
NewTPInfLong()         // Create positive infinity (W)
NewTPNegInfLong()      // Create negative infinity (-W)
```

Similar functions exist for all numeric types (Byte, Short, Int, Real, Float)

### Vector Types
- VEC_BOOL ([]bool)
- VEC_BYTE ([]int8), VEC_SHORT ([]int16), VEC_INT ([]int32), VEC_LONG ([]int64)
- VEC_REAL ([]float32), VEC_FLOAT ([]float64)
- VEC_ENUM ([]uint32)
- VEC_GUID ([]UUID)
- VEC_SYMBOL ([]string)
- VEC_CHAR (string)
- Temporal vectors: VEC_MONTH, VEC_DATE, VEC_TIMESTAMP, etc.

### Complex Types
- LIST ([]TPTypes)
- DICT (keys and values TPTypes)
- TABLE (keys and values TPTypes)
- PATTERN (exprs and arms TPTypes)
- LAMBDA (function with metadata)

## Binary Format

The binary format consists of:

1. **Header** (16 bytes on 64-bit systems):
   - Features (4 bytes): Bit flags for compression, buffering, etc.
   - Reserved (4 bytes): Reserved for future use
   - Length (8 bytes): Length of payload data

2. **Payload**:
   - Type tag (4 bytes): Identifies the data type
   - Data: Type-specific serialized data
   - All multi-byte values are little-endian

## Compression

- Compression is optional via the `compress` parameter in `TPDataSer`
- Uses LZ4 compression algorithm
- Only compresses when data size > 4096 bytes (DEFAULT_UNCOMPRESSED_SIZE_LIMIT)
- Decompression is automatic in `TPDataDe` (detects compression flag in header)

## Testing

Run tests with:

```bash
cd tpsg
go test ./tpserde -v
```

All tests include:
- Scalar type serialization/deserialization
- Vector type serialization/deserialization
- Complex type serialization/deserialization
- Compression/decompression
- Round-trip tests
- Header serialization

## Dependencies

- `github.com/google/uuid` - UUID support
- `github.com/pierrec/lz4/v4` - LZ4 compression
