package kvdb

import (
	"errors"
	"fmt"
	"path"
)

type IQuery interface {
	Seek(i int64) IQuery
	FetchAll(limit int64, asc bool) (keys [][]byte, vals [][]byte, err error)
	Key() []byte
	Close()
	Error() error
	Value() []byte
	Prev() bool
	Next() bool
}

type IBatch interface {
	Put(key, value []byte) error
	Delete(key []byte) error
	Commit() error
}

type DB interface {
	Get([]byte) []byte
	Set([]byte, []byte)
	SetSync([]byte, []byte)
	Delete([]byte)
	DeleteSync([]byte)
	Close()
	Query(key string, beg, end, bitSize int64) IQuery
	GetIter(key string, bitSize int64) IQuery
	NewBatch(sync bool) IBatch
	// For debugging
	Print()
}

//-----------------------------------------------------------------------------

// Database types
const DBBackendMemDB = "memdb"
const DBBackendLevelDB = "leveldb"
const DBBackendRocksDB = "rocksdb"

func NewDB(name string, backend string, dir string) DB {
	switch backend {
	case DBBackendMemDB:
		db, err := dbDrivers["memdb"]("path")
		if err != nil {
			panic(err)
		}
		return db
	case DBBackendLevelDB:
		db, err := dbDrivers["leveldb"](path.Join(dir, name+".db"))
		if err != nil {
			panic(err)
		}
		return db
	case DBBackendRocksDB:
		db, err := dbDrivers["rocksdb"](path.Join(dir, name+".db"))
		if err != nil {
			panic(err)
		}
		return db
	default:
		panic(errors.New(fmt.Sprintf("Unknown DB backend: %v", backend)))
	}
	return nil
}
