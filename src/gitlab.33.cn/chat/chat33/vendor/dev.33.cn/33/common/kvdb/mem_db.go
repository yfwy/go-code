package kvdb

import (
	"fmt"
)

type MemDB struct {
	db map[string][]byte
}

func NewMemDB(path string) (DB, error) {
	database := &MemDB{db: make(map[string][]byte)}
	return database, nil
}

func (db *MemDB) Get(key []byte) []byte {
	return db.db[string(key)]
}

func (db *MemDB) Set(key []byte, value []byte) {
	db.db[string(key)] = value
}

func (db *MemDB) SetSync(key []byte, value []byte) {
	db.db[string(key)] = value
}

func (db *MemDB) Delete(key []byte) {
	delete(db.db, string(key))
}

func (db *MemDB) DeleteSync(key []byte) {
	delete(db.db, string(key))
}

func (db *MemDB) Close() {
	db = nil
}

func (db *MemDB) Print() {
	for key, value := range db.db {
		fmt.Printf("[%X]:\t[%X]\n", []byte(key), value)
	}
}

func (db *MemDB) Query(key string, beg, end, bitSize int64) IQuery {
	return nil
}

func (db *MemDB) GetIter(key string, bitSize int64) IQuery {
	return nil
}

func (db *MemDB) NewBatch(sync bool) IBatch {
	return nil
}

func init() {
	dbDrivers["memdb"] = NewMemDB
}
