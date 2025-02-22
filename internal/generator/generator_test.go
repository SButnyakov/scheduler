package generator

import (
	"scheduler/internal/task"
	"testing"
)

func TestGenerateTask(t *testing.T) {
	generatedTasks := make(map[int]*task.Task)

	for i := 0; i < 100; i++ {
		tsk := GenerateTask()
		if tsk == nil {
			t.Errorf("GenerateTask вернул nil")
		}

		// Проверка, что сгенерированная задача имеет допустимые параметры
		if tsk.GetType() != task.Basic && tsk.GetType() != task.Extended {
			t.Errorf("Недопустимый тип задачи: %v", tsk.GetType())
		}
		if tsk.GetPriority() < task.P0 || tsk.GetPriority() > task.P3 {
			t.Errorf("Недопустимый приоритет задачи: %v", tsk.GetPriority())
		}
		if tsk.GetState() != task.Ready && tsk.GetState() != task.Suspended {
			t.Errorf("Недопустимое состояние задачи: %v", tsk.GetState())
		}

		// Проверка уникальности ID
		if _, exists := generatedTasks[tsk.ID]; exists {
			t.Errorf("Дублирующийся ID задачи: %d", tsk.ID)
		}
		generatedTasks[tsk.ID] = tsk
	}
}
