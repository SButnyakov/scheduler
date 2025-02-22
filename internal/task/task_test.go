package task

import (
	"errors"
	"testing"
)

func TestTask_New_InvalidCases(t *testing.T) {
	tests := []struct {
		tType    TaskType
		priority TaskPriority
		state    TaskState
		wantErr  bool
	}{
		{"invalid", P1, Ready, true},
		{Basic, 5, Ready, true},
		{Basic, P1, "invalid", true},
		{Extended, P2, Running, false},
	}

	for _, tc := range tests {
		_, err := New(tc.tType, tc.priority, tc.state)
		if (err != nil) != tc.wantErr {
			t.Errorf("New(%v, %v, %v) ожидал ошибку: %v, получил: %v", tc.tType, tc.priority, tc.state, tc.wantErr, err)
		}
	}
}

func TestTaskState(t *testing.T) {
	task := Task{ID: 1, tType: Basic, priority: P2, state: Ready}

	// Проверка начального состояния
	if task.GetState() != Ready {
		t.Errorf("Ожидалось состояние %s, получено %s", Ready, task.GetState())
	}

	// Проверка изменения состояния
	err := task.SetState(Running)
	if err != nil {
		t.Errorf("Ошибка при установке состояния: %v", err)
	}
	if task.GetState() != Running {
		t.Errorf("Ожидалось состояние %s, получено %s", Running, task.GetState())
	}

	// Проверка недопустимого состояния
	err = task.SetState("invalid")
	if err == nil {
		t.Errorf("Ожидалась ошибка при установке недопустимого состояния")
	}
	if !errors.Is(err, ErrInvalidState) {
		t.Errorf("Ожидалась ошибка %v, получена %v", ErrInvalidState, err)
	}
}

func TestTaskType(t *testing.T) {
	task := Task{ID: 1, tType: Basic, priority: P2, state: Ready}

	// Проверка корректного изменения типа
	err := task.SetType(Extended)
	if err != nil {
		t.Errorf("Ошибка при установке типа: %v", err)
	}
	if task.GetType() != Extended {
		t.Errorf("Ожидался тип %s, получен %s", Extended, task.GetType())
	}

	// Проверка недопустимого типа
	err = task.SetType("invalid")
	if err == nil {
		t.Errorf("Ожидалась ошибка при установке недопустимого типа")
	}
	if !errors.Is(err, ErrInvalidType) {
		t.Errorf("Ожидалась ошибка %v, получена %v", ErrInvalidType, err)
	}
}

func TestTaskPriority(t *testing.T) {
	task := Task{ID: 1, tType: Basic, priority: P2, state: Ready}

	// Проверка корректного изменения приоритета
	err := task.SetPriority(P3)
	if err != nil {
		t.Errorf("Ошибка при установке приоритета: %v", err)
	}
	if task.GetPriority() != P3 {
		t.Errorf("Ожидался приоритет %d, получен %d", P3, task.GetPriority())
	}

	// Проверка недопустимого приоритета
	err = task.SetPriority(5)
	if err == nil {
		t.Errorf("Ожидалась ошибка при установке недопустимого приоритета")
	}
	if !errors.Is(err, ErrInvalidPriority) {
		t.Errorf("Ожидалась ошибка %v, получена %v", ErrInvalidPriority, err)
	}
}
