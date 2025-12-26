package tpsg

import (
	"fmt"
	"net"
	"tpsg/tpserde"
)

func TestSeqs() {
	// // Testing sequence: TCP client connection
	// LogEvent("Starting TCP client test...")
	// conn, err := net.Dial("tcp", "127.0.0.1:17001")
	// if err != nil {
	// 	LogError(fmt.Sprintf("Failed to connect to 127.0.0.1:17001: %s", err.Error()))
	// } else {
	// 	LogEvent(fmt.Sprintf("Successfully connected to 127.0.0.1:17001, connection: %v", conn))

	// 	// Perform IPC handshake as active side (client)
	// 	features := tpserde.NewFeatures().WithBuffered()
	// 	handshake, err := tpserde.ExchangeHandshake(conn, features)
	// 	if err != nil {
	// 		LogError(fmt.Sprintf("Handshake failed: %s", err.Error()))
	// 		return
	// 	}
	// 	LogEvent(fmt.Sprintf("Handshake completed, server features: buffered=%v, compressed=%v",
	// 		handshake.Features.IsBuffered(), handshake.Features.IsCompressed()))

	// 	// Wait 1 second
	// 	time.Sleep(1 * time.Second)

	// 	// Create dictionary with symbol keys and string values
	// 	keys := tpserde.NewTPVecSymbol([]string{"msgType", "payload"})
	// 	values := tpserde.NewTPList([]tpserde.TPTypes{
	// 		tpserde.NewTPVecChar("test"),
	// 		tpserde.NewTPVecChar("test string"),
	// 	})
	// 	testDict := tpserde.NewTPDict(keys, values)
	// 	binaryData, err := tpserde.TPDataSer(testDict, true)
	// 	if err != nil {
	// 		LogError(fmt.Sprintf("Failed to serialize data: %s", err.Error()))
	// 	} else {
	// 		LogEvent(fmt.Sprintf("Serialized dictionary to binary, size: %d bytes", len(binaryData)))

	// 		// Send binary data over connection
	// 		n, err := conn.Write(binaryData)
	// 		if err != nil {
	// 			LogError(fmt.Sprintf("Failed to send data: %s", err.Error()))
	// 		} else {
	// 			LogEvent(fmt.Sprintf("Sent %d bytes over connection", n))
	// 		}
	// 	}
	// }

	// Testing sequence: TCP client connection with table in payload
	LogEvent("Starting TCP client test with table payload...")
	conn2, err := net.Dial("tcp", "127.0.0.1:17001")
	if err != nil {
		LogError(fmt.Sprintf("Failed to connect to 127.0.0.1:17001: %s", err.Error()))
	} else {
		LogEvent(fmt.Sprintf("Successfully connected to 127.0.0.1:17001, connection: %v", conn2))

		// Perform IPC handshake as active side (client)
		features := tpserde.NewFeatures().WithBuffered()
		handshake, err := tpserde.ExchangeHandshake(conn2, features)
		if err != nil {
			LogError(fmt.Sprintf("Handshake failed: %s", err.Error()))
			return
		}
		LogEvent(fmt.Sprintf("Handshake completed, server features: buffered=%v, compressed=%v",
			handshake.Features.IsBuffered(), handshake.Features.IsCompressed()))

		// // Wait 1 second
		// time.Sleep(1 * time.Second)

		// Create table with columns: x (long vec), y (symbol vec), z (list of strings)
		tableKeys := tpserde.NewTPVecSymbol([]string{"x", "y", "z"})
		tableValues := tpserde.NewTPList([]tpserde.TPTypes{
			tpserde.NewTPVecLong([]int64{1, 2, 3}),
			tpserde.NewTPVecSymbol([]string{"y1", "y2", "y3"}),
			tpserde.NewTPList([]tpserde.TPTypes{
				tpserde.NewTPVecChar("z1"),
				tpserde.NewTPVecChar("z2"),
				tpserde.NewTPVecChar("z3"),
			}),
		})
		table := tpserde.NewTPTable(tableKeys, tableValues)

		// Create dictionary with symbol keys: msgType and payload (table)
		dictKeys := tpserde.NewTPVecSymbol([]string{"msgType", "payload"})
		dictValues := tpserde.NewTPList([]tpserde.TPTypes{
			tpserde.NewTPVecChar("test"),
			table,
		})
		testDict := tpserde.NewTPDict(dictKeys, dictValues)

		binaryData, err := tpserde.TPDataSer(testDict, true)
		if err != nil {
			LogError(fmt.Sprintf("Failed to serialize data: %s", err.Error()))
		} else {
			LogEvent(fmt.Sprintf("Serialized dictionary with table to binary, size: %d bytes", len(binaryData)))

			// Send binary data over connection
			n, err := conn2.Write(binaryData)
			if err != nil {
				LogError(fmt.Sprintf("Failed to send data: %s", err.Error()))
			} else {
				LogEvent(fmt.Sprintf("Sent %d bytes over connection", n))
			}
		}
	}
}
