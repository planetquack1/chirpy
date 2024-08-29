package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	fmt.Println("Ensuring database exists")
	if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Database does not exist")
		initialData := DBStructure{
			Chirps: make(map[int]Chirp),
		}

		file, err := json.Marshal(initialData)
		if err != nil {
			return err
		}

		os.WriteFile("database.json", file, 0666)
		fmt.Println("Wrote to file " + db.path)
		return nil
	}

	fmt.Println("Database exists: " + db.path)
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	fmt.Println("Loading database from file")
	file, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	var data DBStructure
	if err := json.Unmarshal(file, &data); err != nil {
		return DBStructure{}, err
	}

	return data, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(data DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	fmt.Println("Writing database to file")
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, file, 0666)
}
