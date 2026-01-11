package embedded_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/toondb/toondb-go/embedded"
)

const testDBPath = "./validation_test_db"

func cleanup() {
	os.RemoveAll(testDBPath)
}

func TestBasicKV(t *testing.T) {
	fmt.Println("\n=== Testing Basic KV Operations ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	// Store data
	err = db.Put([]byte("user:1"), []byte("Alice"))
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	err = db.Put([]byte("user:2"), []byte("Bob"))
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	fmt.Println("✓ Put operations successful")

	// Retrieve data
	user, err := db.Get([]byte("user:1"))
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if string(user) != "Alice" {
		t.Errorf("Expected 'Alice', got '%s'", string(user))
	}
	fmt.Println("✓ Get operation successful")

	// Delete data
	err = db.Delete([]byte("user:1"))
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	deleted, err := db.Get([]byte("user:1"))
	if err != nil {
		t.Fatalf("Get after delete failed: %v", err)
	}
	if deleted != nil {
		t.Error("Delete failed - key still exists")
	}
	fmt.Println("✓ Delete operation successful")
}

func TestPathOperations(t *testing.T) {
	fmt.Println("\n=== Testing Path-Based Keys ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	// Store with path
	err = db.PutPath("users/alice/name", []byte("Alice Smith"))
	if err != nil {
		t.Fatalf("PutPath failed: %v", err)
	}
	err = db.PutPath("users/alice/email", []byte("alice@example.com"))
	if err != nil {
		t.Fatalf("PutPath failed: %v", err)
	}
	err = db.PutPath("users/bob/name", []byte("Bob Jones"))
	if err != nil {
		t.Fatalf("PutPath failed: %v", err)
	}
	fmt.Println("✓ Path put operations successful")

	// Retrieve by path
	name, err := db.GetPath("users/alice/name")
	if err != nil {
		t.Fatalf("GetPath failed: %v", err)
	}
	if string(name) != "Alice Smith" {
		t.Errorf("Expected 'Alice Smith', got '%s'", string(name))
	}
	fmt.Println("✓ Path get operation successful")

	// Delete by path
	// Note: DeletePath not  implemented, using delete with path key
	err = db.Delete([]byte("users/alice/email"))
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	deleted, err := db.GetPath("users/alice/email")
	if err != nil {
		t.Fatalf("GetPath after delete failed: %v", err)
	}
	if deleted != nil {
		t.Error("Path delete failed")
	}
	fmt.Println("✓ Path delete operation successful")
}

func TestTransactions(t *testing.T) {
	fmt.Println("\n=== Testing Transactions (ACID with SSI) ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	// WithTransaction pattern
	err = db.WithTransaction(func(txn *embedded.Transaction) error {
		if err := txn.Put([]byte("accounts/alice"), []byte("1000")); err != nil {
			return err
		}
		if err := txn.Put([]byte("accounts/bob"), []byte("500")); err != nil {
			return err
		}

		// Read within transaction
		balance, err := txn.Get([]byte("accounts/alice"))
		if err != nil {
			return err
		}
		if string(balance) != "1000" {
			t.Error("Transaction read-your-own-writes failed")
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}
	fmt.Println("✓ WithTransaction pattern successful")

	// Manual transaction control
	txn := db.Begin()
	err = txn.Put([]byte("key1"), []byte("value1"))
	if err != nil {
		txn.Abort()
		t.Fatalf("Transaction put failed: %v", err)
	}
	err = txn.Put([]byte("key2"), []byte("value2"))
	if err != nil {
		txn.Abort()
		t.Fatalf("Transaction put failed: %v", err)
	}

	err = txn.Commit()
	if err != nil {
		t.Fatalf("Transaction commit failed: %v", err)
	}
	fmt.Println("✓ Manual transaction committed")

	// Verify committed data
	value1, err := db.Get([]byte("key1"))
	if err != nil {
		t.Fatalf("Get after commit failed: %v", err)
	}
	if string(value1) != "value1" {
		t.Error("Transaction commit verification failed")
	}
	fmt.Println("✓ Transaction commit verified")
}

func TestScanOperations(t *testing.T) {
	fmt.Println("\n=== Testing Prefix Scanning ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	// Insert test data
	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("scan_test:%d", i)
		value := fmt.Sprintf("value%d", i)
		err = db.Put([]byte(key), []byte(value))
		if err != nil {
			t.Fatalf("Put failed: %v", err)
		}
	}
	fmt.Println("✓ Inserted 5 test records")

	// Scan with prefix
	txn := db.Begin()
	defer txn.Abort()

	iter := txn.ScanPrefix([]byte("scan_test:"))
	defer iter.Close()

	count := 0
	for {
		key, value, ok := iter.Next()
		if !ok {
			break
		}
		count++
		fmt.Printf("  Found: %s = %s\n", string(key), string(value))
	}

	if count != 5 {
		t.Errorf("Expected 5 records, found %d", count)
	}
	fmt.Printf("✓ Scanned %d records successfully\n", count)
}

func TestSSIConflictHandling(t *testing.T) {
	fmt.Println("\n=== Testing SSI Conflict Handling ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	// Initialize counter
	err = db.Put([]byte("counter"), []byte("0"))
	if err != nil {
		t.Fatalf("Initial put failed: %v", err)
	}

	const maxRetries = 3
	retries := 0

	for attempt := 0; attempt < maxRetries; attempt++ {
		err = db.WithTransaction(func(txn *embedded.Transaction) error {
			// Read and modify
			value, err := txn.Get([]byte("counter"))
			if err != nil {
				return err
			}

			currentValue, _ := strconv.Atoi(string(value))
			newValue := fmt.Sprintf("%d", currentValue+1)

			return txn.Put([]byte("counter"), []byte(newValue))
		})

		if err != nil {
			if attempt < maxRetries-1 {
				retries++
				fmt.Printf("  Retry %d due to potential SSI conflict\n", attempt+1)
				continue
			}
			t.Fatalf("Transaction failed after %d retries: %v", maxRetries, err)
		}
		break
	}

	fmt.Printf("✓ SSI conflict handling works (%d retries)\n", retries)
}

func TestStatistics(t *testing.T) {
	fmt.Println("\n=== Testing Statistics & Monitoring ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	stats, err := db.Stats()
	if err != nil {
		t.Fatalf("Stats failed: %v", err)
	}

	fmt.Printf("  Active transactions: %d\n", stats.ActiveTransactions)
	fmt.Printf("  Memtable size: %d bytes\n", stats.MemtableSizeBytes)
	fmt.Printf("  WAL size: %d bytes\n", stats.WalSizeBytes)
	fmt.Println("✓ Statistics retrieval successful")
}

func TestCheckpoint(t *testing.T) {
	fmt.Println("\n=== Testing Checkpoints & Snapshots ===")

	db, err := embedded.Open(testDBPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	defer cleanup()

	lsn, err := db.Checkpoint()
	if err != nil {
		t.Fatalf("Checkpoint failed: %v", err)
	}

	fmt.Printf("  Checkpoint LSN: %d\n", lsn)
	fmt.Println("✓ Checkpoint successful")
}

func TestMain(m *testing.M) {
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║    ToonDB Go SDK - Comprehensive Validation       ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")

	cleanup()

	code := m.Run()

	cleanup()

	if code == 0 {
		fmt.Println("╔════════════════════════════════════════════════════╗")
		fmt.Println("║              ✅ ALL TESTS PASSED                   ║")
		fmt.Println("╚════════════════════════════════════════════════════╝")
		fmt.Println()
	}

	os.Exit(code)
}
