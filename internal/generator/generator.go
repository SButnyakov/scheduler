package generator

import (
	"log"
	"log/slog"
	"math/rand"
	"scheduler/internal/scheduler"
	"scheduler/internal/task"
)

func Generate(scheduler *scheduler.Scheduler) {
	for i := 0; i < 10; i++ {
		t, err := generateTask()
		if err != nil {
			slog.Error("failed to generate task")
			continue
		}
		scheduler.AddNewTask(t)
	}
}

func generateTask() (*task.Task, error) {
	types := []task.TaskType{task.Basic, task.Extended}
	priorities := []task.TaskPriority{task.P0, task.P1, task.P2, task.P3}
	//states := []task.TaskState{task.Ready, task.Suspended}

	typeIndex := rand.Intn(len(types))
	priorityIndex := rand.Intn(len(priorities))
	//stateIndex := rand.Intn(len(states))

	t, err := task.New(types[typeIndex], priorities[priorityIndex], task.Suspended)
	if err != nil {
		return nil, err
	}

	log.Printf("Generated task: ID=%d, Type=%s, Priority=%d", t.ID, t.GetType(), t.GetPriority())

	return t, nil
}
