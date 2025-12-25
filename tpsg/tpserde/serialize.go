package tpserde

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "io"

    "github.com/google/uuid"
    "github.com/pierrec/lz4/v4"
)

// TPDataSer serializes TPTypes to theplatform binary format
// compress: if true, compresses data with LZ4 when size exceeds DEFAULT_UNCOMPRESSED_SIZE_LIMIT
func TPDataSer(data TPTypes, compress bool) (TPBinary, error) {
    // Create buffer for payload
    payloadBuf := &bytes.Buffer{}

    // Write AST to payload
    if err := writeAST(payloadBuf, data); err != nil {
        return nil, fmt.Errorf("failed to serialize AST: %w", err)
    }

    payload := payloadBuf.Bytes()

    // Compress if needed
    header := NewHeader()
    if compress && len(payload) > DEFAULT_UNCOMPRESSED_SIZE_LIMIT {
        compressed, err := compressLZ4(payload)
        if err != nil {
            return nil, fmt.Errorf("failed to compress data: %w", err)
        }
        payload = compressed
        header = header.WithCompressed()
    }

    header.SetLen(uint64(len(payload)))

    // Create buffer for final output
    output := &bytes.Buffer{}

    // Write header
    if err := header.Serialize(output); err != nil {
        return nil, fmt.Errorf("failed to serialize header: %w", err)
    }

    // Write payload
    if _, err := output.Write(payload); err != nil {
        return nil, fmt.Errorf("failed to write payload: %w", err)
    }

    return TPBinary(output.Bytes()), nil
}

// compressLZ4 compresses data using LZ4
func compressLZ4(data []byte) ([]byte, error) {
    buf := &bytes.Buffer{}
    writer := lz4.NewWriter(buf)

    if _, err := writer.Write(data); err != nil {
        return nil, err
    }

    if err := writer.Close(); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

// writeAST writes a single TPTypes value to the writer
func writeAST(writer io.Writer, data TPTypes) error {
    // Write type tag
    if err := binary.Write(writer, binary.LittleEndian, data.Tp); err != nil {
        return fmt.Errorf("failed to write type tag: %w", err)
    }

    // Dispatch based on type
    switch data.Tp {
    case NIL, ANY:
        // No data to write
        return nil

    // Scalar types
    case SC_BOOL:
        return writeScalarBool(writer, data.Data.(TPBool))
    case SC_BYTE:
        return writeScalarByte(writer, data.Data.(TPByte))
    case SC_SHORT:
        return writeScalarShort(writer, data.Data.(TPShort))
    case SC_INT:
        return writeScalarInt32(writer, data.Data.(TPInt))
    case SC_MONTH:
        return writeScalarInt32(writer, TPInt(data.Data.(TPMonth)))
    case SC_DATE:
        return writeScalarInt32(writer, TPInt(data.Data.(TPDate)))
    case SC_MINUTE:
        return writeScalarInt32(writer, TPInt(data.Data.(TPMinute)))
    case SC_SECOND:
        return writeScalarInt32(writer, TPInt(data.Data.(TPSecond)))
    case SC_TIME:
        return writeScalarInt32(writer, TPInt(data.Data.(TPTime)))
    case SC_ENUM:
        return writeScalarEnum(writer, data.Data.(TPEnum))
    case SC_LONG:
        return writeScalarInt64(writer, data.Data.(TPLong))
    case SC_TIMESTAMP:
        return writeScalarInt64(writer, TPLong(data.Data.(TPTimestamp)))
    case SC_DATETIME:
        return writeScalarInt64(writer, TPLong(data.Data.(TPDatetime)))
    case SC_TIMESPAN:
        return writeScalarInt64(writer, TPLong(data.Data.(TPTimespan)))
    case SC_REAL:
        return writeScalarReal(writer, data.Data.(TPReal))
    case SC_FLOAT:
        return writeScalarFloat(writer, data.Data.(TPFloat))
    case SC_GUID:
        return writeScalarGUID(writer, data.Data.(TPGUID))
    case SC_SYMBOL, SC_SHADOW:
        return writeScalarSymbol(writer, data.Data.(TPSymbol))

    // Vector types
    case VEC_BOOL:
        return writeVectorBool(writer, data.Data.(TPVecBool))
    case VEC_BYTE:
        return writeVectorByte(writer, data.Data.(TPVecByte))
    case VEC_SHORT:
        return writeVectorShort(writer, data.Data.(TPVecShort))
    case VEC_INT:
        return writeVectorInt32(writer, []int32(data.Data.(TPVecInt)))
    case VEC_MONTH:
        return writeVectorInt32(writer, []int32(data.Data.(TPVecMonth)))
    case VEC_DATE:
        return writeVectorInt32(writer, []int32(data.Data.(TPVecDate)))
    case VEC_MINUTE:
        return writeVectorInt32(writer, []int32(data.Data.(TPVecMinute)))
    case VEC_SECOND:
        return writeVectorInt32(writer, []int32(data.Data.(TPVecSecond)))
    case VEC_TIME:
        return writeVectorInt32(writer, []int32(data.Data.(TPVecTime)))
    case VEC_ENUM:
        return writeVectorEnum(writer, data.Data.(TPVecEnum))
    case VEC_LONG:
        return writeVectorInt64(writer, []int64(data.Data.(TPVecLong)))
    case VEC_TIMESTAMP:
        return writeVectorInt64(writer, []int64(data.Data.(TPVecTimestamp)))
    case VEC_DATETIME:
        return writeVectorInt64(writer, []int64(data.Data.(TPVecDatetime)))
    case VEC_TIMESPAN:
        return writeVectorInt64(writer, []int64(data.Data.(TPVecTimespan)))
    case VEC_REAL:
        return writeVectorReal(writer, data.Data.(TPVecReal))
    case VEC_FLOAT:
        return writeVectorFloat(writer, data.Data.(TPVecFloat))
    case VEC_GUID:
        return writeVectorGUID(writer, data.Data.(TPVecGUID))
    case VEC_SYMBOL, VEC_SHADOW:
        return writeVectorSymbol(writer, data.Data.(TPVecSymbol))
    case VEC_CHAR:
        return writeVectorChar(writer, data.Data.(TPVecChar))

    // Complex types
    case LIST, RETURN, LIST_EXPR:
        return writeList(writer, data.Data.(TPList))
    case DICT:
        return writeDict(writer, data.Data.(TPDict))
    case DICT_TABLE:
        return writeDict(writer, data.Data.(TPDict))
    case TABLE:
        return writeTable(writer, data.Data.(TPTable))
    case PATTERN:
        return writePattern(writer, data.Data.(TPPattern))
    case LAMBDA, CLOSURE:
        return writeLambda(writer, data.Data.(TPLambda))
    case LAMBDA_REC:
        // No data to write
        return nil
    case REAGENT:
        // No data to write
        return nil

    default:
        return fmt.Errorf("unsupported type for serialization: 0x%x", data.Tp)
    }
}

// Scalar writers
func writeScalarBool(writer io.Writer, v TPBool) error {
    return binary.Write(writer, binary.LittleEndian, bool(v))
}

func writeScalarByte(writer io.Writer, v TPByte) error {
    return binary.Write(writer, binary.LittleEndian, int8(v))
}

func writeScalarShort(writer io.Writer, v TPShort) error {
    return binary.Write(writer, binary.LittleEndian, int16(v))
}

func writeScalarInt32(writer io.Writer, v TPInt) error {
    return binary.Write(writer, binary.LittleEndian, int32(v))
}

func writeScalarEnum(writer io.Writer, v TPEnum) error {
    return binary.Write(writer, binary.LittleEndian, uint32(v))
}

func writeScalarInt64(writer io.Writer, v TPLong) error {
    return binary.Write(writer, binary.LittleEndian, int64(v))
}

func writeScalarReal(writer io.Writer, v TPReal) error {
    return binary.Write(writer, binary.LittleEndian, float32(v))
}

func writeScalarFloat(writer io.Writer, v TPFloat) error {
    return binary.Write(writer, binary.LittleEndian, float64(v))
}

func writeScalarGUID(writer io.Writer, v TPGUID) error {
    guid := uuid.UUID(v)
    _, err := writer.Write(guid[:])
    return err
}

func writeScalarSymbol(writer io.Writer, v TPSymbol) error {
    s := string(v)
    if len(s) > 255 {
        return fmt.Errorf("symbol too long: %d bytes (max 255)", len(s))
    }
    // Write length (u8)
    if err := binary.Write(writer, binary.LittleEndian, uint8(len(s))); err != nil {
        return err
    }
    // Write string bytes
    _, err := writer.Write([]byte(s))
    return err
}

// Vector writers
func writeVectorBool(writer io.Writer, v TPVecBool) error {
    vec := []bool(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorByte(writer io.Writer, v TPVecByte) error {
    vec := []int8(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorShort(writer io.Writer, v TPVecShort) error {
    vec := []int16(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorInt32(writer io.Writer, vec []int32) error {
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorEnum(writer io.Writer, v TPVecEnum) error {
    vec := []uint32(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorInt64(writer io.Writer, vec []int64) error {
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorReal(writer io.Writer, v TPVecReal) error {
    vec := []float32(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorFloat(writer io.Writer, v TPVecFloat) error {
    vec := []float64(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write data
    return binary.Write(writer, binary.LittleEndian, vec)
}

func writeVectorGUID(writer io.Writer, v TPVecGUID) error {
    vec := []uuid.UUID(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write each GUID
    for _, guid := range vec {
        if _, err := writer.Write(guid[:]); err != nil {
            return err
        }
    }
    return nil
}

func writeVectorSymbol(writer io.Writer, v TPVecSymbol) error {
    vec := []string(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(vec))); err != nil {
        return err
    }
    // Write each symbol
    for _, s := range vec {
        if len(s) > 255 {
            return fmt.Errorf("symbol too long: %d bytes (max 255)", len(s))
        }
        // Write symbol length (u8)
        if err := binary.Write(writer, binary.LittleEndian, uint8(len(s))); err != nil {
            return err
        }
        // Write symbol bytes
        if _, err := writer.Write([]byte(s)); err != nil {
            return err
        }
    }
    return nil
}

func writeVectorChar(writer io.Writer, v TPVecChar) error {
    s := string(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(s))); err != nil {
        return err
    }
    // Write string bytes
    _, err := writer.Write([]byte(s))
    return err
}

// Complex type writers
func writeList(writer io.Writer, v TPList) error {
    list := []TPTypes(v)
    // Write length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(list))); err != nil {
        return err
    }
    // Write each element
    for _, item := range list {
        if err := writeAST(writer, item); err != nil {
            return err
        }
    }
    return nil
}

func writeDict(writer io.Writer, v TPDict) error {
    // Write keys
    if err := writeAST(writer, v.Keys); err != nil {
        return err
    }
    // Write values
    return writeAST(writer, v.Values)
}

func writeTable(writer io.Writer, v TPTable) error {
    // Write keys
    if err := writeAST(writer, v.Keys); err != nil {
        return err
    }
    // Write values
    return writeAST(writer, v.Values)
}

func writePattern(writer io.Writer, v TPPattern) error {
    // Write exprs
    if err := writeAST(writer, v.Exprs); err != nil {
        return err
    }
    // Write arms
    return writeAST(writer, v.Arms)
}

func writeLambda(writer io.Writer, v TPLambda) error {
    // Write lambda text length
    if err := binary.Write(writer, binary.LittleEndian, uint32(len(v.Text))); err != nil {
        return err
    }
    // Write lambda text
    if _, err := writer.Write([]byte(v.Text)); err != nil {
        return err
    }
    // Write cargs
    if err := binary.Write(writer, binary.LittleEndian, v.Cargs); err != nil {
        return err
    }
    // Write clocals
    if err := binary.Write(writer, binary.LittleEndian, v.Clocals); err != nil {
        return err
    }
    // Write meta
    if err := writeAST(writer, v.Meta.Bind); err != nil {
        return err
    }
    if err := writeAST(writer, v.Meta.Channels); err != nil {
        return err
    }
    if err := writeAST(writer, v.Meta.Args); err != nil {
        return err
    }
    if err := writeAST(writer, v.Meta.Locals); err != nil {
        return err
    }
    // Write body
    if err := writeAST(writer, v.Body); err != nil {
        return err
    }
    // Write upvals
    return writeAST(writer, v.Upvals)
}
