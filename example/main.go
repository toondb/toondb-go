package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sochdb/sochdb-go/embedded"
)

func main() {
	fmt.Println("SochDB Go SDK v0.4.0 - Embedded Mode Example")
	fmt.Println("=============================================")
	fmt.Println()

	dbPath := "./example_db"

	// Clean up old database
	os.RemoveAll(dbPath)

	// Open database in embedded mode (no server needed)
	fmt.Println("Opening database in embedded mode...")
	db, err := embedded.Open(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	fmt.Println("✅ Database opened")
	fmt.Println()

	// Example 1: Basic Key-Value Operations
	fmt.Println("1. Basic Key-Value Operations:")
	basicKVExample(db)
	fmt.Println()

	// Example 2: Path-Based Keys
	fmt.Println("2. Path-Based Operations:")
	pathExample(db)
	fmt.Println()

	// Example 3: Transactions
	fmt.Println("3. ACID Transactions:")
	transactionExample(db)
	fmt.Println()

	// Example 4: Prefix Scanning
	fmt.Println("4. Prefix Scanning:")
	scanExample(db)
	fmt.Println()

	// Example 5: Statistics
	fmt.Println("5. Database Statistics:")
	statsExample(db)
	fmt.Println()

	fmt.Println("✅ All examples completed successfully!")
}

func basicKVExample(db *embedded.Database) {
	// Put key-value pairs
	err := db.Put([]byte("user:1"), []byte(`{"name": "Alice", "age": 30}`))
	if err != nil {
		log.Printf("Put failed: %v", err)
		return
	}
	fmt.Println("  ✅ Put: user:1 -> {name: Alice, age: 30}")

	err = db.Put([]byte("user:2"), []byte(`{"name": "Bob", "age": 25}`))
	if err != nil {
		log.Printf("Put failed: %v", err)
		return
	}
	fmt.Println("  ✅ Put: user:2 -> {name: Bob, age: 25}")

	// Get key-value
	value, err := db.Get([]byte("user:1"))
	if err != nil {
		log.Printf("Get failed: %v", err)
		return
	}
	fmt.Printf("  ✅ Get: user:1 -> %s\n", string(value))

	// Delete key
	err = db.Delete([]byte("user:2"))
	if err != nil {
		log.Printf("Delete failed: %v", err)
		return
	}
	fmt.Println("  ✅ Delete: user:2")
}

func pathExample(db *embedded.Database) {
	// Path-based keys for hierarchical data
	err := db.PutPath("users/alice/email", []byte("alice@example.com"))
	if err != nil {
		log.Printf("PutPath failed: %v", err)
		return
	}
	fmt.Println("  ✅ Put: users/alice/email -> alice@example.com")

	err = db.PutPath("users/alice/age", []byte("30"))
	if err != nil {
		log.Printf("PutPath failed: %v", err)
		return
	}
	fmt.Println("  ✅ Put: users/alice/age -> 30")

	// Get path-based value
	email, err := db.GetPath("users/alice/email")
	if err != nil {
		log.Printf("GetPath failed: %v", err)
		return
	}
	fmt.Printf("  ✅ Get: users/alice/email -> %s\n", string(email))
}

func transactionExample(db *embedded.Database) {
	// Using WithTransaction for automatic commit/rollback
	err := db.WithTransaction(func(txn *embedded.Transaction) error {
		err := txn.Put([]byte("counter"), []byte("0"))
		if err != nil {
			return err
		}
		err = txn.Put([]byte("timestamp"), []byte("2026-01-12"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Printf("Transaction failed: %v", err)
		return
	}
	fmt.Println("  ✅ Transaction committed (2 writes)")

	// Manual transaction control
	txn := db.Begin()
	err = txn.Put([]byte("item:1"), []byte("value1"))
	if err != nil {
		txn.Abort()
		log.Printf("Transaction failed: %v", err)
		return
	}
	err = txn.Commit()
	if err != nil {
		log.Printf("Commit failed: %v", err)
		return
	}
	fmt.Println("  ✅ Manual transaction committed")
}

func scanExample(db *embedded.Database) {
	// Insert test data
	for i := 1; i <= 5; i++ {
		key := fmt.Sprintf("product:%d", i)
		value := fmt.Sprintf("Item %d", i)
		err := db.Put([]byte(key), []byte(value))
		if err != nil {
			log.Printf("Put failed: %v", err)
			return
		}
	}
	fmt.Println("  ✅ Inserted 5 products")

	// Scan with prefix using transaction
	txn := db.Begin()
	defer txn.Abort()

	iter := txn.ScanPrefix([]byte("product:"))
	defer iter.Close()

	count := 0
	for {
		key, value, ok := iter.Next()
		if !ok {
			break
		}
		count++
		if count <= 3 { // Show first 3
			fmt.Printf("    %s -> %s\n", string(key), string(value))
		}
	}

	if err := iter.Err(); err != nil {
		log.Printf("Scan failed: %v", err)
		return
	}
	fmt.Printf("  ✅ Scanned %d items\n", count)
}

func statsExample(db *embedded.Database) {
	stats, err := db.Stats()
	if err != nil {
		log.Printf("Stats failed: %v", err)
		return
	}
	fmt.Printf("  📊 Active transactions: %d\n", stats.ActiveTransactions)
	fmt.Printf("  📊 Memtable size: %d bytes\n", stats.MemtableSizeBytes)
	fmt.Printf("  📊 WAL size: %d bytes\n", stats.WalSizeBytes)
	fmt.Printf("  📊 Last checkpoint LSN: %d\n", stats.LastCheckpointLsn)
}
