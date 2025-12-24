package tpsg

import (
    "sync"
    "testing"
    "time"
)

// TestGKVSBasicOperations tests basic Set, Get, and Delete operations
func TestGKVSBasicOperations(t *testing.T) {
    gkvs := NewGKVS()

    // Test Set and Get
    gkvs.Set("key1", NewGKVSString("value1"))
    result := gkvs.Get("key1")
    if result.Type != GKVSTypeString || result.String != "value1" {
        t.Errorf("Expected String='value1', got Type=%v, String='%s'", result.Type, result.String)
    }

    // Test Get non-existent key
    result = gkvs.Get("nonexistent")
    if result.Type != GKVSTypeNone {
        t.Errorf("Expected Type=GKVSTypeNone for non-existent key, got Type=%v", result.Type)
    }

    // Test Delete
    deleted := gkvs.Delete("key1")
    if deleted.Type != GKVSTypeString || deleted.String != "value1" {
        t.Errorf("Expected deleted value String='value1', got Type=%v, String='%s'", deleted.Type, deleted.String)
    }

    // Verify key is deleted
    result = gkvs.Get("key1")
    if result.Type != GKVSTypeNone {
        t.Errorf("Expected Type=GKVSTypeNone after deletion, got Type=%v", result.Type)
    }

    // Test Delete non-existent key
    deleted = gkvs.Delete("nonexistent")
    if deleted.Type != GKVSTypeNone {
        t.Errorf("Expected Type=GKVSTypeNone when deleting non-existent key, got Type=%v", deleted.Type)
    }
}

// TestGKVSAllTypes tests all supported value types
func TestGKVSAllTypes(t *testing.T) {
    gkvs := NewGKVS()

    // Test Int8
    gkvs.Set("int8", NewGKVSInt8(-42))
    if result := gkvs.Get("int8"); result.Type != GKVSTypeInt8 || result.Int8 != -42 {
        t.Errorf("Int8 test failed")
    }

    // Test UInt8
    gkvs.Set("uint8", NewGKVSUInt8(255))
    if result := gkvs.Get("uint8"); result.Type != GKVSTypeUInt8 || result.UInt8 != 255 {
        t.Errorf("UInt8 test failed")
    }

    // Test Int16
    gkvs.Set("int16", NewGKVSInt16(-1000))
    if result := gkvs.Get("int16"); result.Type != GKVSTypeInt16 || result.Int16 != -1000 {
        t.Errorf("Int16 test failed")
    }

    // Test UInt16
    gkvs.Set("uint16", NewGKVSUInt16(65535))
    if result := gkvs.Get("uint16"); result.Type != GKVSTypeUInt16 || result.UInt16 != 65535 {
        t.Errorf("UInt16 test failed")
    }

    // Test Int32
    gkvs.Set("int32", NewGKVSInt32(-100000))
    if result := gkvs.Get("int32"); result.Type != GKVSTypeInt32 || result.Int32 != -100000 {
        t.Errorf("Int32 test failed")
    }

    // Test UInt32
    gkvs.Set("uint32", NewGKVSUInt32(4294967295))
    if result := gkvs.Get("uint32"); result.Type != GKVSTypeUInt32 || result.UInt32 != 4294967295 {
        t.Errorf("UInt32 test failed")
    }

    // Test Int64
    gkvs.Set("int64", NewGKVSInt64(-9223372036854775807))
    if result := gkvs.Get("int64"); result.Type != GKVSTypeInt64 || result.Int64 != -9223372036854775807 {
        t.Errorf("Int64 test failed")
    }

    // Test UInt64
    gkvs.Set("uint64", NewGKVSUInt64(18446744073709551615))
    if result := gkvs.Get("uint64"); result.Type != GKVSTypeUInt64 || result.UInt64 != 18446744073709551615 {
        t.Errorf("UInt64 test failed")
    }

    // Test Float32
    gkvs.Set("float32", NewGKVSFloat32(3.14159))
    if result := gkvs.Get("float32"); result.Type != GKVSTypeFloat32 || result.Float32 != 3.14159 {
        t.Errorf("Float32 test failed")
    }

    // Test Float64
    gkvs.Set("float64", NewGKVSFloat64(2.718281828459045))
    if result := gkvs.Get("float64"); result.Type != GKVSTypeFloat64 || result.Float64 != 2.718281828459045 {
        t.Errorf("Float64 test failed")
    }

    // Test String
    gkvs.Set("string", NewGKVSString("hello world"))
    if result := gkvs.Get("string"); result.Type != GKVSTypeString || result.String != "hello world" {
        t.Errorf("String test failed")
    }

    // Test TUserCreds
    creds := TUserCreds{Username: "testuser", Password: "testpass"}
    gkvs.Set("creds", NewGKVSTUserCreds(creds))
    if result := gkvs.Get("creds"); result.Type != GKVSTypeTUserCreds || result.TUserCreds.Username != "testuser" {
        t.Errorf("TUserCreds test failed")
    }

    // Test TConfigTOML
    config := TConfigTOML{TCP: 8080, WS: 8081}
    gkvs.Set("config", NewGKVSTConfigTOML(config))
    if result := gkvs.Get("config"); result.Type != GKVSTypeTConfigTOML || result.TConfigTOML.TCP != 8080 {
        t.Errorf("TConfigTOML test failed")
    }

    // Test None
    gkvs.Set("none", NewGKVSNone())
    if result := gkvs.Get("none"); result.Type != GKVSTypeNone {
        t.Errorf("None test failed")
    }
}

// TestGKVSConcurrentAccess tests thread-safe concurrent operations
func TestGKVSConcurrentAccess(t *testing.T) {
    gkvs := NewGKVS()
    var wg sync.WaitGroup

    // Goroutine 1: Set initial value
    wg.Add(1)
    go func() {
        defer wg.Done()
        gkvs.Set("test1", NewGKVSUInt16(1))
    }()

    // Goroutine 2: Wait, then read and update
    wg.Add(1)
    go func() {
        defer wg.Done()
        time.Sleep(100 * time.Millisecond)
        value := gkvs.Get("test1").UInt16
        if value != 1 {
            t.Errorf("Goroutine 2: Expected test1=1, got %d", value)
        }
        gkvs.Set("test1", NewGKVSUInt16(2))
    }()

    // Goroutine 3: Wait longer, then read
    wg.Add(1)
    go func() {
        defer wg.Done()
        time.Sleep(200 * time.Millisecond)
        value := gkvs.Get("test1").UInt16
        if value != 2 {
            t.Errorf("Goroutine 3: Expected test1=2, got %d", value)
        }
    }()

    wg.Wait()
}

// TestGKVSConcurrentStress performs stress testing with many concurrent operations
func TestGKVSConcurrentStress(t *testing.T) {
    gkvs := NewGKVS()
    var wg sync.WaitGroup

    // Spawn 100 goroutines writing different keys
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            key := "key_" + string(rune('0'+id%10))
            gkvs.Set(key, NewGKVSInt32(int32(id)))
        }(i)
    }

    // Spawn 100 goroutines reading keys
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            key := "key_" + string(rune('0'+id%10))
            _ = gkvs.Get(key)
        }(i)
    }

    wg.Wait()
}

// TestGKVSSetOverwrite tests overwriting existing keys
func TestGKVSSetOverwrite(t *testing.T) {
    gkvs := NewGKVS()

    // Set initial value
    gkvs.Set("key", NewGKVSString("initial"))
    if result := gkvs.Get("key"); result.String != "initial" {
        t.Errorf("Expected 'initial', got '%s'", result.String)
    }

    // Overwrite with different type
    gkvs.Set("key", NewGKVSInt32(42))
    result := gkvs.Get("key")
    if result.Type != GKVSTypeInt32 || result.Int32 != 42 {
        t.Errorf("Expected Int32=42, got Type=%v, Int32=%d", result.Type, result.Int32)
    }
}
