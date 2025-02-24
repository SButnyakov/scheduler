package generator

import (
	"log"
	"math/rand"
	"scheduler/internal/task"
)

func GenerateTask() (*task.Task, error) {
	types := []task.TaskType{task.Basic, task.Extended}
	priorities := []task.TaskPriority{task.P0, task.P1, task.P2, task.P3}

	typeIndex := rand.Intn(len(types))
	priorityIndex := rand.Intn(len(priorities))

	t, err := task.New(types[typeIndex], priorities[priorityIndex], task.Suspended)
	if err != nil {
		return nil, err
	}

	log.Printf("Generated task | ID=%d | type=%s | p=%d", t.ID, t.GetType(), t.GetPriority())

	return t, nil
}
