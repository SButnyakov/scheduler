package generator

import (
	"log"
	"math/rand"
	"scheduler/internal/task"
)

func GenerateTask() *task.Task {
	types := []task.TaskType{task.Basic, task.Extended}
	priorities := []task.TaskPriority{task.P0, task.P1, task.P2, task.P3}
	//states := []task.TaskState{task.Ready, task.Suspended}

	typeIndex := rand.Intn(len(types))
	priorityIndex := rand.Intn(len(priorities))
	//stateIndex := rand.Intn(len(states))

	t, err := task.New(types[typeIndex], priorities[priorityIndex], task.Suspended)
	if err != nil {
		return nil
	}

	log.Printf("Generated task: ID=%d, Type=%s, Priority=%d", t.ID, t.GetType(), t.GetPriority())

	return t
}
