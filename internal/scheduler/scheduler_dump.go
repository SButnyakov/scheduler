package scheduler

import (
	"scheduler/internal/task"
	"sort"
	"time"
)

type scheduler_dump struct {
	Name           string
	CurrentTask    *task.Task
	ReadyQueues    [][]task.Task
	SuspendedQueue []task.Task
	WaitingQueues  [][]task.Task
	Timestamp      time.Time
}

type byTimestamp []scheduler_dump

func (a byTimestamp) Len() int           { return len(a) }
func (a byTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTimestamp) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }

func (s *Scheduler) dump(name string) {
	s.rMu.Lock()
	defer s.rMu.Unlock()
	s.sMu.Lock()
	defer s.sMu.Unlock()
	s.wMu.Lock()
	defer s.wMu.Unlock()

	dump := scheduler_dump{
		Name:           name,
		ReadyQueues:    make([][]task.Task, len(s.readyQueues)),
		SuspendedQueue: make([]task.Task, len(s.suspendedQueue)),
		WaitingQueues:  make([][]task.Task, len(s.waitingQueues)),
		Timestamp:      time.Now(),
	}
	if s.currentTask != nil {
		dump.CurrentTask = s.currentTask.Copy()
	}

	for i := range s.readyQueues {
		dump.ReadyQueues[i] = make([]task.Task, len(s.readyQueues[i]))
		for j, t := range s.readyQueues[i] {
			dump.ReadyQueues[i][j] = *t.Copy()
		}
	}
	for i, t := range s.suspendedQueue {
		dump.SuspendedQueue[i] = *t.Copy()
	}
	for i := range s.waitingQueues {
		dump.WaitingQueues[i] = make([]task.Task, len(s.waitingQueues[i]))
		for j, t := range s.waitingQueues[i] {
			dump.WaitingQueues[i][j] = *t.Copy()
		}
	}

	s.Dumps = append(s.Dumps, dump)
}

func (s *Scheduler) sortDumpsByTimestamp() {
	sort.Sort(byTimestamp(s.Dumps))
}
