package tpserde

import (
    "bytes"
    "testing"
)

func TestHandshakeSerializeDeserialize(t *testing.T) {
    // Create handshake
    handshake := NewHandshake().WithBuffered()

    // Serialize
    buf := &bytes.Buffer{}
    if err := handshake.Serialize(buf); err != nil {
        t.Fatalf("Failed to serialize handshake: %v", err)
    }

    // Check size (should be 8 bytes: 4 for version + 4 for features)
    if buf.Len() != 8 {
        t.Errorf("Expected handshake size 8 bytes, got %d", buf.Len())
    }

    // Deserialize
    result, err := DeserializeHandshake(buf)
    if err != nil {
        t.Fatalf("Failed to deserialize handshake: %v", err)
    }

    // Verify version
    if result.Version != IPC_VERSION {
        t.Errorf("Expected version %d, got %d", IPC_VERSION, result.Version)
    }

    // Verify features
    if !result.Features.IsBuffered() {
        t.Error("Expected buffered feature to be set")
    }
}

func TestHandshakeVersion(t *testing.T) {
    // Version 0.1.0 should be encoded as: 0 << 20 | 1 << 10 | 0 = 1024
    expectedVersion := uint32(1024)
    if IPC_VERSION != expectedVersion {
        t.Errorf("Expected IPC_VERSION to be %d, got %d", expectedVersion, IPC_VERSION)
    }
}

func TestHandshakeFeatures(t *testing.T) {
    tests := []struct {
        name     string
        builder  func(Handshake) Handshake
        check    func(Handshake) bool
        checkMsg string
    }{
        {
            name:     "Compressed",
            builder:  func(h Handshake) Handshake { return h.WithCompressed() },
            check:    func(h Handshake) bool { return h.Features.IsCompressed() },
            checkMsg: "compressed feature should be set",
        },
        {
            name:     "Buffered",
            builder:  func(h Handshake) Handshake { return h.WithBuffered() },
            check:    func(h Handshake) bool { return h.Features.IsBuffered() },
            checkMsg: "buffered feature should be set",
        },
        {
            name:     "Unsupported",
            builder:  func(h Handshake) Handshake { return h.WithUnsupported() },
            check:    func(h Handshake) bool { return h.IsUnsupported() },
            checkMsg: "unsupported feature should be set",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handshake := tt.builder(NewHandshake())
            if !tt.check(handshake) {
                t.Error(tt.checkMsg)
            }
        })
    }
}

func TestExchangeHandshake(t *testing.T) {
    // Simulate a client-server handshake exchange
    // Create a combined buffer that both sides can read/write
    combinedBuf := &bytes.Buffer{}

    // Client sends handshake with buffered feature
    clientHandshake := NewHandshake().WithBuffered()
    if err := clientHandshake.Serialize(combinedBuf); err != nil {
        t.Fatalf("Client failed to send handshake: %v", err)
    }

    // Server receives handshake
    receivedHandshake, err := DeserializeHandshake(combinedBuf)
    if err != nil {
        t.Fatalf("Server failed to receive handshake: %v", err)
    }

    // Verify server received buffered feature
    if !receivedHandshake.Features.IsBuffered() {
        t.Error("Server should have received buffered feature")
    }

    // Server sends response handshake
    serverHandshake := NewHandshake().WithBuffered()
    if err := serverHandshake.Serialize(combinedBuf); err != nil {
        t.Fatalf("Server failed to send response handshake: %v", err)
    }

    // Client receives response
    responseHandshake, err := DeserializeHandshake(combinedBuf)
    if err != nil {
        t.Fatalf("Client failed to receive response handshake: %v", err)
    }

    // Verify client received buffered feature
    if !responseHandshake.Features.IsBuffered() {
        t.Error("Client should have received buffered feature")
    }
}

func TestExchangeHandshakeHelper(t *testing.T) {
    // Create buffers for bidirectional communication
    clientToServer := &bytes.Buffer{}
    serverToClient := &bytes.Buffer{}

    // Test that handshakes can be created with features
    clientFeatures := NewFeatures().WithBuffered()
    serverFeatures := NewFeatures().WithBuffered()

    clientHandshake := Handshake{Version: IPC_VERSION, Features: clientFeatures}
    serverHandshake := Handshake{Version: IPC_VERSION, Features: serverFeatures}

    // Serialize client handshake
    if err := clientHandshake.Serialize(clientToServer); err != nil {
        t.Fatalf("Failed to serialize client handshake: %v", err)
    }

    // Deserialize on server side
    receivedClientHandshake, err := DeserializeHandshake(clientToServer)
    if err != nil {
        t.Fatalf("Failed to deserialize client handshake: %v", err)
    }

    if !receivedClientHandshake.Features.IsBuffered() {
        t.Error("Server should have received buffered feature from client")
    }

    // Serialize server handshake
    if err := serverHandshake.Serialize(serverToClient); err != nil {
        t.Fatalf("Failed to serialize server handshake: %v", err)
    }

    // Deserialize on client side
    receivedServerHandshake, err := DeserializeHandshake(serverToClient)
    if err != nil {
        t.Fatalf("Failed to deserialize server handshake: %v", err)
    }

    if !receivedServerHandshake.Features.IsBuffered() {
        t.Error("Client should have received buffered feature from server")
    }
}

func TestHandshakeUnsupported(t *testing.T) {
    handshake := NewHandshake().WithUnsupported()

    if !handshake.IsUnsupported() {
        t.Error("Handshake should have unsupported flag set")
    }

    // Serialize and deserialize
    buf := &bytes.Buffer{}
    if err := handshake.Serialize(buf); err != nil {
        t.Fatalf("Failed to serialize handshake: %v", err)
    }

    result, err := DeserializeHandshake(buf)
    if err != nil {
        t.Fatalf("Failed to deserialize handshake: %v", err)
    }

    if !result.IsUnsupported() {
        t.Error("Deserialized handshake should have unsupported flag set")
    }
}
