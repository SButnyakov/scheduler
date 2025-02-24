package scheduler

import (
	"fmt"
	"scheduler/internal/generator"
	"scheduler/internal/task"
	"scheduler/internal/utils"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_Requirements(t *testing.T) {
	testsAmount := 5
	var wg sync.WaitGroup
	wg.Add(testsAmount)
	for i := 0; i < testsAmount; i++ {
		go func() {
			defer wg.Done()
			tasksAmount := 5
			s := New()
			tasks := make([]*task.Task, 0, tasksAmount)

			for i := 0; i < tasksAmount; i++ {
				task, err := generator.GenerateTask()
				assert.NoError(t, err)
				s.AddNewTask(task)
				tasks = append(tasks, task)
			}

			tasksMsg := ""
			for _, v := range tasks {
				tasksMsg += fmt.Sprintf("task | ID=%d | type=%s | priority=%d\n", v.ID, v.GetType(), v.GetPriority())
			}
			assert.True(t, true)

			s.Run()
			<-s.StopChan

			for _, v := range tasks {
				err := checkFirstInstanceIsSuspended(v, s.Dumps)
				assert.NoError(t, err, tasksMsg)
				err = checkStatesStoryCorrectness(v, s.Dumps)
				assert.NoError(t, err, tasksMsg)
			}

			err := checkReadyQueuesLimit(s.Dumps)
			assert.NoError(t, err)

			err = checkRunningTasksLimit(s.Dumps)
			assert.NoError(t, err)

			err = checkTaskPositionInReadyQueue(s.Dumps)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
}

func TestScheduler_TaskExecutionOrder(t *testing.T) {
	s := New()

	// Create tasks with specified priorities
	taskP3, err := task.New(task.Basic, task.P3, task.Suspended)
	assert.NoError(t, err)
	taskP2, err := task.New(task.Basic, task.P2, task.Suspended)
	assert.NoError(t, err)
	taskP1_1, err := task.New(task.Basic, task.P1, task.Suspended)
	assert.NoError(t, err)
	taskP1_2, err := task.New(task.Basic, task.P1, task.Suspended)
	assert.NoError(t, err)
	taskP0, err := task.New(task.Basic, task.P0, task.Suspended)
	assert.NoError(t, err)

	// Add tasks to the scheduler
	s.AddNewTask(taskP3)
	s.AddNewTask(taskP2)
	s.AddNewTask(taskP1_1)
	s.AddNewTask(taskP1_2)
	s.AddNewTask(taskP0)

	// Run the scheduler
	s.Run()
	<-s.StopChan

	// Check the order of task execution
	expectedOrder := []int{taskP3.ID, taskP2.ID, taskP1_1.ID, taskP1_2.ID, taskP0.ID}
	actualOrder := getTaskExecutionOrder(s.Dumps)

	assert.Equal(t, expectedOrder, actualOrder, "Tasks were not executed in the correct order")
}

func TestScheduler_TaskExecutionOrderWithWaiting(t *testing.T) {
	s := New()

	// Create tasks with specified priorities and types
	task1, err := task.New(task.Basic, task.P3, task.Suspended)
	assert.NoError(t, err)
	task2, err := task.New(task.Extended, task.P2, task.Suspended)
	assert.NoError(t, err)
	task3, err := task.New(task.Basic, task.P1, task.Suspended)
	assert.NoError(t, err)

	// Add tasks to the scheduler
	s.AddNewTask(task1)
	s.AddNewTask(task2)
	s.AddNewTask(task3)

	// Run the scheduler
	s.Run()
	<-s.StopChan

	// Check the order of task execution
	expectedOrder := []int{task1.ID, task2.ID, task3.ID, task2.ID, task3.ID}
	actualOrder := getTaskExecutionOrder(s.Dumps)

	assert.Equal(t, expectedOrder, actualOrder, "Tasks were not executed in the correct order")
}

func checkFirstInstanceIsSuspended(t *task.Task, dumps []scheduler_dump) error {
	for _, dump := range dumps {
		for _, v := range dump.SuspendedQueue {
			if v.ID == t.ID {
				if v.GetState() != task.Suspended {
					return fmt.Errorf("task first met in suspended but with wrong state | ID=%d | state=%s\n", v.ID, v.GetState())
				}
				return nil
			}
		}
		for _, queue := range dump.ReadyQueues {
			for _, v := range queue {
				if v.ID == t.ID {
					return fmt.Errorf("task first met in ready not suspended | ID=%d | state=%s\n", v.ID, v.GetState())
				}
			}
		}
		for _, queue := range dump.WaitingQueues {
			for _, v := range queue {
				if v.ID == t.ID {
					return fmt.Errorf("task first met in waiting not suspended | ID=%d | state=%s\n", v.ID, v.GetState())
				}
			}
		}
	}
	return nil
}

func checkStatesStoryCorrectness(t *task.Task, dumps []scheduler_dump) error {
	var prevState task.TaskState

	for _, dump := range dumps {
		dumpTask := findTaskInDump(t, dump)
		if dumpTask.ID == -1 {
			continue // Task not found in this dump
		}

		currentState := dumpTask.GetState()
		if currentState == prevState {
			continue
		}

		if prevState != "" {
			switch prevState {
			case task.Suspended:
				if currentState != task.Ready {
					return fmt.Errorf("invalid state transition | ID=%d | from=%s to=%s\n", dumpTask.ID, prevState, currentState)
				}
			case task.Ready:
				if currentState != task.Running {
					return fmt.Errorf("invalid state transition | ID=%d | from=%s to=%s\n", dumpTask.ID, prevState, currentState)
				}
			case task.Running:
				if currentState != task.Waiting && currentState != task.Ready && currentState != task.Suspended {
					return fmt.Errorf("invalid state transition | ID=%d | from=%s to=%s\n", dumpTask.ID, prevState, currentState)
				}
			case task.Waiting:
				if currentState != task.Ready {
					return fmt.Errorf("invalid state transition | ID=%d | from=%s to=%s\n", dumpTask.ID, prevState, currentState)
				}
			}
		}

		prevState = currentState
	}

	return nil
}

func checkReadyQueuesLimit(dumps []scheduler_dump) error {
	for _, dump := range dumps {
		totalReadyTasks := utils.LMDA(dump.ReadyQueues)
		if totalReadyTasks > MaxReadyTasks {
			return fmt.Errorf("number of tasks in readyQueues exceeds limit | Timestamp=%s | Count=%d\n", dump.Timestamp, totalReadyTasks)
		}
	}
	return nil
}

func getTaskExecutionOrder(dumps []scheduler_dump) []int {
	executionOrder := []int{}
	prevTaskID := -1

	for _, dump := range dumps {
		if dump.CurrentTask != nil && dump.CurrentTask.ID != prevTaskID {
			executionOrder = append(executionOrder, dump.CurrentTask.ID)
			prevTaskID = dump.CurrentTask.ID
		}
	}

	return executionOrder
}

func findTaskInDump(t *task.Task, dump scheduler_dump) task.Task {
	if dump.CurrentTask != nil && dump.CurrentTask.ID == t.ID {
		return *dump.CurrentTask
	}
	for _, v := range dump.SuspendedQueue {
		if v.ID == t.ID {
			return v
		}
	}
	for _, v := range dump.ReadyQueues[t.GetPriority()] {
		if v.ID == t.ID {
			return v
		}
	}
	for _, v := range dump.WaitingQueues[t.GetPriority()] {
		if v.ID == t.ID {
			return v
		}
	}
	return task.Task{ID: -1}
}

func checkRunningTasksLimit(dumps []scheduler_dump) error {
	for _, dump := range dumps {
		runningCount := 0

		if dump.CurrentTask != nil && dump.CurrentTask.GetState() == task.Running {
			runningCount++
		}

		for _, queue := range dump.ReadyQueues {
			for _, t := range queue {
				if t.GetState() == task.Running {
					runningCount++
				}
			}
		}

		for _, t := range dump.SuspendedQueue {
			if t.GetState() == task.Running {
				runningCount++
			}
		}

		for _, queue := range dump.WaitingQueues {
			for _, t := range queue {
				if t.GetState() == task.Running {
					runningCount++
				}
			}
		}

		if runningCount > 1 {
			return fmt.Errorf("more than one task with state Running in dump | Timestamp=%s | RunningCount=%d\n", dump.Timestamp, runningCount)
		}
	}
	return nil
}

func checkTaskPositionInReadyQueue(dumps []scheduler_dump) error {
	for _, dump := range dumps {
		var id int
		var priority task.TaskPriority
		var queue []task.Task

		if _, err := fmt.Sscanf(dump.Name, "task waiting -> ready | ID=%d", &id); err == nil {
			priority = findTaskPriorityByID(dump, id)
			queue = dump.ReadyQueues[priority]
			if len(queue) == 0 || queue[0].ID != id {
				return fmt.Errorf("task with ID=%d is not at the beginning of its priority readyQueue in dump: %s", id, dump.Name)
			}
		} else if _, err := fmt.Sscanf(dump.Name, "task running -> ready | ID=%d", &id); err == nil {
			priority = findTaskPriorityByID(dump, id)
			queue = dump.ReadyQueues[priority]
			if len(queue) == 0 || queue[0].ID != id {
				return fmt.Errorf("task with ID=%d is not at the beginning of its priority readyQueue in dump: %s", id, dump.Name)
			}
		} else if _, err := fmt.Sscanf(dump.Name, "task suspended -> ready | ID=%d", &id); err == nil {
			priority = findTaskPriorityByID(dump, id)
			queue = dump.ReadyQueues[priority]
			if len(queue) == 0 || queue[len(queue)-1].ID != id {
				return fmt.Errorf("task with ID=%d is not at the end of its priority readyQueue in dump: %s", id, dump.Name)
			}
		}
	}
	return nil
}

func findTaskPriorityByID(dump scheduler_dump, id int) task.TaskPriority {
	for priority, queue := range dump.ReadyQueues {
		for _, t := range queue {
			if t.ID == id {
				return task.TaskPriority(priority)
			}
		}
	}
	return -1
}
