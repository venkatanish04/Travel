package storage

import "testing"

func TestMySQLConfigFromEnvUsesLocalDefaults(t *testing.T) {
	config := MySQLConfigFromEnv()
	if config.Username != "root" {
		t.Fatalf("Username = %q, want %q", config.Username, "root")
	}
	if config.Password != "root" {
		t.Fatalf("Password = %q, want %q", config.Password, "root")
	}
	if config.Host != "127.0.0.1" {
		t.Fatalf("Host = %q, want %q", config.Host, "127.0.0.1")
	}
	if config.Port != "3306" {
		t.Fatalf("Port = %q, want %q", config.Port, "3306")
	}
	if config.Database != "travelraft" {
		t.Fatalf("Database = %q, want %q", config.Database, "travelraft")
	}
}
