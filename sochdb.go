// Package sochdb provides Go SDK for SochDB v0.4.0
//
// Dual-mode architecture: Embedded (FFI) + Server (gRPC/IPC)
//
// Architecture: Flexible Deployment
// ==================================
// This SDK supports BOTH modes:
//
// 1. Embedded Mode (FFI) - For single-process apps:
//   - Direct FFI bindings to Rust libraries (via CGO)
//   - No server required - just import and run
//   - Best for: Local development, simple apps, edge deployments
//
// 2. Server Mode (gRPC) - For distributed systems:
//   - Thin client connecting to sochdb-grpc server
//   - Best for: Production, multi-language, scalability
//
// Example (Server Mode - gRPC):
//
//	import "github.com/sochdb/sochdb-go"
//
//	// Connect to server
//	client, err := sochdb.GrpcConnect("localhost:50051")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Use the client
//	err = client.PutKv(context.Background(), "namespace", "key", []byte("value"))
//
// Example (Embedded Mode - IPC):
//
//	import "github.com/sochdb/sochdb-go"
//
//	// Connect via Unix socket
//	client, err := sochdb.Connect("/path/to/db/sochdb.sock")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Use the client
//	value, err := client.Get([]byte("key"))
package sochdb

// Version is the current SDK version.
const Version = "0.3.6"
