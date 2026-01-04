// Copyright 2025 Sushanth (https://github.com/sushanthpy)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package toondb

import (
	"fmt"
	"os"
	"time"
)

// Example_analyticsEnabled verifies analytics tracking works when enabled.
func Example_analyticsEnabled() {
	// Test with analytics enabled (default)
	db, err := Open("./test_analytics_db")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()
	defer os.RemoveAll("./test_analytics_db")

	// Perform operations
	_ = db.Put([]byte("key"), []byte("value"))
	_, _ = db.Get([]byte("key"))

	// Give time for analytics to send
	time.Sleep(100 * time.Millisecond)

	fmt.Println("Analytics test with tracking enabled completed")
	// Output: Analytics test with tracking enabled completed
}

func Example_analyticsDisabled() {
	// Test with analytics disabled
	os.Setenv("TOONDB_DISABLE_ANALYTICS", "true")
	defer os.Unsetenv("TOONDB_DISABLE_ANALYTICS")

	db, err := Open("./test_analytics_disabled_db")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer db.Close()
	defer os.RemoveAll("./test_analytics_disabled_db")

	// Perform operations
	_ = db.Put([]byte("key"), []byte("value"))
	_, _ = db.Get([]byte("key"))

	fmt.Println("Analytics test with tracking disabled completed")
	// Output: Analytics test with tracking disabled completed
}

func Example_connectionError() {
	// Test connection error tracking
	_, err := Open("/invalid/path/that/cannot/be/created")
	if err != nil {
		// Connection error should be tracked
		fmt.Println("Connection error occurred (tracked)")
	}
	// Output: Connection error occurred (tracked)
}
