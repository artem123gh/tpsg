package tpserde

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"

    "github.com/google/uuid"
    "github.com/pierrec/lz4/v4"
)

// TPBinary represents binary data in theplatform format
type TPBinary []byte

// TPDataDe deserializes theplatform binary data to TPTypes
// Automatically detects and decompresses LZ4 compressed data
func TPDataDe(data TPBinary) (TPTypes, error) {
    reader := bytes.NewReader(data)

    // Read header
    header, err := DeserializeHeader(reader)
    if err != nil {
        return TPTypes{}, fmt.Errorf("failed to deserialize header: %w", err)
    }
    if header == nil {
        return TPTypes{}, fmt.Errorf("unexpected nil header")
    }

    // Check for unsupported features
    if header.Features.IsUnsupported() {
        return TPTypes{}, fmt.Errorf("unsupported features in header")
    }

    // Read payload
    var payload []byte
    if header.Len > 0 {
        payload = make([]byte, header.Len)
        if _, err := io.ReadFull(reader, payload); err != nil {
            return TPTypes{}, fmt.Errorf("failed to read payload: %w", err)
        }
    } else {
        // If len is 0, read remaining data
        payload, err = io.ReadAll(reader)
        if err != nil {
            return TPTypes{}, fmt.Errorf("failed to read payload: %w", err)
        }
    }

    // Decompress if needed
    if header.IsCompressed() {
        decompressed, err := decompressLZ4(payload)
        if err != nil {
            return TPTypes{}, fmt.Errorf("failed to decompress LZ4 data: %w", err)
        }
        payload = decompressed
    }

    // Deserialize AST
    payloadReader := bytes.NewReader(payload)
    return readAST(payloadReader)
}

// decompressLZ4 decompresses LZ4 compressed data
func decompressLZ4(data []byte) ([]byte, error) {
    reader := lz4.NewReader(bytes.NewReader(data))
    return io.ReadAll(reader)
}

// readAST reads a single AST value from the reader
func readAST(reader io.Reader) (TPTypes, error) {
    // Read type tag
    var tp uint32
    if err := binary.Read(reader, binary.LittleEndian, &tp); err != nil {
        return TPTypes{}, fmt.Errorf("failed to read type tag: %w", err)
    }

    // Dispatch based on type
    switch tp {
    case NIL:
        return NewTPNil(), nil
    case ANY:
        return NewTPAny(), nil

    // Scalar types
    case SC_BOOL:
        return readScalarBool(reader)
    case SC_BYTE:
        return readScalarByte(reader)
    case SC_SHORT:
        return readScalarShort(reader)
    case SC_INT, SC_MONTH, SC_DATE, SC_MINUTE, SC_SECOND, SC_TIME:
        return readScalarInt32(reader, tp)
    case SC_ENUM:
        return readScalarEnum(reader)
    case SC_LONG, SC_TIMESTAMP, SC_DATETIME, SC_TIMESPAN:
        return readScalarInt64(reader, tp)
    case SC_REAL:
        return readScalarReal(reader)
    case SC_FLOAT:
        return readScalarFloat(reader)
    case SC_GUID:
        return readScalarGUID(reader)
    case SC_SYMBOL, SC_SHADOW:
        return readScalarSymbol(reader, tp)

    // Vector types
    case VEC_BOOL:
        return readVectorBool(reader)
    case VEC_BYTE:
        return readVectorByte(reader)
    case VEC_SHORT:
        return readVectorShort(reader)
    case VEC_INT, VEC_MONTH, VEC_DATE, VEC_MINUTE, VEC_SECOND, VEC_TIME:
        return readVectorInt32(reader, tp)
    case VEC_ENUM:
        return readVectorEnum(reader)
    case VEC_LONG, VEC_TIMESTAMP, VEC_DATETIME, VEC_TIMESPAN:
        return readVectorInt64(reader, tp)
    case VEC_REAL:
        return readVectorReal(reader)
    case VEC_FLOAT:
        return readVectorFloat(reader)
    case VEC_GUID:
        return readVectorGUID(reader)
    case VEC_SYMBOL, VEC_SHADOW:
        return readVectorSymbol(reader, tp)
    case VEC_CHAR:
        return readVectorChar(reader)

    // Complex types
    case LIST, RETURN, LIST_EXPR:
        return readList(reader, tp)
    case DICT:
        return readDict(reader)
    case DICT_TABLE:
        return readDictTable(reader)
    case TABLE:
        return readTable(reader)
    case PATTERN:
        return readPattern(reader)
    case LAMBDA, CLOSURE:
        return readLambda(reader, tp)
    case LAMBDA_REC:
        return TPTypes{Tp: LAMBDA_REC, Data: nil}, nil
    case REAGENT:
        return TPTypes{Tp: REAGENT, Data: nil}, nil

    default:
        return TPTypes{}, fmt.Errorf("unsupported type: 0x%x", tp)
    }
}

// Scalar readers
func readScalarBool(reader io.Reader) (TPTypes, error) {
    var v bool
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    return NewTPBool(v), nil
}

func readScalarByte(reader io.Reader) (TPTypes, error) {
    var v int8
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    return NewTPByte(v), nil
}

func readScalarShort(reader io.Reader) (TPTypes, error) {
    var v int16
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    return NewTPShort(v), nil
}

func readScalarInt32(reader io.Reader, tp uint32) (TPTypes, error) {
    var v int32
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    switch tp {
    case SC_MONTH:
        return NewTPMonth(v), nil
    case SC_DATE:
        return NewTPDate(v), nil
    case SC_MINUTE:
        return NewTPMinute(v), nil
    case SC_SECOND:
        return NewTPSecond(v), nil
    case SC_TIME:
        return NewTPTime(v), nil
    default:
        return NewTPInt(v), nil
    }
}

func readScalarEnum(reader io.Reader) (TPTypes, error) {
    var v uint32
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    return NewTPEnum(v), nil
}

func readScalarInt64(reader io.Reader, tp uint32) (TPTypes, error) {
    var v int64
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    switch tp {
    case SC_TIMESTAMP:
        return NewTPTimestamp(v), nil
    case SC_DATETIME:
        return NewTPDatetime(v), nil
    case SC_TIMESPAN:
        return NewTPTimespan(v), nil
    default:
        return NewTPLong(v), nil
    }
}

func readScalarReal(reader io.Reader) (TPTypes, error) {
    var v float32
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    return NewTPReal(v), nil
}

func readScalarFloat(reader io.Reader) (TPTypes, error) {
    var v float64
    if err := binary.Read(reader, binary.LittleEndian, &v); err != nil {
        return TPTypes{}, err
    }
    return NewTPFloat(v), nil
}

func readScalarGUID(reader io.Reader) (TPTypes, error) {
    var buf [16]byte
    if _, err := io.ReadFull(reader, buf[:]); err != nil {
        return TPTypes{}, err
    }
    guid, err := uuid.FromBytes(buf[:])
    if err != nil {
        return TPTypes{}, err
    }
    return NewTPGUID(guid), nil
}

func readScalarSymbol(reader io.Reader, tp uint32) (TPTypes, error) {
    // Read length (u8)
    var length uint8
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    // Read string bytes
    buf := make([]byte, length)
    if _, err := io.ReadFull(reader, buf); err != nil {
        return TPTypes{}, err
    }
    return NewTPSymbol(string(buf)), nil
}

// Vector readers
func readVectorBool(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]bool, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecBool(vec), nil
}

func readVectorByte(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]int8, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecByte(vec), nil
}

func readVectorShort(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]int16, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecShort(vec), nil
}

func readVectorInt32(reader io.Reader, tp uint32) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]int32, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    switch tp {
    case VEC_MONTH:
        return NewTPVecMonth(vec), nil
    case VEC_DATE:
        return NewTPVecDate(vec), nil
    case VEC_MINUTE:
        return NewTPVecMinute(vec), nil
    case VEC_SECOND:
        return NewTPVecSecond(vec), nil
    case VEC_TIME:
        return NewTPVecTime(vec), nil
    default:
        return NewTPVecInt(vec), nil
    }
}

func readVectorEnum(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]uint32, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecEnum(vec), nil
}

func readVectorInt64(reader io.Reader, tp uint32) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]int64, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    switch tp {
    case VEC_TIMESTAMP:
        return NewTPVecTimestamp(vec), nil
    case VEC_DATETIME:
        return NewTPVecDatetime(vec), nil
    case VEC_TIMESPAN:
        return NewTPVecTimespan(vec), nil
    default:
        return NewTPVecLong(vec), nil
    }
}

func readVectorReal(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]float32, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecReal(vec), nil
}

func readVectorFloat(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]float64, length)
    if err := binary.Read(reader, binary.LittleEndian, &vec); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecFloat(vec), nil
}

func readVectorGUID(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]uuid.UUID, length)
    for i := uint32(0); i < length; i++ {
        var buf [16]byte
        if _, err := io.ReadFull(reader, buf[:]); err != nil {
            return TPTypes{}, err
        }
        guid, err := uuid.FromBytes(buf[:])
        if err != nil {
            return TPTypes{}, err
        }
        vec[i] = guid
    }
    return NewTPVecGUID(vec), nil
}

func readVectorSymbol(reader io.Reader, tp uint32) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    vec := make([]string, length)
    for i := uint32(0); i < length; i++ {
        // Read symbol length (u8)
        var symLen uint8
        if err := binary.Read(reader, binary.LittleEndian, &symLen); err != nil {
            return TPTypes{}, err
        }
        // Read symbol bytes
        buf := make([]byte, symLen)
        if _, err := io.ReadFull(reader, buf); err != nil {
            return TPTypes{}, err
        }
        vec[i] = string(buf)
    }
    return NewTPVecSymbol(vec), nil
}

func readVectorChar(reader io.Reader) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    buf := make([]byte, length)
    if _, err := io.ReadFull(reader, buf); err != nil {
        return TPTypes{}, err
    }
    return NewTPVecChar(string(buf)), nil
}

// Complex type readers
func readList(reader io.Reader, tp uint32) (TPTypes, error) {
    var length uint32
    if err := binary.Read(reader, binary.LittleEndian, &length); err != nil {
        return TPTypes{}, err
    }
    list := make([]TPTypes, length)
    for i := uint32(0); i < length; i++ {
        item, err := readAST(reader)
        if err != nil {
            return TPTypes{}, err
        }
        list[i] = item
    }
    result := NewTPList(list)
    result.Tp = tp // Preserve original type (LIST, RETURN, or LIST_EXPR)
    return result, nil
}

func readDict(reader io.Reader) (TPTypes, error) {
    keys, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    values, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    return NewTPDict(keys, values), nil
}

func readDictTable(reader io.Reader) (TPTypes, error) {
    keys, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    values, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    return TPTypes{Tp: DICT_TABLE, Data: TPDict{Keys: keys, Values: values}}, nil
}

func readTable(reader io.Reader) (TPTypes, error) {
    keys, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    values, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    return NewTPTable(keys, values), nil
}

func readPattern(reader io.Reader) (TPTypes, error) {
    exprs, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    arms, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    return NewTPPattern(exprs, arms), nil
}

func readLambda(reader io.Reader, tp uint32) (TPTypes, error) {
    // Read lambda text length
    var textLen uint32
    if err := binary.Read(reader, binary.LittleEndian, &textLen); err != nil {
        return TPTypes{}, err
    }
    // Read lambda text
    textBuf := make([]byte, textLen)
    if _, err := io.ReadFull(reader, textBuf); err != nil {
        return TPTypes{}, err
    }
    text := string(textBuf)

    // Read cargs
    var cargs uint16
    if err := binary.Read(reader, binary.LittleEndian, &cargs); err != nil {
        return TPTypes{}, err
    }

    // Read clocals
    var clocals uint16
    if err := binary.Read(reader, binary.LittleEndian, &clocals); err != nil {
        return TPTypes{}, err
    }

    // Read meta
    bind, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    channels, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    args, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }
    locals, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }

    meta := TPLambdaMeta{
        Bind:     bind,
        Channels: channels,
        Args:     args,
        Locals:   locals,
    }

    // Read body
    body, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }

    // Read upvals
    upvals, err := readAST(reader)
    if err != nil {
        return TPTypes{}, err
    }

    result := NewTPLambda(text, cargs, clocals, meta, body, upvals)
    result.Tp = tp // Preserve original type (LAMBDA or CLOSURE)
    return result, nil
}
