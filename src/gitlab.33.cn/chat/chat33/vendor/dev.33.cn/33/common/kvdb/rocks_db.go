// +build rocksdb

package kvdb

import (
	"bytes"
	"fmt"
	"github.com/tecbot/gorocksdb"
	. "github.com/tendermint/go-common"
	"log"
	"math"
	"path"
	"runtime"
)

const OpenFileLimit = 64

type RocksDB struct {
	db *gorocksdb.DB
	ro *gorocksdb.ReadOptions
	wo *gorocksdb.WriteOptions
}

func NewRocksDB(name string) (DB, error) {
	dbPath := path.Join(name)
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(128 * 1024 * 1024))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	cpu := runtime.NumCPU() / 2
	if cpu == 0 {
		cpu = 1
	}
	opts.IncreaseParallelism(cpu)
	opts.SetMaxOpenFiles(OpenFileLimit)
	opts.SetWriteBufferSize(32 * 1024 * 1024)
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, dbPath)
	if err != nil {
		return nil, err
	}
	database := &RocksDB{db: db}
	database.ro = gorocksdb.NewDefaultReadOptions()
	database.wo = gorocksdb.NewDefaultWriteOptions()
	return database, nil
}

func (db *RocksDB) Get(key []byte) []byte {
	res, err := db.db.Get(db.ro, key)
	if err != nil {
		PanicCrisis(err)
	}
	if res.Size() == 0 {
		return nil
	}
	return res.Data()
}

func (db *RocksDB) Set(key []byte, value []byte) {
	err := db.db.Put(db.wo, key, value)
	if err != nil {
		PanicCrisis(err)
	}
}

func (db *RocksDB) SetSync(key []byte, value []byte) {
	wo := *db.wo
	wo.SetSync(true)
	err := db.db.Put(&wo, key, value)
	if err != nil {
		PanicCrisis(err)
	}
}

func (db *RocksDB) Delete(key []byte) {
	err := db.db.Delete(db.wo, key)
	if err != nil {
		PanicCrisis(err)
	}
}

func (db *RocksDB) DeleteSync(key []byte) {
	wo := *db.wo
	wo.SetSync(true)
	err := db.db.Delete(&wo, key)
	if err != nil {
		PanicCrisis(err)
	}
}

func (db *RocksDB) DB() *gorocksdb.DB {
	return db.db
}

func (db *RocksDB) Close() {
	db.db.Close()
}

func (db *RocksDB) Print() {
	iter := db.db.NewIterator(nil)
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		key := iter.Key().Data()
		value := iter.Value().Data()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
}

func (db *RocksDB) Query(key string, beg, end, bitSize int64) IQuery {
	iter := db.db.NewIterator(gorocksdb.NewDefaultReadOptions())
	fmtpad := "20"
	if bitSize == 32 {
		fmtpad = "10"
	}
	minkey := fmt.Sprintf("%s%0"+fmtpad+"d", key, beg)
	maxkey := fmt.Sprintf("%s%0"+fmtpad+"d", key, end)
	return &rdbQuery{iter, fmtpad, key, []byte(minkey), []byte(maxkey), false, true}
}

func (db *RocksDB) GetIter(key string, bitSize int64) IQuery {
	return db.Query(key, 0, math.MaxInt64, bitSize)
}

type rdbQuery struct {
	iter   *gorocksdb.Iterator
	fmtpad string
	key    string
	beg    []byte
	end    []byte
	isseek bool
	first  bool
}

func (q *rdbQuery) Seek(i int64) IQuery {
	//i = 0 表示最后一个
	seekprev := false
	if i == 0 { //seek to end
		i = math.MaxInt64
		seekprev = true
	}
	seekkey := fmt.Sprintf("%s%0"+q.fmtpad+"d", q.key, i)
	q.iter.Seek([]byte(seekkey))
	if seekprev {
		q.iter.Prev()
	}
	q.isseek = true
	return q
}

func (q *rdbQuery) FetchAll(limit int64, asc bool) (keys [][]byte, vals [][]byte, err error) {
	iter := q.iter
	defer iter.Close()
	i := 0
	if asc {
		for q.Next() {
			key := iter.Key().Data()
			val := iter.Value().Data()
			keys = append(keys, []byte(string(key)))
			vals = append(vals, val)
			i++
			if limit > 0 && i >= int(limit) {
				break
			}
		}
	} else {
		for q.Prev() {
			key := iter.Key().Data()
			val := iter.Value().Data()
			keys = append(keys, []byte(string(key)))
			vals = append(vals, val)
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

func (q *rdbQuery) Next() bool {
	if !q.isseek {
		q.iter.Seek(q.beg)
		q.isseek = true
	}
	if !q.first {
		q.iter.Next()
	}
	q.first = false
	end := q.iter.Key().Data()
	log.Println("next", toStr(end))
	if !bytes.HasPrefix(end, []byte(q.key)) || bytes.Compare(end, q.end) > 0 {
		return false
	}
	return q.iter.Valid()
}

func (q *rdbQuery) Key() []byte {
	return q.iter.Key().Data()
}

func (q *rdbQuery) Close() {
	q.iter.Close()
}

func (q *rdbQuery) Error() error {
	return q.iter.Err()
}

func (q *rdbQuery) Value() []byte {
	return q.iter.Value().Data()
}

func (q *rdbQuery) Prev() bool {
	if !q.isseek {
		q.iter.Seek(q.end)
		q.isseek = true
	}
	if !q.first {
		q.iter.Prev()
	}
	q.first = false
	end := q.iter.Key().Data()
	log.Println("prev", toStr(end))
	if !bytes.HasPrefix(end, []byte(q.key)) || bytes.Compare(q.beg, end) > 0 {
		return false
	}
	return q.iter.Valid()
}

type rdbBatch struct {
	db *gorocksdb.DB
	b  *gorocksdb.WriteBatch
	wo *gorocksdb.WriteOptions
}

//因为是金融系统，这个选项还是有必要，避免发生丢数据的情况
//而我们设计整个系统尽量采用批量提交，所以性能应该还是可以的
func (d *RocksDB) NewBatch(sync bool) IBatch {
	wop := gorocksdb.NewDefaultWriteOptions()
	wop.SetSync(sync)
	return &rdbBatch{db: d.db, b: gorocksdb.NewWriteBatch(), wo: wop}
}

func (b *rdbBatch) Put(key, value []byte) error {
	b.b.Put(key, value)
	return nil
}

func (b *rdbBatch) Delete(key []byte) error {
	b.b.Delete(key)
	return nil
}

func (b *rdbBatch) Commit() error {
	return b.db.Write(b.wo, b.b)
}

func init() {
	dbDrivers["rocksdb"] = NewRocksDB
}
