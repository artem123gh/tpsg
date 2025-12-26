package tpsg

import (
    "fmt"
    "net"
    "time"
    "tpsg/tpserde"
)

func TestSeqs() {
    // Testing sequence: TCP client connection
    LogEvent("Starting TCP client test...")
    conn, err := net.Dial("tcp", "127.0.0.1:17001")
    if err != nil {
        LogError(fmt.Sprintf("Failed to connect to 127.0.0.1:17001: %s", err.Error()))
    } else {
        LogEvent(fmt.Sprintf("Successfully connected to 127.0.0.1:17001, connection: %v", conn))

        // Wait 1 second
        time.Sleep(1 * time.Second)

        // Create string variable and serialize it
        testString := "test"
        testData := tpserde.NewTPVecChar(testString)
        binaryData, err := tpserde.TPDataSer(testData, false)
        if err != nil {
            LogError(fmt.Sprintf("Failed to serialize data: %s", err.Error()))
        } else {
            LogEvent(fmt.Sprintf("Serialized '%s' to binary, size: %d bytes", testString, len(binaryData)))

            // Send binary data over connection
            n, err := conn.Write(binaryData)
            if err != nil {
                LogError(fmt.Sprintf("Failed to send data: %s", err.Error()))
            } else {
                LogEvent(fmt.Sprintf("Sent %d bytes over connection", n))
            }
        }
    }
}
