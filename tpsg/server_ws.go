package tpsg

import (
    "fmt"
    "net"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func RunWSServer(port uint16) {
    go func() {
        addr := fmt.Sprintf(":%d", port)

        listener, err := net.Listen("tcp", addr)
        if err != nil {
            LogError(fmt.Sprintf("Failed to start WebSocket server on port %d: %v", port, err))
            return
        }
        defer listener.Close()

        LogEvent(fmt.Sprintf("WebSocket server started on port %d", port))

        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            conn, err := upgrader.Upgrade(w, r, nil)
            if err != nil {
                LogError(fmt.Sprintf("Failed to upgrade connection: %v", err))
                return
            }
            LogEvent(fmt.Sprintf("New WebSocket connection from %s", conn.RemoteAddr().String()))
            go HandleWSConnection(conn)
        })

        err = http.Serve(listener, nil)
        if err != nil {
            LogError(fmt.Sprintf("WebSocket server error: %v", err))
        }
    }()
}

func HandleWSConnection(conn *websocket.Conn) {
    defer conn.Close()
    defer LogEvent(fmt.Sprintf("WebSocket connection closed: %s", conn.RemoteAddr().String()))

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                LogError(fmt.Sprintf("WebSocket read error: %v", err))
            }
            break
        }

        LogEvent(fmt.Sprintf("Received WebSocket message from %s", conn.RemoteAddr().String()))

        go func(mt int, msg []byte) {
            response := ProcessWSRequest(string(msg))
            err := conn.WriteMessage(mt, []byte(response))
            if err != nil {
                LogError(fmt.Sprintf("WebSocket write error: %v", err))
            }
        }(messageType, message)
    }
}

func ProcessWSRequest(request string) string {
    LogEvent(fmt.Sprintf("Processing WebSocket request: %s", request))
    return fmt.Sprintf("Echo: %s", request)
}
