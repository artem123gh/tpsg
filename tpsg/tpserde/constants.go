package tpserde

import (
    "math"
)

// NULL values for each numeric type
// These match theplatform's NULL representations
const (
    NULL_BOOL   bool    = false
    NULL_BYTE   int8    = -128      // 0x80 as i8
    NULL_SHORT  int16   = -32768    // 0x8000 as i16
    NULL_INT    int32   = -2147483648 // 0x80000000 as i32
    NULL_LONG   int64   = -9223372036854775808 // 0x8000000000000000 as i64
    NULL_ENUM   uint32  = 0x80000000
    NULL_SYM    int64   = 0
    NULL_DOM    uint32  = 0x80000000
)

// NULL values for floats are negative NaN
var (
    NULL_REAL  float32 = float32(math.NaN()) * -1  // negative NaN
    NULL_FLOAT float64 = math.NaN() * -1           // negative NaN
)

// Positive infinity values
const (
    INF_BYTE  int8  = 127          // 0x7F
    INF_SHORT int16 = 32767        // 0x7FFF
    INF_INT   int32 = 2147483647   // 0x7FFFFFFF
    INF_LONG  int64 = 9223372036854775807 // 0x7FFFFFFFFFFFFFFF
)

var (
    INF_REAL  float32 = float32(math.Inf(1))
    INF_FLOAT float64 = math.Inf(1)
)

// Negative infinity values
const (
    NEG_INF_BYTE  int8  = -127
    NEG_INF_SHORT int16 = -32767
    NEG_INF_INT   int32 = -2147483647
    NEG_INF_LONG  int64 = -9223372036854775807
)

var (
    NEG_INF_REAL  float32 = float32(math.Inf(-1))
    NEG_INF_FLOAT float64 = math.Inf(-1)
)

// Helper functions to check for NULL values
func IsNullByte(v int8) bool {
    return v == NULL_BYTE
}

func IsNullShort(v int16) bool {
    return v == NULL_SHORT
}

func IsNullInt(v int32) bool {
    return v == NULL_INT
}

func IsNullLong(v int64) bool {
    return v == NULL_LONG
}

func IsNullReal(v float32) bool {
    return math.IsNaN(float64(v)) && math.Signbit(float64(v))
}

func IsNullFloat(v float64) bool {
    return math.IsNaN(v) && math.Signbit(v)
}

func IsNullEnum(v uint32) bool {
    return v == NULL_ENUM
}

// Helper functions to check for infinity values
func IsInfByte(v int8) bool {
    return v == INF_BYTE || v == NEG_INF_BYTE
}

func IsInfShort(v int16) bool {
    return v == INF_SHORT || v == NEG_INF_SHORT
}

func IsInfInt(v int32) bool {
    return v == INF_INT || v == NEG_INF_INT
}

func IsInfLong(v int64) bool {
    return v == INF_LONG || v == NEG_INF_LONG
}

func IsInfReal(v float32) bool {
    return math.IsInf(float64(v), 0)
}

func IsInfFloat(v float64) bool {
    return math.IsInf(v, 0)
}

// Helper functions to check for positive infinity
func IsPosInfByte(v int8) bool {
    return v == INF_BYTE
}

func IsPosInfShort(v int16) bool {
    return v == INF_SHORT
}

func IsPosInfInt(v int32) bool {
    return v == INF_INT
}

func IsPosInfLong(v int64) bool {
    return v == INF_LONG
}

func IsPosInfReal(v float32) bool {
    return math.IsInf(float64(v), 1)
}

func IsPosInfFloat(v float64) bool {
    return math.IsInf(v, 1)
}

// Helper functions to check for negative infinity
func IsNegInfByte(v int8) bool {
    return v == NEG_INF_BYTE
}

func IsNegInfShort(v int16) bool {
    return v == NEG_INF_SHORT
}

func IsNegInfInt(v int32) bool {
    return v == NEG_INF_INT
}

func IsNegInfLong(v int64) bool {
    return v == NEG_INF_LONG
}

func IsNegInfReal(v float32) bool {
    return math.IsInf(float64(v), -1)
}

func IsNegInfFloat(v float64) bool {
    return math.IsInf(v, -1)
}

// Constructor functions for NULL values
func NewTPNullByte() TPTypes {
    return NewTPByte(NULL_BYTE)
}

func NewTPNullShort() TPTypes {
    return NewTPShort(NULL_SHORT)
}

func NewTPNullInt() TPTypes {
    return NewTPInt(NULL_INT)
}

func NewTPNullLong() TPTypes {
    return NewTPLong(NULL_LONG)
}

func NewTPNullReal() TPTypes {
    return NewTPReal(NULL_REAL)
}

func NewTPNullFloat() TPTypes {
    return NewTPFloat(NULL_FLOAT)
}

func NewTPNullEnum() TPTypes {
    return NewTPEnum(NULL_ENUM)
}

// Constructor functions for positive infinity
func NewTPInfByte() TPTypes {
    return NewTPByte(INF_BYTE)
}

func NewTPInfShort() TPTypes {
    return NewTPShort(INF_SHORT)
}

func NewTPInfInt() TPTypes {
    return NewTPInt(INF_INT)
}

func NewTPInfLong() TPTypes {
    return NewTPLong(INF_LONG)
}

func NewTPInfReal() TPTypes {
    return NewTPReal(INF_REAL)
}

func NewTPInfFloat() TPTypes {
    return NewTPFloat(INF_FLOAT)
}

// Constructor functions for negative infinity
func NewTPNegInfByte() TPTypes {
    return NewTPByte(NEG_INF_BYTE)
}

func NewTPNegInfShort() TPTypes {
    return NewTPShort(NEG_INF_SHORT)
}

func NewTPNegInfInt() TPTypes {
    return NewTPInt(NEG_INF_INT)
}

func NewTPNegInfLong() TPTypes {
    return NewTPLong(NEG_INF_LONG)
}

func NewTPNegInfReal() TPTypes {
    return NewTPReal(NEG_INF_REAL)
}

func NewTPNegInfFloat() TPTypes {
    return NewTPFloat(NEG_INF_FLOAT)
}
