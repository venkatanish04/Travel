package storage

import (
	"encoding/json"
	"os"
)

type Database struct{ path string }

func Open(path string) *Database { return &Database{path: path} }
func (d *Database) Save(value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(d.path, data, 0644)
}
func (d *Database) Load(value any) error {
	data, err := os.ReadFile(d.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
