package utils

import (
	"scheduler/internal/task"
	"sync"
)

// length of multidimensional array
func LMDA[V any](arr [][]V) int {
	res := 0
	for _, v := range arr {
		res += len(v)
	}
	return res
}

func PopNextFromQueues[T any](queues [][]*T, mu *sync.Mutex) *T {
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	var res *T
	for i := task.P3; i >= task.P0; i-- {
		if len(queues[i]) > 0 {
			res = queues[i][0]
			queues[i] = queues[i][:1]
			return res
		}
	}
	return nil
}
