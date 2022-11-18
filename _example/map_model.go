package _example

import (
	"fmt"
	"reflect"

	"github.com/anishathalye/porcupine"
)

const (
	GetOp = iota
	PutOp = iota
	DelOp = iota
)

type MapInput struct {
	Operation uint8
	Key       int
	Val       int
}

type MapOutput struct {
	Key   int
	Val   int
	Found bool
}

type MapState struct {
	m map[int]int
}

func (s MapState) Clone() MapState {
	sx := make(map[int]int, len(s.m))
	for k, v := range s.m {
		sx[k] = v
	}
	return MapState{m: sx}
}

func (s MapState) Equals(otherState MapState) bool {
	return reflect.DeepEqual(s.m, otherState.m)
}

func (s MapState) String() string {
	return fmt.Sprintf("%v", s.m)
}

func NewMapState() MapState {
	return MapState{m: map[int]int{}}
}

var MapModel = porcupine.Model[MapState, MapInput, MapOutput]{
	Init: NewMapState,
	Partition: func(history []porcupine.Operation[MapInput, MapOutput]) [][]porcupine.Operation[MapInput, MapOutput] {
		indexMap := map[int]int{}
		var partitions [][]porcupine.Operation[MapInput, MapOutput]
		for _, op := range history {
			ix, found := indexMap[op.Input.Key]
			if !found {
				partitions = append(partitions, []porcupine.Operation[MapInput, MapOutput]{op})
				indexMap[op.Input.Key] = len(partitions) - 1
			} else {
				partitions[ix] = append(partitions[ix], op)
			}
		}
		return partitions
	},
	Step: func(state MapState, input MapInput, output MapOutput) (bool, MapState) {
		stateVal, found := state.m[input.Key]

		switch input.Operation {
		case GetOp:
			if !output.Found {
				return !found, state
			} else if stateVal == output.Val {
				return true, state
			}
			break
		case PutOp:
			state.m[input.Key] = input.Val
			return true, state
		case DelOp:
			delete(state.m, input.Key)
			return true, state
		}

		return false, state
	},
	DescribeOperation: func(i MapInput, o MapOutput) string {
		opName := ""
		switch i.Operation {
		case GetOp:
			opName = "Get"
			break
		case PutOp:
			opName = "Put"
			break
		case DelOp:
			opName = "Del"
			break
		}

		return fmt.Sprintf("%s(%d) -> %d", opName, i.Key, o.Val)
	},
}
