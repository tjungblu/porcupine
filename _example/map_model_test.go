package _example

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/anishathalye/porcupine"
)

func TestHappyPath(t *testing.T) {
	db := NewDatabase()

	key := 42
	client := NewDatabaseRecorder(db, 0)
	for i := 0; i < 100; i++ {
		_, _ = client.Get(key)
		client.Put(key, i)
		if rand.Float32() < 0.25 {
			client.Del(key)
		}
	}

	result, info := porcupine.CheckOperationsVerbose(MapModel, client.operations, 0)
	err := porcupine.VisualizePath(MapModel, info, t.Name()+"_porcupine.html")
	if err != nil {
		t.Fatal(err)
	}
	if result != porcupine.Ok {
		t.Fatal("unexpected result does not match OK state")
	}
}

func TestMultiGoroutines(t *testing.T) {
	numGoRoutines := 4
	d := NewDatabase()

	var operations []porcupine.Operation[MapInput, MapOutput]
	var opsLock sync.Mutex
	wg := sync.WaitGroup{}

	for n := 0; n < numGoRoutines; n++ {
		wg.Add(1)
		go func(db *db, id int) {
			client := NewDatabaseRecorder(db, id)

			for j := 0; j < 100; j++ {
				for i := 0; i < 10; i++ {
					_, _ = client.Get(i)
					client.Put(i, i)
					if rand.Float32() < 0.25 {
						client.Del(i)
					}
				}
			}

			opsLock.Lock()
			defer opsLock.Unlock()

			operations = append(operations, client.operations...)

			wg.Done()
		}(d, n)
	}

	wg.Wait()

	result, info := porcupine.CheckOperationsVerbose(MapModel, operations, 0)
	err := porcupine.VisualizePath(MapModel, info, t.Name()+"_porcupine.html")
	if err != nil {
		t.Fatal(err)
	}
	if result != porcupine.Ok {
		t.Fatalf("unexpected result does not match OK state: %v", result)
	}
}
