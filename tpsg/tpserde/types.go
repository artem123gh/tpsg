package tpserde

import (
    "github.com/google/uuid"
)

// TPTypes represents theplatform data types in Go
type TPTypes struct {
    Tp   uint32      // Type tag
    Data interface{} // Actual data
}

// Scalar types
type TPNil struct{}
type TPAny struct{}
type TPBool bool
type TPByte int8
type TPShort int16
type TPInt int32
type TPLong int64
type TPReal float32
type TPFloat float64
type TPEnum uint32
type TPGUID uuid.UUID
type TPSymbol string
type TPChar byte

// Temporal types
type TPMonth int32
type TPDate int32
type TPMinute int32
type TPSecond int32
type TPTime int32
type TPTimestamp int64
type TPDatetime int64
type TPTimespan int64

// Vector types
type TPVecBool []bool
type TPVecByte []int8
type TPVecShort []int16
type TPVecInt []int32
type TPVecLong []int64
type TPVecReal []float32
type TPVecFloat []float64
type TPVecEnum []uint32
type TPVecGUID []uuid.UUID
type TPVecSymbol []string
type TPVecChar string // String is essentially []byte in Go

// Temporal vector types
type TPVecMonth []int32
type TPVecDate []int32
type TPVecMinute []int32
type TPVecSecond []int32
type TPVecTime []int32
type TPVecTimestamp []int64
type TPVecDatetime []int64
type TPVecTimespan []int64

// Complex types
type TPList []TPTypes
type TPDict struct {
    Keys   TPTypes
    Values TPTypes
}
type TPTable struct {
    Keys   TPTypes
    Values TPTypes
}
type TPPattern struct {
    Exprs TPTypes
    Arms  TPTypes
}
type TPLambda struct {
    Text    string
    Cargs   uint16
    Clocals uint16
    Meta    TPLambdaMeta
    Body    TPTypes
    Upvals  TPTypes
}
type TPLambdaMeta struct {
    Bind     TPTypes
    Channels TPTypes
    Args     TPTypes
    Locals   TPTypes
}

// Constructor functions for scalar types
func NewTPNil() TPTypes {
    return TPTypes{Tp: NIL, Data: TPNil{}}
}

func NewTPAny() TPTypes {
    return TPTypes{Tp: ANY, Data: TPAny{}}
}

func NewTPBool(v bool) TPTypes {
    return TPTypes{Tp: SC_BOOL, Data: TPBool(v)}
}

func NewTPByte(v int8) TPTypes {
    return TPTypes{Tp: SC_BYTE, Data: TPByte(v)}
}

func NewTPShort(v int16) TPTypes {
    return TPTypes{Tp: SC_SHORT, Data: TPShort(v)}
}

func NewTPInt(v int32) TPTypes {
    return TPTypes{Tp: SC_INT, Data: TPInt(v)}
}

func NewTPLong(v int64) TPTypes {
    return TPTypes{Tp: SC_LONG, Data: TPLong(v)}
}

func NewTPReal(v float32) TPTypes {
    return TPTypes{Tp: SC_REAL, Data: TPReal(v)}
}

func NewTPFloat(v float64) TPTypes {
    return TPTypes{Tp: SC_FLOAT, Data: TPFloat(v)}
}

func NewTPEnum(v uint32) TPTypes {
    return TPTypes{Tp: SC_ENUM, Data: TPEnum(v)}
}

func NewTPGUID(v uuid.UUID) TPTypes {
    return TPTypes{Tp: SC_GUID, Data: TPGUID(v)}
}

func NewTPSymbol(v string) TPTypes {
    return TPTypes{Tp: SC_SYMBOL, Data: TPSymbol(v)}
}

// Temporal type constructors
func NewTPMonth(v int32) TPTypes {
    return TPTypes{Tp: SC_MONTH, Data: TPMonth(v)}
}

func NewTPDate(v int32) TPTypes {
    return TPTypes{Tp: SC_DATE, Data: TPDate(v)}
}

func NewTPMinute(v int32) TPTypes {
    return TPTypes{Tp: SC_MINUTE, Data: TPMinute(v)}
}

func NewTPSecond(v int32) TPTypes {
    return TPTypes{Tp: SC_SECOND, Data: TPSecond(v)}
}

func NewTPTime(v int32) TPTypes {
    return TPTypes{Tp: SC_TIME, Data: TPTime(v)}
}

func NewTPTimestamp(v int64) TPTypes {
    return TPTypes{Tp: SC_TIMESTAMP, Data: TPTimestamp(v)}
}

func NewTPDatetime(v int64) TPTypes {
    return TPTypes{Tp: SC_DATETIME, Data: TPDatetime(v)}
}

func NewTPTimespan(v int64) TPTypes {
    return TPTypes{Tp: SC_TIMESPAN, Data: TPTimespan(v)}
}

// Vector type constructors
func NewTPVecBool(v []bool) TPTypes {
    return TPTypes{Tp: VEC_BOOL, Data: TPVecBool(v)}
}

func NewTPVecByte(v []int8) TPTypes {
    return TPTypes{Tp: VEC_BYTE, Data: TPVecByte(v)}
}

func NewTPVecShort(v []int16) TPTypes {
    return TPTypes{Tp: VEC_SHORT, Data: TPVecShort(v)}
}

func NewTPVecInt(v []int32) TPTypes {
    return TPTypes{Tp: VEC_INT, Data: TPVecInt(v)}
}

func NewTPVecLong(v []int64) TPTypes {
    return TPTypes{Tp: VEC_LONG, Data: TPVecLong(v)}
}

func NewTPVecReal(v []float32) TPTypes {
    return TPTypes{Tp: VEC_REAL, Data: TPVecReal(v)}
}

func NewTPVecFloat(v []float64) TPTypes {
    return TPTypes{Tp: VEC_FLOAT, Data: TPVecFloat(v)}
}

func NewTPVecEnum(v []uint32) TPTypes {
    return TPTypes{Tp: VEC_ENUM, Data: TPVecEnum(v)}
}

func NewTPVecGUID(v []uuid.UUID) TPTypes {
    return TPTypes{Tp: VEC_GUID, Data: TPVecGUID(v)}
}

func NewTPVecSymbol(v []string) TPTypes {
    return TPTypes{Tp: VEC_SYMBOL, Data: TPVecSymbol(v)}
}

func NewTPVecChar(v string) TPTypes {
    return TPTypes{Tp: VEC_CHAR, Data: TPVecChar(v)}
}

// Temporal vector constructors
func NewTPVecMonth(v []int32) TPTypes {
    return TPTypes{Tp: VEC_MONTH, Data: TPVecMonth(v)}
}

func NewTPVecDate(v []int32) TPTypes {
    return TPTypes{Tp: VEC_DATE, Data: TPVecDate(v)}
}

func NewTPVecMinute(v []int32) TPTypes {
    return TPTypes{Tp: VEC_MINUTE, Data: TPVecMinute(v)}
}

func NewTPVecSecond(v []int32) TPTypes {
    return TPTypes{Tp: VEC_SECOND, Data: TPVecSecond(v)}
}

func NewTPVecTime(v []int32) TPTypes {
    return TPTypes{Tp: VEC_TIME, Data: TPVecTime(v)}
}

func NewTPVecTimestamp(v []int64) TPTypes {
    return TPTypes{Tp: VEC_TIMESTAMP, Data: TPVecTimestamp(v)}
}

func NewTPVecDatetime(v []int64) TPTypes {
    return TPTypes{Tp: VEC_DATETIME, Data: TPVecDatetime(v)}
}

func NewTPVecTimespan(v []int64) TPTypes {
    return TPTypes{Tp: VEC_TIMESPAN, Data: TPVecTimespan(v)}
}

// Complex type constructors
func NewTPList(v []TPTypes) TPTypes {
    return TPTypes{Tp: LIST, Data: TPList(v)}
}

func NewTPDict(keys, values TPTypes) TPTypes {
    return TPTypes{Tp: DICT, Data: TPDict{Keys: keys, Values: values}}
}

func NewTPTable(keys, values TPTypes) TPTypes {
    return TPTypes{Tp: TABLE, Data: TPTable{Keys: keys, Values: values}}
}

func NewTPPattern(exprs, arms TPTypes) TPTypes {
    return TPTypes{Tp: PATTERN, Data: TPPattern{Exprs: exprs, Arms: arms}}
}

func NewTPLambda(text string, cargs, clocals uint16, meta TPLambdaMeta, body, upvals TPTypes) TPTypes {
    return TPTypes{Tp: LAMBDA, Data: TPLambda{Text: text, Cargs: cargs, Clocals: clocals, Meta: meta, Body: body, Upvals: upvals}}
}
