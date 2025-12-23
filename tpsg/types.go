package main

type Status bool

const (
	Success Status = true
	Failure Status = false
)

// Type for user creds.
type TUserCreds struct {
	Username string
	Password string
}

// Type for TOML config.
type TConfigTOML struct {
	TCP uint16 `toml:"TCP"`
	WS  uint16 `toml:"WS"`
}

type GKVSValueType int

const (
	GKVSTypeNone GKVSValueType = iota
	GKVSTypeInt8
	GKVSTypeUInt8
	GKVSTypeInt16
	GKVSTypeUInt16
	GKVSTypeInt32
	GKVSTypeUInt32
	GKVSTypeInt64
	GKVSTypeUInt64
	GKVSTypeFloat32
	GKVSTypeFloat64
	GKVSTypeString
	GKVSTypeTUserCreds
	GKVSTypeTConfigTOML
)

// GKVSTypes is an enum-like type for supported value types for global key-value storage - GKVS.
type GKVSTypes struct {
	Type        GKVSValueType
	Int8        int8
	UInt8       uint8
	Int16       int16
	UInt16      uint16
	Int32       int32
	UInt32      uint32
	Int64       int64
	UInt64      uint64
	Float32     float32
	Float64     float64
	String      string
	TUserCreds  TUserCreds
	TConfigTOML TConfigTOML
}

func NewGKVSNone() GKVSTypes {
	return GKVSTypes{Type: GKVSTypeNone}
}

func NewGKVSInt8(value int8) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeInt8, Int8: value}
}

func NewGKVSUInt8(value uint8) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeUInt8, UInt8: value}
}

func NewGKVSInt16(value int16) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeInt16, Int16: value}
}

func NewGKVSUInt16(value uint16) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeUInt16, UInt16: value}
}

func NewGKVSInt32(value int32) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeInt32, Int32: value}
}

func NewGKVSUInt32(value uint32) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeUInt32, UInt32: value}
}

func NewGKVSInt64(value int64) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeInt64, Int64: value}
}

func NewGKVSUInt64(value uint64) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeUInt64, UInt64: value}
}

func NewGKVSFloat32(value float32) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeFloat32, Float32: value}
}

func NewGKVSFloat64(value float64) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeFloat64, Float64: value}
}

func NewGKVSString(value string) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeString, String: value}
}

func NewGKVSTUserCreds(value TUserCreds) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeTUserCreds, TUserCreds: value}
}

func NewGKVSTConfigTOML(value TConfigTOML) GKVSTypes {
	return GKVSTypes{Type: GKVSTypeTConfigTOML, TConfigTOML: value}
}
