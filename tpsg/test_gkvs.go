package main

import (
    "fmt"
    "time"
)

func TestGKVS() {
    LogEvent("Starting GKVS concurrent access test")

    // Goroutine 1: Set uint16 value 1 to key "test1"
    go func() {
        TConfig.Set("test1", NewGKVSUInt16(1))
        LogEvent("Goroutine 1: Set test1 = 1")
    }()

    // Goroutine 2: Wait 5 seconds, retrieve test1, print it, then set value 2
    go func() {
        time.Sleep(5 * time.Second)
        value := TConfig.Get("test1").UInt16
        fmt.Printf("Goroutine 2: Retrieved test1 = %d\n", value)
        LogEvent(fmt.Sprintf("Goroutine 2: Retrieved test1 = %d", value))

        TConfig.Set("test1", NewGKVSUInt16(2))
        LogEvent("Goroutine 2: Set test1 = 2")
    }()

    // Goroutine 3: Wait 10 seconds, retrieve test1, print it
    go func() {
        time.Sleep(10 * time.Second)
        value := TConfig.Get("test1").UInt16
        fmt.Printf("Goroutine 3: Retrieved test1 = %d\n", value)
        LogEvent(fmt.Sprintf("Goroutine 3: Retrieved test1 = %d", value))
    }()

    // Wait for all goroutines to complete (longest is 10 seconds + execution time)
    time.Sleep(12 * time.Second)

    LogEvent("GKVS concurrent access test completed")
}
