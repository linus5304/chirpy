package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RefreshTokens: map[string]RefreshToken{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStructure := DBStructure{}
	data, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

// func generateUniqueId(dbStructure DBStructure) int {
// 	maxId := 0
// 	for id := range dbStructure.Chirps {
// 		if id > maxId {
// 			maxId = id
// 		}
// 	}
// 	maxId = maxId + 1
// 	return maxId
// }
