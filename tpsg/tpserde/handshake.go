package tpserde

import (
    "encoding/binary"
    "fmt"
    "io"
)

// IPC protocol version (0.1.0)
// Encoded as: major << 20 | minor << 10 | patch
const IPC_VERSION uint32 = 0<<20 | 1<<10 | 0

// Handshake represents the IPC protocol handshake
// Layout matches Rust #[repr(C)] struct Handshake:
//   version: u32
//   features: u32
type Handshake struct {
    Version  uint32
    Features Features
}

// NewHandshake creates a new handshake with default IPC version and no features
func NewHandshake() Handshake {
    return Handshake{
        Version:  IPC_VERSION,
        Features: NewFeatures(),
    }
}

// WithCompressed adds the compressed feature flag
func (h Handshake) WithCompressed() Handshake {
    h.Features = h.Features.WithCompressed()
    return h
}

// WithBuffered adds the buffered feature flag
func (h Handshake) WithBuffered() Handshake {
    h.Features = h.Features.WithBuffered()
    return h
}

// WithUnsupported adds the unsupported feature flag
func (h Handshake) WithUnsupported() Handshake {
    h.Features = h.Features.WithUnsupported()
    return h
}

// IsUnsupported checks if the handshake has unsupported features
func (h Handshake) IsUnsupported() bool {
    return h.Features.IsUnsupported()
}

// Serialize writes the handshake to the writer in little-endian format
func (h Handshake) Serialize(writer io.Writer) error {
    // Write version (u32)
    if err := binary.Write(writer, binary.LittleEndian, h.Version); err != nil {
        return fmt.Errorf("failed to write handshake version: %w", err)
    }
    // Write features (u32)
    if err := binary.Write(writer, binary.LittleEndian, h.Features.Flags()); err != nil {
        return fmt.Errorf("failed to write handshake features: %w", err)
    }
    return nil
}

// DeserializeHandshake reads the handshake from the reader in little-endian format
func DeserializeHandshake(reader io.Reader) (Handshake, error) {
    var version uint32
    var flags uint32

    // Read version
    if err := binary.Read(reader, binary.LittleEndian, &version); err != nil {
        return Handshake{}, fmt.Errorf("failed to read handshake version: %w", err)
    }

    // Read features
    if err := binary.Read(reader, binary.LittleEndian, &flags); err != nil {
        return Handshake{}, fmt.Errorf("failed to read handshake features: %w", err)
    }

    handshake := Handshake{
        Version:  version,
        Features: Features{flags: flags},
    }

    return handshake, nil
}

// ExchangeHandshake performs handshake exchange for active side (client)
// Sends handshake, then receives response handshake
func ExchangeHandshake(stream io.ReadWriter, features Features) (Handshake, error) {
    // Create and send handshake with requested features
    handshake := Handshake{
        Version:  IPC_VERSION,
        Features: features,
    }

    if err := handshake.Serialize(stream); err != nil {
        return Handshake{}, fmt.Errorf("failed to send handshake: %w", err)
    }

    // Receive response handshake
    responseHandshake, err := DeserializeHandshake(stream)
    if err != nil {
        return Handshake{}, fmt.Errorf("failed to receive handshake: %w", err)
    }

    // Check if features are unsupported
    if responseHandshake.IsUnsupported() {
        return Handshake{}, fmt.Errorf("ipc: unsupported features")
    }

    return responseHandshake, nil
}

// ResponseHandshake performs handshake exchange for passive side (server)
// Receives handshake, then sends response handshake
func ResponseHandshake(stream io.ReadWriter, features Features) (Handshake, error) {
    // Receive handshake
    receivedHandshake, err := DeserializeHandshake(stream)
    if err != nil {
        return Handshake{}, fmt.Errorf("failed to receive handshake: %w", err)
    }

    // Create and send response handshake with our features
    handshake := Handshake{
        Version:  IPC_VERSION,
        Features: features,
    }

    if err := handshake.Serialize(stream); err != nil {
        return Handshake{}, fmt.Errorf("failed to send handshake: %w", err)
    }

    // Check if features are unsupported
    if receivedHandshake.IsUnsupported() {
        return Handshake{}, fmt.Errorf("ipc: unsupported features")
    }

    return receivedHandshake, nil
}
