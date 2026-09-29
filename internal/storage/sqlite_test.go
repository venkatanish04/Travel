package storage

import (
	"path/filepath"
	"testing"
)

func TestOpenDatabaseInitializesSchema(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "travelraft.db")
	db, err := OpenDatabase(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := InitializeSchema(db); err != nil {
		t.Fatal(err)
	}

	var tableCount int
	if err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name IN ('users', 'vehicles', 'seats', 'bookings')").Scan(&tableCount); err != nil {
		t.Fatal(err)
	}
	if tableCount != 4 {
		t.Fatalf("schema table count = %d, want 4", tableCount)
	}
}
