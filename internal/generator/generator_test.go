package generator

import (
	"scheduler/internal/task"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTask(t *testing.T) {
	for i := 0; i < 100; i++ { // Generate multiple tasks to cover randomness
		t.Run("GenerateTask", func(t *testing.T) {
			generatedTask, err := GenerateTask()
			assert.NoError(t, err)
			assert.NotNil(t, generatedTask)

			// Check that the generated task has a valid type
			assert.Contains(t, []task.TaskType{task.Basic, task.Extended}, generatedTask.GetType())

			// Check that the generated task has a valid priority
			assert.Contains(t, []task.TaskPriority{task.P0, task.P1, task.P2, task.P3}, generatedTask.GetPriority())

			// Check that the generated task is in the suspended state
			assert.Equal(t, task.Suspended, generatedTask.GetState())
		})
	}
}
