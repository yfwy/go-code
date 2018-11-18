package kvdb

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"path"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelDB struct {
	db *leveldb.DB
}

func NewLevelDB(name string) (DB, error) {
	dbPath := path.Join(name)
	// Calculate the cache and file descriptor allowance for this particular database
	cache := 64
	handles := 64
	log.Printf("Allotted %dMB cache and %d file handles to %s\n", cache, handles, dbPath)

	// Open the db and recover any potential corruptions
	db, err := leveldb.OpenFile(dbPath, &opt.Options{
		OpenFilesCacheCapacity: handles,
		BlockCacheCapacity:     cache / 2 * opt.MiB,
		WriteBuffer:            cache / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(dbPath, nil)
	}
	// (Re)check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	database := &LevelDB{db: db}
	return database, nil
}

func (db *LevelDB) Get(key []byte) []byte {
	res, err := db.db.Get(key, nil)
	if err != nil {
		if err == errors.ErrNotFound {
			return nil
		} else {
			panic(err)
		}
	}
	return res
}

func (db *LevelDB) Set(key []byte, value []byte) {
	err := db.db.Put(key, value, nil)
	if err != nil {
		panic(err)
	}
}

func (db *LevelDB) SetSync(key []byte, value []byte) {
	err := db.db.Put(key, value, &opt.WriteOptions{Sync: true})
	if err != nil {
		panic(err)
	}
}

func (db *LevelDB) Delete(key []byte) {
	err := db.db.Delete(key, nil)
	if err != nil {
		panic(err)
	}
}

func (db *LevelDB) DeleteSync(key []byte) {
	err := db.db.Delete(key, &opt.WriteOptions{Sync: true})
	if err != nil {
		panic(err)
	}
}

func (db *LevelDB) DB() *leveldb.DB {
	return db.db
}

func (db *LevelDB) Close() {
	db.db.Close()
}

func (db *LevelDB) Print() {
	iter := db.db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

type ldbBatch struct {
	db *leveldb.DB
	b  *leveldb.Batch
	wo *opt.WriteOptions
}

//因为是金融系统，这个选项还是有必要，避免发生丢数据的情况
//而我们设计整个系统尽量采用批量提交，所以性能应该还是可以的
func (d *LevelDB) NewBatch(sync bool) IBatch {
	wop := &opt.WriteOptions{Sync: true}
	batch := new(leveldb.Batch)
	return &ldbBatch{db: d.db, b: batch, wo: wop}
}

func (b *ldbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	return nil
}

func (b *ldbBatch) Delete(key []byte) error {
	b.b.Delete(key)
	return nil
}

func (b *ldbBatch) Commit() error {
	return b.db.Write(b.b, b.wo)
}

func (db *LevelDB) Query(key string, beg, end, bitSize int64) IQuery {
	iter := db.db.NewIterator(nil, nil)
	fmtpad := "20"
	if bitSize == 32 {
		fmtpad = "10"
	}
	minkey := fmt.Sprintf("%s%0"+fmtpad+"d", key, beg)
	maxkey := fmt.Sprintf("%s%0"+fmtpad+"d", key, end)
	return &ldbQuery{iter, fmtpad, key, []byte(minkey), []byte(maxkey), false, true}
}

func (db *LevelDB) GetIter(key string, bitSize int64) IQuery {
	return db.Query(key, 0, math.MaxInt64, bitSize)
}

type ldbQuery struct {
	iter   iterator.Iterator
	fmtpad string
	key    string
	beg    []byte
	end    []byte
	isseek bool
	first  bool
}

func (q *ldbQuery) Seek(i int64) IQuery {
	//i = 0 表示最后一个
	seekprev := false
	if i == 0 { //seek to end
		i = math.MaxInt64
		seekprev = true
	}
	seekkey := fmt.Sprintf("%s%0"+q.fmtpad+"d", q.key, i)
	//log.Println("seek:", seekkey)
	q.iter.Seek([]byte(seekkey))
	if seekprev {
		q.iter.Prev()
	}
	q.isseek = true
	return q
}

func (q *ldbQuery) FetchAll(limit int64, asc bool) (keys [][]byte, vals [][]byte, err error) {
	iter := q.iter
	defer iter.Release()
	i := 0
	if asc {
		for q.Next() {
			key := iter.Key()
			val := iter.Value()
			//log.Println("fetch.asc ", string(key),val)
			keys = append(keys, []byte(string(key)))
			vals = append(vals, []byte(string(val)))
			i++
			if limit > 0 && i >= int(limit) {
				break
			}
		}
	} else {
		for q.Prev() {
			key := iter.Key()
			val := iter.Value()
			//log.Println("fetch.desc ", string(key))
			keys = append(keys, []byte(string(key)))
			vals = append(vals, []byte(string(val)))
			i++
			if limit > 0 && i >= int(limit) {
				break
			}
		}
	}
	if err := q.Error(); err != nil {
		return nil, nil, err
	}
	return keys, vals, nil
}

func (q *ldbQuery) Next() bool {
	if !q.isseek {
		//log.Println("next seek to:", string(q.beg))
		q.iter.Seek(q.beg)
		q.isseek = true
	}
	if !q.first {
		q.iter.Next()
	}
	q.first = false
	end := q.iter.Key()
	//log.Println("next", toStr(end), q.key)
	if !bytes.HasPrefix(end, []byte(q.key)) || bytes.Compare(end, q.end) > 0 {
		//log.Println("next return false")
		return false
	}
	return q.iter.Valid()
}

func (q *ldbQuery) Key() []byte {
	return q.iter.Key()
}

func (q *ldbQuery) Close() {
	q.iter.Release()
}

func (q *ldbQuery) Error() error {
	return q.iter.Error()
}

func (q *ldbQuery) Value() []byte {
	return q.iter.Value()
}

func toStr(b []byte) string {
	for i := 0; i < len(b); i++ {
		if b[i] == byte(7) {
			b = b[:i]
			b = append(b, []byte("<end>")...)
			break
		}
	}
	return string(b)
}

func (q *ldbQuery) Prev() bool {
	if !q.isseek {
		q.iter.Seek(q.end)
		q.isseek = true
	}
	if !q.first {
		q.iter.Prev()
	}
	q.first = false
	end := q.iter.Key()
	//log.Println("prev", toStr(end))
	if !bytes.HasPrefix(end, []byte(q.key)) || bytes.Compare(q.beg, end) > 0 {
		//log.Println("prev return false")
		return false
	}
	return q.iter.Valid()
}

var dbDrivers = make(map[string]func(name string) (DB, error))

func init() {
	dbDrivers["leveldb"] = NewLevelDB
}
