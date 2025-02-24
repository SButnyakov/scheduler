package task

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, Basic, task.GetType())
	assert.Equal(t, TaskPriority(P1), task.GetPriority())
	assert.Equal(t, Ready, task.GetState())
}

func TestSetType(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)

	err = task.SetType(Extended)
	assert.NoError(t, err)
	assert.Equal(t, Extended, task.GetType())

	err = task.SetType("invalid")
	assert.Error(t, err)
}

func TestSetPriority(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)

	err = task.SetPriority(P2)
	assert.NoError(t, err)
	assert.Equal(t, TaskPriority(P2), task.GetPriority())

	err = task.SetPriority(10)
	assert.Error(t, err)
}

func TestSetState(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)

	err = task.SetState(Running)
	assert.NoError(t, err)
	assert.Equal(t, Running, task.GetState())

	err = task.SetState("invalid")
	assert.Error(t, err)
}

func TestDo(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)

	go task.Do()
	time.Sleep(TaskSleepTime * 3)
	task.Interrupt()
	time.Sleep(TaskSleepTime)

	assert.GreaterOrEqual(t, task.progress, 2)
}

func TestDoExtended(t *testing.T) {
	task, err := New(Extended, P1, Ready)
	assert.NoError(t, err)

	go task.Do()
	time.Sleep(TaskSleepTime * 3)
	assert.GreaterOrEqual(t, task.progress, 2)
	select {
	case <-task.WaitChan:
	default:
		assert.True(t, false)
	}
}

func TestInterrupt(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)

	go task.Do()
	time.Sleep(TaskSleepTime / 2)
	task.Interrupt()
	time.Sleep(TaskSleepTime)

	assert.Equal(t, 1, task.progress)
}

func TestCopy(t *testing.T) {
	task, err := New(Basic, P1, Ready)
	assert.NoError(t, err)

	taskCopy := task.Copy()
	assert.Equal(t, task.ID, taskCopy.ID)
	assert.Equal(t, task.GetType(), taskCopy.GetType())
	assert.Equal(t, task.GetPriority(), taskCopy.GetPriority())
	assert.Equal(t, task.GetState(), taskCopy.GetState())
	assert.Equal(t, task.progress, taskCopy.progress)
}

func TestIsValidState(t *testing.T) {
	assert.True(t, isValidState(Running))
	assert.True(t, isValidState(Ready))
	assert.True(t, isValidState(Waiting))
	assert.True(t, isValidState(Suspended))
	assert.False(t, isValidState("invalid"))
}

func TestIsValidType(t *testing.T) {
	assert.True(t, isValidType(Basic))
	assert.True(t, isValidType(Extended))
	assert.False(t, isValidType("invalid"))
}

func TestIsValidPriority(t *testing.T) {
	assert.True(t, isValidPriority(P0))
	assert.True(t, isValidPriority(P1))
	assert.True(t, isValidPriority(P2))
	assert.True(t, isValidPriority(P3))
	assert.False(t, isValidPriority(10))
}
