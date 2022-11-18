package _example

import (
	"sync"
	"time"

	"github.com/anishathalye/porcupine"
)

type db struct {
	rwLock *sync.RWMutex
	db     map[int]int
}

type DatabaseClientRecorder struct {
	clientId int
	db       *db

	operations []porcupine.Operation[MapInput, MapOutput]
}

func (d *db) Put(key, val int) {
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	d.db[key] = val
}

func (d *db) Get(key int) (int, bool) {
	d.rwLock.RLock()
	defer d.rwLock.RUnlock()

	v, ok := d.db[key]
	return v, ok
}

func (d *db) Del(key int) {
	d.rwLock.Lock()
	defer d.rwLock.Unlock()

	delete(d.db, key)
}

func (d *DatabaseClientRecorder) Get(key int) (int, bool) {
	start := time.Now()
	val, found := d.db.Get(key)
	end := time.Now()
	d.operations = append(d.operations, porcupine.Operation[MapInput, MapOutput]{
		ClientId: d.clientId,
		Input: MapInput{
			Operation: GetOp,
			Key:       key,
		},
		Call: start.UnixNano(),
		Output: MapOutput{
			Key:   key,
			Val:   val,
			Found: found,
		},
		Return: end.UnixNano(),
	})

	return val, found
}

func (d *DatabaseClientRecorder) Put(key, value int) {
	start := time.Now()
	d.db.Put(key, value)
	end := time.Now()
	d.operations = append(d.operations, porcupine.Operation[MapInput, MapOutput]{
		ClientId: d.clientId,
		Input: MapInput{
			Operation: PutOp,
			Key:       key,
			Val:       value,
		},
		Call: start.UnixNano(),
		Output: MapOutput{
			Key: key,
			Val: value,
		},
		Return: end.UnixNano(),
	})

}

func (d *DatabaseClientRecorder) Del(key int) {
	start := time.Now()
	d.db.Del(key)
	end := time.Now()
	d.operations = append(d.operations, porcupine.Operation[MapInput, MapOutput]{
		ClientId: d.clientId,
		Input: MapInput{
			Operation: DelOp,
			Key:       key,
			Val:       0,
		},
		Call: start.UnixNano(),
		Output: MapOutput{
			Key: key,
			Val: 0,
		},
		Return: end.UnixNano(),
	})
}

func NewDatabase() *db {
	return &db{
		rwLock: &sync.RWMutex{},
		db:     map[int]int{},
	}
}

func NewDatabaseRecorder(db *db, clientId int) *DatabaseClientRecorder {
	return &DatabaseClientRecorder{
		clientId:   clientId,
		db:         db,
		operations: []porcupine.Operation[MapInput, MapOutput]{},
	}
}
