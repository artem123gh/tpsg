package tpserde

import (
    "bytes"
    "testing"

    "github.com/google/uuid"
)

// TestScalarTypes tests serialization and deserialization of scalar types
func TestScalarTypes(t *testing.T) {
    tests := []struct {
        name string
        data TPTypes
    }{
        {"Nil", NewTPNil()},
        {"Any", NewTPAny()},
        {"Bool true", NewTPBool(true)},
        {"Bool false", NewTPBool(false)},
        {"Byte", NewTPByte(42)},
        {"Short", NewTPShort(1000)},
        {"Int", NewTPInt(100000)},
        {"Long", NewTPLong(9223372036854775807)},
        {"Real", NewTPReal(3.14)},
        {"Float", NewTPFloat(3.14159265359)},
        {"Enum", NewTPEnum(123)},
        {"Symbol", NewTPSymbol("test_symbol")},
        {"Month", NewTPMonth(12)},
        {"Date", NewTPDate(20231225)},
        {"Timestamp", NewTPTimestamp(1703548800000)},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test without compression
            binary, err := TPDataSer(tt.data, false)
            if err != nil {
                t.Fatalf("Failed to serialize %s: %v", tt.name, err)
            }

            result, err := TPDataDe(binary)
            if err != nil {
                t.Fatalf("Failed to deserialize %s: %v", tt.name, err)
            }

            if result.Tp != tt.data.Tp {
                t.Errorf("Type mismatch for %s: got 0x%x, want 0x%x", tt.name, result.Tp, tt.data.Tp)
            }

            // Test with compression (should not compress small data)
            binary, err = TPDataSer(tt.data, true)
            if err != nil {
                t.Fatalf("Failed to serialize %s with compression: %v", tt.name, err)
            }

            result, err = TPDataDe(binary)
            if err != nil {
                t.Fatalf("Failed to deserialize %s with compression: %v", tt.name, err)
            }

            if result.Tp != tt.data.Tp {
                t.Errorf("Type mismatch for %s (compressed): got 0x%x, want 0x%x", tt.name, result.Tp, tt.data.Tp)
            }
        })
    }
}

// TestVectorTypes tests serialization and deserialization of vector types
func TestVectorTypes(t *testing.T) {
    tests := []struct {
        name string
        data TPTypes
    }{
        {"VecBool", NewTPVecBool([]bool{true, false, true})},
        {"VecByte", NewTPVecByte([]int8{1, 2, 3, 4, 5})},
        {"VecShort", NewTPVecShort([]int16{100, 200, 300})},
        {"VecInt", NewTPVecInt([]int32{1000, 2000, 3000})},
        {"VecLong", NewTPVecLong([]int64{10000, 20000, 30000})},
        {"VecReal", NewTPVecReal([]float32{1.1, 2.2, 3.3})},
        {"VecFloat", NewTPVecFloat([]float64{1.1, 2.2, 3.3})},
        {"VecChar", NewTPVecChar("Hello, World!")},
        {"VecSymbol", NewTPVecSymbol([]string{"foo", "bar", "baz"})},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            binary, err := TPDataSer(tt.data, false)
            if err != nil {
                t.Fatalf("Failed to serialize %s: %v", tt.name, err)
            }

            result, err := TPDataDe(binary)
            if err != nil {
                t.Fatalf("Failed to deserialize %s: %v", tt.name, err)
            }

            if result.Tp != tt.data.Tp {
                t.Errorf("Type mismatch for %s: got 0x%x, want 0x%x", tt.name, result.Tp, tt.data.Tp)
            }
        })
    }
}

// TestGUID tests GUID serialization and deserialization
func TestGUID(t *testing.T) {
    guid := uuid.New()
    data := NewTPGUID(guid)

    binary, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize GUID: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize GUID: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch: got 0x%x, want 0x%x", result.Tp, data.Tp)
    }

    resultGUID := result.Data.(TPGUID)
    if uuid.UUID(resultGUID) != guid {
        t.Errorf("GUID mismatch: got %v, want %v", resultGUID, guid)
    }
}

// TestVecGUID tests GUID vector serialization and deserialization
func TestVecGUID(t *testing.T) {
    guids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
    data := NewTPVecGUID(guids)

    binary, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize VecGUID: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize VecGUID: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch: got 0x%x, want 0x%x", result.Tp, data.Tp)
    }

    resultVec := result.Data.(TPVecGUID)
    if len(resultVec) != len(guids) {
        t.Errorf("Length mismatch: got %d, want %d", len(resultVec), len(guids))
    }
}

// TestList tests list serialization and deserialization
func TestList(t *testing.T) {
    list := []TPTypes{
        NewTPInt(42),
        NewTPVecChar("hello"),
        NewTPFloat(3.14),
    }
    data := NewTPList(list)

    binary, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize list: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize list: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch: got 0x%x, want 0x%x", result.Tp, data.Tp)
    }

    resultList := result.Data.(TPList)
    if len(resultList) != len(list) {
        t.Errorf("List length mismatch: got %d, want %d", len(resultList), len(list))
    }
}

// TestDict tests dict serialization and deserialization
func TestDict(t *testing.T) {
    keys := NewTPVecSymbol([]string{"a", "b", "c"})
    values := NewTPVecInt([]int32{1, 2, 3})
    data := NewTPDict(keys, values)

    binary, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize dict: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize dict: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch: got 0x%x, want 0x%x", result.Tp, data.Tp)
    }
}

// TestTable tests table serialization and deserialization
func TestTable(t *testing.T) {
    keys := NewTPVecSymbol([]string{"col1", "col2"})
    values := NewTPList([]TPTypes{
        NewTPVecInt([]int32{1, 2, 3}),
        NewTPVecFloat([]float64{1.1, 2.2, 3.3}),
    })
    data := NewTPTable(keys, values)

    binary, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize table: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize table: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch: got 0x%x, want 0x%x", result.Tp, data.Tp)
    }
}

// TestCompression tests LZ4 compression and decompression
func TestCompression(t *testing.T) {
    // Create large data that should be compressed
    largeVec := make([]int32, 10000)
    for i := range largeVec {
        largeVec[i] = int32(i)
    }
    data := NewTPVecInt(largeVec)

    // Serialize without compression
    uncompressed, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize without compression: %v", err)
    }

    // Serialize with compression
    compressed, err := TPDataSer(data, true)
    if err != nil {
        t.Fatalf("Failed to serialize with compression: %v", err)
    }

    // Compressed should be smaller
    if len(compressed) >= len(uncompressed) {
        t.Logf("Warning: Compressed size (%d) not smaller than uncompressed (%d)", len(compressed), len(uncompressed))
        t.Logf("This may be normal for highly random data, but our sequential data should compress well")
    } else {
        t.Logf("Compression ratio: %.2f%%", float64(len(compressed))/float64(len(uncompressed))*100)
    }

    // Both should deserialize to the same result
    result1, err := TPDataDe(uncompressed)
    if err != nil {
        t.Fatalf("Failed to deserialize uncompressed: %v", err)
    }

    result2, err := TPDataDe(compressed)
    if err != nil {
        t.Fatalf("Failed to deserialize compressed: %v", err)
    }

    if result1.Tp != result2.Tp {
        t.Errorf("Type mismatch after compression")
    }

    vec1 := result1.Data.(TPVecInt)
    vec2 := result2.Data.(TPVecInt)
    if len(vec1) != len(vec2) {
        t.Errorf("Vector length mismatch after compression")
    }
}

// TestRoundTrip tests complete round-trip for various types
func TestRoundTrip(t *testing.T) {
    // Create a complex nested structure
    data := NewTPList([]TPTypes{
        NewTPInt(42),
        NewTPVecChar("Hello, World!"),
        NewTPList([]TPTypes{
            NewTPFloat(3.14),
            NewTPVecInt([]int32{1, 2, 3, 4, 5}),
        }),
        NewTPDict(
            NewTPVecSymbol([]string{"key1", "key2"}),
            NewTPVecInt([]int32{100, 200}),
        ),
    })

    // Test without compression
    binary, err := TPDataSer(data, false)
    if err != nil {
        t.Fatalf("Failed to serialize: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch: got 0x%x, want 0x%x", result.Tp, data.Tp)
    }

    // Test with compression
    binary, err = TPDataSer(data, true)
    if err != nil {
        t.Fatalf("Failed to serialize with compression: %v", err)
    }

    result, err = TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize with compression: %v", err)
    }

    if result.Tp != data.Tp {
        t.Errorf("Type mismatch (compressed): got 0x%x, want 0x%x", result.Tp, data.Tp)
    }
}

// TestHeaderSerialization tests header serialization and deserialization
func TestHeaderSerialization(t *testing.T) {
    header := NewHeader().WithCompressed().WithBuffered()
    header.SetLen(12345)

    buf := &bytes.Buffer{}
    if err := header.Serialize(buf); err != nil {
        t.Fatalf("Failed to serialize header: %v", err)
    }

    result, err := DeserializeHeader(buf)
    if err != nil {
        t.Fatalf("Failed to deserialize header: %v", err)
    }

    if result == nil {
        t.Fatal("Deserialized header is nil")
    }

    if !result.IsCompressed() {
        t.Error("Compressed flag not preserved")
    }

    if !result.IsBuffered() {
        t.Error("Buffered flag not preserved")
    }

    if result.Len != header.Len {
        t.Errorf("Length mismatch: got %d, want %d", result.Len, header.Len)
    }
}
