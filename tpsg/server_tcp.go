package tpsg

import (
    "bufio"
    "fmt"
    "net"
)

// RunTCPServer starts a TCP server on the specified port in a goroutine
func RunTCPServer(port uint16) {
    go func() {
        address := fmt.Sprintf(":%d", port)
        listener, err := net.Listen("tcp", address)
        if err != nil {
            LogError(fmt.Sprintf("Failed to start TCP server on port %d: %s", port, err.Error()))
            return
        }
        defer listener.Close()

        LogEvent(fmt.Sprintf("TCP server started on port %d", port))

        for {
            conn, err := listener.Accept()
            if err != nil {
                LogError(fmt.Sprintf("Failed to accept TCP connection: %s", err.Error()))
                continue
            }

            LogEvent(fmt.Sprintf("TCP connection accepted from %s", conn.RemoteAddr().String()))
            go HandleTCPConnection(conn)
        }
    }()
}

// HandleTCPConnection processes TCP requests synchronously until connection is closed
func HandleTCPConnection(conn net.Conn) {
    defer conn.Close()
    defer LogEvent(fmt.Sprintf("TCP connection closed from %s", conn.RemoteAddr().String()))

    reader := bufio.NewReader(conn)

    for {
        // Read request line by line until newline
        request, err := reader.ReadString('\n')
        if err != nil {
            // Connection closed or error occurred
            if err.Error() != "EOF" {
                LogError(fmt.Sprintf("Error reading from TCP connection %s: %s", conn.RemoteAddr().String(), err.Error()))
            }
            return
        }

        // Process request and get response
        response := ProcessTCPRequest(request)

        // Send response back to client
        _, err = conn.Write([]byte(response))
        if err != nil {
            LogError(fmt.Sprintf("Error writing to TCP connection %s: %s", conn.RemoteAddr().String(), err.Error()))
            return
        }
    }
}

// ProcessTCPRequest is a placeholder function that echoes the request back to the client
func ProcessTCPRequest(request string) string {
    LogEvent(fmt.Sprintf("TCP request received: %s", request))
    // Echo response
    response := fmt.Sprintf("Echo: %s", request)
    return response
}
