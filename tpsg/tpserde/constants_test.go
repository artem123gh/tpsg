package tpserde

import (
    "testing"
)

// TestNullValues tests NULL value constants and helpers
func TestNullValues(t *testing.T) {
    // Test NULL detection
    if !IsNullByte(NULL_BYTE) {
        t.Error("NULL_BYTE not detected as null")
    }
    if !IsNullShort(NULL_SHORT) {
        t.Error("NULL_SHORT not detected as null")
    }
    if !IsNullInt(NULL_INT) {
        t.Error("NULL_INT not detected as null")
    }
    if !IsNullLong(NULL_LONG) {
        t.Error("NULL_LONG not detected as null")
    }
    if !IsNullReal(NULL_REAL) {
        t.Error("NULL_REAL not detected as null")
    }
    if !IsNullFloat(NULL_FLOAT) {
        t.Error("NULL_FLOAT not detected as null")
    }
    if !IsNullEnum(NULL_ENUM) {
        t.Error("NULL_ENUM not detected as null")
    }

    // Test non-NULL values
    if IsNullByte(0) {
        t.Error("0 incorrectly detected as NULL_BYTE")
    }
    if IsNullInt(0) {
        t.Error("0 incorrectly detected as NULL_INT")
    }
}

// TestInfinityValues tests infinity constants and helpers
func TestInfinityValues(t *testing.T) {
    // Test positive infinity detection
    if !IsPosInfByte(INF_BYTE) {
        t.Error("INF_BYTE not detected as positive infinity")
    }
    if !IsPosInfShort(INF_SHORT) {
        t.Error("INF_SHORT not detected as positive infinity")
    }
    if !IsPosInfInt(INF_INT) {
        t.Error("INF_INT not detected as positive infinity")
    }
    if !IsPosInfLong(INF_LONG) {
        t.Error("INF_LONG not detected as positive infinity")
    }
    if !IsPosInfReal(INF_REAL) {
        t.Error("INF_REAL not detected as positive infinity")
    }
    if !IsPosInfFloat(INF_FLOAT) {
        t.Error("INF_FLOAT not detected as positive infinity")
    }

    // Test negative infinity detection
    if !IsNegInfByte(NEG_INF_BYTE) {
        t.Error("NEG_INF_BYTE not detected as negative infinity")
    }
    if !IsNegInfShort(NEG_INF_SHORT) {
        t.Error("NEG_INF_SHORT not detected as negative infinity")
    }
    if !IsNegInfInt(NEG_INF_INT) {
        t.Error("NEG_INF_INT not detected as negative infinity")
    }
    if !IsNegInfLong(NEG_INF_LONG) {
        t.Error("NEG_INF_LONG not detected as negative infinity")
    }
    if !IsNegInfReal(NEG_INF_REAL) {
        t.Error("NEG_INF_REAL not detected as negative infinity")
    }
    if !IsNegInfFloat(NEG_INF_FLOAT) {
        t.Error("NEG_INF_FLOAT not detected as negative infinity")
    }

    // Test general infinity detection (both positive and negative)
    if !IsInfLong(INF_LONG) {
        t.Error("INF_LONG not detected as infinity")
    }
    if !IsInfLong(NEG_INF_LONG) {
        t.Error("NEG_INF_LONG not detected as infinity")
    }

    // Test non-infinity values
    if IsInfInt(0) {
        t.Error("0 incorrectly detected as infinity")
    }
    if IsInfLong(100) {
        t.Error("100 incorrectly detected as infinity")
    }
}

// TestNullSerialization tests that NULL values serialize and deserialize correctly
func TestNullSerialization(t *testing.T) {
    tests := []struct {
        name string
        data TPTypes
    }{
        {"NullByte", NewTPNullByte()},
        {"NullShort", NewTPNullShort()},
        {"NullInt", NewTPNullInt()},
        {"NullLong", NewTPNullLong()},
        {"NullReal", NewTPNullReal()},
        {"NullFloat", NewTPNullFloat()},
        {"NullEnum", NewTPNullEnum()},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Serialize
            binary, err := TPDataSer(tt.data, false)
            if err != nil {
                t.Fatalf("Failed to serialize %s: %v", tt.name, err)
            }

            // Deserialize
            result, err := TPDataDe(binary)
            if err != nil {
                t.Fatalf("Failed to deserialize %s: %v", tt.name, err)
            }

            // Check type matches
            if result.Tp != tt.data.Tp {
                t.Errorf("Type mismatch for %s: got 0x%x, want 0x%x", tt.name, result.Tp, tt.data.Tp)
            }

            // Verify NULL values are preserved
            switch result.Tp {
            case SC_BYTE:
                if !IsNullByte(int8(result.Data.(TPByte))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            case SC_SHORT:
                if !IsNullShort(int16(result.Data.(TPShort))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            case SC_INT:
                if !IsNullInt(int32(result.Data.(TPInt))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            case SC_LONG:
                if !IsNullLong(int64(result.Data.(TPLong))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            case SC_REAL:
                if !IsNullReal(float32(result.Data.(TPReal))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            case SC_FLOAT:
                if !IsNullFloat(float64(result.Data.(TPFloat))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            case SC_ENUM:
                if !IsNullEnum(uint32(result.Data.(TPEnum))) {
                    t.Errorf("%s: NULL value not preserved", tt.name)
                }
            }
        })
    }
}

// TestInfinitySerialization tests that infinity values serialize and deserialize correctly
func TestInfinitySerialization(t *testing.T) {
    tests := []struct {
        name string
        data TPTypes
        isPosInf bool
    }{
        {"PosInfByte", NewTPInfByte(), true},
        {"NegInfByte", NewTPNegInfByte(), false},
        {"PosInfShort", NewTPInfShort(), true},
        {"NegInfShort", NewTPNegInfShort(), false},
        {"PosInfInt", NewTPInfInt(), true},
        {"NegInfInt", NewTPNegInfInt(), false},
        {"PosInfLong", NewTPInfLong(), true},
        {"NegInfLong", NewTPNegInfLong(), false},
        {"PosInfReal", NewTPInfReal(), true},
        {"NegInfReal", NewTPNegInfReal(), false},
        {"PosInfFloat", NewTPInfFloat(), true},
        {"NegInfFloat", NewTPNegInfFloat(), false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Serialize
            binary, err := TPDataSer(tt.data, false)
            if err != nil {
                t.Fatalf("Failed to serialize %s: %v", tt.name, err)
            }

            // Deserialize
            result, err := TPDataDe(binary)
            if err != nil {
                t.Fatalf("Failed to deserialize %s: %v", tt.name, err)
            }

            // Check type matches
            if result.Tp != tt.data.Tp {
                t.Errorf("Type mismatch for %s: got 0x%x, want 0x%x", tt.name, result.Tp, tt.data.Tp)
            }

            // Verify infinity values are preserved
            switch result.Tp {
            case SC_BYTE:
                v := int8(result.Data.(TPByte))
                if tt.isPosInf && !IsPosInfByte(v) {
                    t.Errorf("%s: positive infinity not preserved", tt.name)
                }
                if !tt.isPosInf && !IsNegInfByte(v) {
                    t.Errorf("%s: negative infinity not preserved", tt.name)
                }
            case SC_SHORT:
                v := int16(result.Data.(TPShort))
                if tt.isPosInf && !IsPosInfShort(v) {
                    t.Errorf("%s: positive infinity not preserved", tt.name)
                }
                if !tt.isPosInf && !IsNegInfShort(v) {
                    t.Errorf("%s: negative infinity not preserved", tt.name)
                }
            case SC_INT:
                v := int32(result.Data.(TPInt))
                if tt.isPosInf && !IsPosInfInt(v) {
                    t.Errorf("%s: positive infinity not preserved", tt.name)
                }
                if !tt.isPosInf && !IsNegInfInt(v) {
                    t.Errorf("%s: negative infinity not preserved", tt.name)
                }
            case SC_LONG:
                v := int64(result.Data.(TPLong))
                if tt.isPosInf && !IsPosInfLong(v) {
                    t.Errorf("%s: positive infinity not preserved", tt.name)
                }
                if !tt.isPosInf && !IsNegInfLong(v) {
                    t.Errorf("%s: negative infinity not preserved", tt.name)
                }
            case SC_REAL:
                v := float32(result.Data.(TPReal))
                if tt.isPosInf && !IsPosInfReal(v) {
                    t.Errorf("%s: positive infinity not preserved", tt.name)
                }
                if !tt.isPosInf && !IsNegInfReal(v) {
                    t.Errorf("%s: negative infinity not preserved", tt.name)
                }
            case SC_FLOAT:
                v := float64(result.Data.(TPFloat))
                if tt.isPosInf && !IsPosInfFloat(v) {
                    t.Errorf("%s: positive infinity not preserved", tt.name)
                }
                if !tt.isPosInf && !IsNegInfFloat(v) {
                    t.Errorf("%s: negative infinity not preserved", tt.name)
                }
            }
        })
    }
}

// TestVectorWithSpecialValues tests vectors containing NULL and infinity values
func TestVectorWithSpecialValues(t *testing.T) {
    // Vector with NULLs and infinities
    vecLong := NewTPVecLong([]int64{
        100,
        NULL_LONG,
        INF_LONG,
        NEG_INF_LONG,
        0,
        -500,
    })

    binary, err := TPDataSer(vecLong, false)
    if err != nil {
        t.Fatalf("Failed to serialize vector with special values: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize vector with special values: %v", err)
    }

    resultVec := result.Data.(TPVecLong)
    if len(resultVec) != 6 {
        t.Errorf("Vector length mismatch: got %d, want 6", len(resultVec))
    }

    // Check special values preserved
    if !IsNullLong(int64(resultVec[1])) {
        t.Error("NULL value not preserved in vector")
    }
    if !IsPosInfLong(int64(resultVec[2])) {
        t.Error("Positive infinity not preserved in vector")
    }
    if !IsNegInfLong(int64(resultVec[3])) {
        t.Error("Negative infinity not preserved in vector")
    }

    // Check regular values
    if resultVec[0] != 100 {
        t.Errorf("Regular value not preserved: got %d, want 100", resultVec[0])
    }
}

// TestFloatSpecialValues tests float special values (NaN, Inf)
func TestFloatSpecialValues(t *testing.T) {
    vecFloat := NewTPVecFloat([]float64{
        1.5,
        NULL_FLOAT,
        INF_FLOAT,
        NEG_INF_FLOAT,
        0.0,
        -2.5,
    })

    binary, err := TPDataSer(vecFloat, false)
    if err != nil {
        t.Fatalf("Failed to serialize float vector: %v", err)
    }

    result, err := TPDataDe(binary)
    if err != nil {
        t.Fatalf("Failed to deserialize float vector: %v", err)
    }

    resultVec := result.Data.(TPVecFloat)
    if len(resultVec) != 6 {
        t.Errorf("Vector length mismatch: got %d, want 6", len(resultVec))
    }

    // Check special values preserved
    if !IsNullFloat(float64(resultVec[1])) {
        t.Error("NULL float not preserved in vector")
    }
    if !IsPosInfFloat(float64(resultVec[2])) {
        t.Error("Positive infinity float not preserved in vector")
    }
    if !IsNegInfFloat(float64(resultVec[3])) {
        t.Error("Negative infinity float not preserved in vector")
    }
}
