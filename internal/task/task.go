package task

import (
	"log"
	"scheduler/internal/utils"
	"sync"
	"time"
)

type TaskState string

type TaskType string

type TaskPriority int

const (
	ArgsStateIndex int = 0
)

const (
	Running   TaskState = "running"
	Ready     TaskState = "ready"
	Waiting   TaskState = "waiting"
	Suspended TaskState = "suspended"
)

const (
	Basic    TaskType = "basic"
	Extended TaskType = "extended"
)

const (
	P0 = iota
	P1
	P2
	P3
)

const (
	ProgressLimit int           = 5
	TaskSleepTime time.Duration = time.Duration(2 * time.Second)
)

type Task struct {
	ID            int
	tType         TaskType
	priority      TaskPriority
	state         TaskState
	progress      int
	interruptChan chan struct{}
	DoneChan      chan struct{}
	WaitChan      chan struct{}
	EventChan     chan struct{}
}

var nextTaskID = 0
var mu sync.Mutex

func New(tType TaskType, priority TaskPriority, state TaskState) (*Task, error) {
	mu.Lock()
	defer mu.Unlock()
	t := Task{ID: nextTaskID}
	if err := t.SetPriority(priority); err != nil {
		return nil, err
	}
	if err := t.SetState(state); err != nil {
		return nil, err
	}
	if err := t.SetType(tType); err != nil {
		return nil, err
	}
	nextTaskID++
	t.interruptChan = make(chan struct{})
	t.DoneChan = make(chan struct{})
	t.WaitChan = make(chan struct{})
	t.EventChan = make(chan struct{}, 1)
	return &t, nil
}

func (t *Task) Do() {
	for {
		select {
		case <-t.interruptChan:
			log.Printf("Task %d interruption\n", t.ID)
			return
		default:
			if t.progress >= ProgressLimit {
				t.DoneChan <- struct{}{}
				return
			}
			if t.tType == Extended && t.progress == ProgressLimit/2 && !utils.IsChannelClosed(t.WaitChan) {
				log.Printf("Task %d waiting...\n", t.ID)
				t.WaitChan <- struct{}{}
				return
			}
			if t.tType == Basic && t.progress == ProgressLimit/2 && !utils.IsChannelClosed(t.EventChan) {
				log.Printf("Task %d event release.\n", t.ID)
				t.EventChan <- struct{}{}
			}
			t.progress++
			log.Printf("Task %d progress: %d/%d\n", t.ID, t.progress, ProgressLimit)
			time.Sleep(time.Duration(TaskSleepTime))
		}
	}
}

func (t *Task) Interrupt() {
	t.interruptChan <- struct{}{}
}

func (t *Task) SetType(newType TaskType) error {
	if !isValidType(newType) {
		return ErrInvalidType
	}
	t.tType = newType
	return nil
}

func (t *Task) GetType() TaskType {
	return t.tType
}

func (t *Task) SetPriority(newPriority TaskPriority) error {
	if !isValidPriority(newPriority) {
		return ErrInvalidPriority
	}
	t.priority = newPriority
	return nil
}

func (t *Task) GetPriority() TaskPriority {
	return t.priority
}

func (t *Task) SetState(newState TaskState) error {
	if !isValidState(newState) {
		return ErrInvalidState
	}
	t.state = newState
	return nil
}

func (t *Task) GetState() TaskState {
	return t.state
}

func isValidState(newState TaskState) bool {
	return newState == Running ||
		newState == Ready ||
		newState == Waiting ||
		newState == Suspended
}

func isValidType(newType TaskType) bool {
	return newType == Basic ||
		newType == Extended
}

func isValidPriority(newP TaskPriority) bool {
	return newP >= P0 && newP <= P3
}
