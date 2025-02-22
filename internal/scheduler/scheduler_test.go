package scheduler

import (
	"scheduler/internal/task"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_NewTask(t *testing.T) {
	// Создаем новый планировщик
	scheduler := &Scheduler{}

	// Создаем реальную задачу с приоритетом
	realTask, err := task.New(task.Basic, task.P0, task.Ready)
	assert.NoError(t, err, "task.New shouldn't return any errors")

	// Вызываем функцию NewTask
	scheduler.AddNewTask(realTask)

	// Проверяем, что состояние задачи изменилось на Suspended
	assert.Equal(t, task.Suspended, realTask.GetState(), "Task state should be Suspended")

	// Проверяем, что задача добавлена в соответствующую очередь приостановленных задач
	testTask := scheduler.suspendedQueues[realTask.GetPriority()][0]
	assert.Contains(t, scheduler.suspendedQueues[0], realTask, "Task should be in the suspended queue")
	assert.Equal(t, task.Suspended, testTask.GetState(), "Task state should be equal suspended")
}
