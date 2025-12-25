package tpserde

import (
    "encoding/binary"
    "io"
)

// Feature flags
const (
    FEATURE_COMPRESSED   uint32 = 1
    FEATURE_BUFFERED     uint32 = 1 << 1
    FEATURE_UNSUPPORTED  uint32 = 1 << 31
)

// Default uncompressed size limit (in bytes)
// If AST exceeds this limit and compression is enabled, data will be compressed
const DEFAULT_UNCOMPRESSED_SIZE_LIMIT = 4096

// Features holds feature flags for the message
type Features struct {
    flags uint32
}

func NewFeatures() Features {
    return Features{flags: 0}
}

func (f Features) WithCompressed() Features {
    f.flags |= FEATURE_COMPRESSED
    return f
}

func (f Features) WithBuffered() Features {
    f.flags |= FEATURE_BUFFERED
    return f
}

func (f Features) IsCompressed() bool {
    return (f.flags & FEATURE_COMPRESSED) == FEATURE_COMPRESSED
}

func (f Features) IsBuffered() bool {
    return (f.flags & FEATURE_BUFFERED) == FEATURE_BUFFERED
}

func (f Features) IsUnsupported() bool {
    return (f.flags & FEATURE_UNSUPPORTED) == FEATURE_UNSUPPORTED
}

func (f Features) Flags() uint32 {
    return f.flags
}

// Header represents the serialized message header
// Layout matches Rust #[repr(C)] struct Header:
//   features: u32
//   reserved: u32
//   len: usize (u64 on 64-bit systems)
type Header struct {
    Features Features
    Reserved uint32
    Len      uint64 // Using uint64 for cross-platform compatibility
}

func NewHeader() Header {
    return Header{
        Features: NewFeatures(),
        Reserved: 0,
        Len:      0,
    }
}

func (h Header) WithCompressed() Header {
    h.Features = h.Features.WithCompressed()
    return h
}

func (h Header) WithBuffered() Header {
    h.Features = h.Features.WithBuffered()
    return h
}

func (h *Header) SetLen(len uint64) {
    h.Len = len
}

func (h Header) IsCompressed() bool {
    return h.Features.IsCompressed()
}

func (h Header) IsBuffered() bool {
    return h.Features.IsBuffered()
}

// Serialize writes the header to the writer in little-endian format
func (h Header) Serialize(writer io.Writer) error {
    // Write features (u32)
    if err := binary.Write(writer, binary.LittleEndian, h.Features.Flags()); err != nil {
        return err
    }
    // Write reserved (u32)
    if err := binary.Write(writer, binary.LittleEndian, h.Reserved); err != nil {
        return err
    }
    // Write len (u64)
    if err := binary.Write(writer, binary.LittleEndian, h.Len); err != nil {
        return err
    }
    return nil
}

// Deserialize reads the header from the reader in little-endian format
// Returns nil header and nil error on EOF (connection closed)
func DeserializeHeader(reader io.Reader) (*Header, error) {
    var flags uint32
    var reserved uint32
    var length uint64

    // Read features
    if err := binary.Read(reader, binary.LittleEndian, &flags); err != nil {
        if err == io.EOF {
            return nil, nil // Connection closed
        }
        return nil, err
    }

    // Read reserved
    if err := binary.Read(reader, binary.LittleEndian, &reserved); err != nil {
        return nil, err
    }

    // Read len
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return nil, err
    }

    header := &Header{
        Features: Features{flags: flags},
        Reserved: reserved,
        Len:      length,
    }

    return header, nil
}
