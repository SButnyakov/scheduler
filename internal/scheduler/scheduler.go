package scheduler

import (
	"fmt"
	"log"
	"log/slog"
	"scheduler/internal/task"
	"scheduler/internal/utils"
	"sync"
	"time"
)

const (
	WaitTaskTime  time.Duration = 100 * time.Millisecond
	MaxReadyTasks int           = 5
	ReadyLimit    int           = MaxReadyTasks - 2
)

type Scheduler struct {
	currentTask    *task.Task
	readyQueues    [][]*task.Task
	suspendedQueue []*task.Task
	waitingQueues  [][]*task.Task
	sMu            sync.Mutex
	rMu            sync.Mutex
	wMu            sync.Mutex
	ctMu           sync.Mutex
	StopChan       chan struct{}
	interruptChan  chan struct{}
	Dumps          []scheduler_dump
	noTasksCount   int
}

func New() *Scheduler {
	s := Scheduler{}
	s.StopChan = make(chan struct{})
	s.interruptChan = make(chan struct{})
	s.readyQueues = make([][]*task.Task, 4)
	s.suspendedQueue = make([]*task.Task, 0)
	s.waitingQueues = make([][]*task.Task, 4)
	for i := task.P0; i < task.P3; i++ {
		s.readyQueues[i] = make([]*task.Task, 0)
		s.waitingQueues[i] = make([]*task.Task, 0)
	}
	s.Dumps = make([]scheduler_dump, 0)
	return &s
}

func (s *Scheduler) Run() {
	go s.processTasks()
	go s.manageQueues()
}

func (s *Scheduler) processTasks() {
	for {
		if s.currentTask == nil {
			t := s.popNextFromReady()
			if t == nil {
				time.Sleep(WaitTaskTime)
				s.noTasksCount++
				if s.noTasksCount >= 10 {
					s.sortDumpsByTimestamp()
					log.Println("PROCESS EXITING...")
					close(s.StopChan)
					return
				}
				continue
			}
			s.noTasksCount = 0
			s.currentTask = t
			s.currentTask.SetState(task.Running)
			s.dump(fmt.Sprintf("task ready -> running | ID=%d\n", t.ID))
			go s.currentTask.Do()
		}
		select {
		case <-s.interruptChan:
			slog.Debug("PROCESS TASK INTERRUPTING")
			s.currentTask.Interrupt()
			slog.Debug("PROCESS TASK INTERRUPTED")
			s.currentTask.SetState(task.Ready)
			s.prependToReady(s.currentTask)
			id := s.currentTask.ID
			s.nilCurrentTask()
			s.dump(fmt.Sprintf("task running -> ready | ID=%d\n", id))
		case <-s.currentTask.DoneChan:
			slog.Debug("PROCESS TASK DONE")
			s.currentTask.SetState(task.Suspended)
			id := s.currentTask.ID
			s.nilCurrentTask()
			s.dump(fmt.Sprintf("task running -> suspended | ID=%d\n", id))
		case _, ok := <-s.currentTask.WaitChan:
			if !ok {
				continue
			}
			slog.Debug("PROCESS TASK WAIT")
			s.currentTask.SetState(task.Waiting)
			s.appendToWaiting(s.currentTask)
			close(s.currentTask.WaitChan)
			id := s.currentTask.ID
			s.nilCurrentTask()
			s.dump(fmt.Sprintf("task running -> waiting | ID=%d\n", id))
		}
		slog.Debug("PROCESS", slog.Any("SUS", s.suspendedQueue))
		slog.Debug("PROCESS", slog.Any("REA", s.readyQueues))
		slog.Debug("PROCESS", slog.Any("WAI", s.waitingQueues))
	}
}

func (s *Scheduler) manageQueues() {
	for {
		select {
		case <-s.StopChan:
			log.Println("MANAGE EXITING...")
			return
		default:
		}

		if utils.LMDA(s.readyQueues) < ReadyLimit && len(s.suspendedQueue) > 0 {
			slog.Debug("MANAGE", slog.Any("SUS", s.suspendedQueue))
			slog.Debug("MANAGE", slog.Any("REA", s.readyQueues))
			t := s.popNextFromSuspended()
			s.appendToReady(t)
			slog.Debug("MANAGE", slog.Any("SUS", s.suspendedQueue))
			slog.Debug("MANAGE", slog.Any("REA", s.readyQueues))
			s.dump(fmt.Sprintf("task suspended -> ready | ID=%d\n", t.ID))
		}
		s.ctMu.Lock()
		if s.currentTask != nil {
			select {
			case _, ok := <-s.currentTask.EventChan:
				if !ok {
					s.ctMu.Unlock()
					continue
				}
				t := s.popNextFromWaiting()
				if t == nil {
					s.ctMu.Unlock()
					continue
				}
				close(s.currentTask.EventChan)
				s.ctMu.Unlock()
				s.prependToReady(t)
				s.dump(fmt.Sprintf("task waiting -> ready | ID=%d\n", t.ID))
			default:
				s.ctMu.Unlock()
				continue
			}
		} else {
			s.ctMu.Unlock()
		}
	}
}

func (s *Scheduler) AddNewTask(t *task.Task) {
	s.sMu.Lock()
	t.SetState(task.Suspended)
	s.suspendedQueue = append(s.suspendedQueue, t)
	s.sMu.Unlock()
	s.dump(fmt.Sprintf("task -> suspended | ID=%d\n", t.ID))
}

func (s *Scheduler) appendToReady(t *task.Task) {
	s.rMu.Lock()
	defer s.rMu.Unlock()
	t.SetState(task.Ready)
	s.readyQueues[t.GetPriority()] = append(s.readyQueues[t.GetPriority()], t)
	s.checkInterruption(t)
}

func (s *Scheduler) prependToReady(t *task.Task) {
	s.rMu.Lock()
	defer s.rMu.Unlock()
	t.SetState(task.Ready)
	s.readyQueues[t.GetPriority()] = append([]*task.Task{t}, s.readyQueues[t.GetPriority()]...)
	s.checkInterruption(t)
}

func (s *Scheduler) appendToWaiting(t *task.Task) {
	s.wMu.Lock()
	defer s.wMu.Unlock()
	s.waitingQueues[t.GetPriority()] = append(s.waitingQueues[t.GetPriority()], t)
}

func (s *Scheduler) popNextFromSuspended() *task.Task {
	s.sMu.Lock()
	defer s.sMu.Unlock()
	if len(s.suspendedQueue) == 0 {
		return nil
	}
	res := s.suspendedQueue[0]
	s.suspendedQueue = s.suspendedQueue[1:]
	return res
}

func (s *Scheduler) popNextFromReady() *task.Task {
	s.rMu.Lock()
	defer s.rMu.Unlock()
	for i := task.P3; i >= task.P0; i-- {
		if len(s.readyQueues[i]) > 0 {
			res := s.readyQueues[i][0]
			s.readyQueues[i] = s.readyQueues[i][1:]
			return res
		}
	}
	return nil
}

func (s *Scheduler) popNextFromWaiting() *task.Task {
	s.wMu.Lock()
	defer s.wMu.Unlock()
	for i := task.P3; i >= task.P0; i-- {
		if len(s.waitingQueues[i]) > 0 {
			res := s.waitingQueues[i][0]
			s.waitingQueues[i] = s.waitingQueues[i][1:]
			return res
		}
	}
	return nil
}

func (s *Scheduler) nilCurrentTask() {
	s.ctMu.Lock()
	defer s.ctMu.Unlock()
	s.currentTask = nil
}

func (s *Scheduler) checkInterruption(t *task.Task) {
	if s.currentTask != nil && s.currentTask.GetPriority() < t.GetPriority() {
		select {
		case s.interruptChan <- struct{}{}:
			log.Println("Interrupt signal sent")
		default:
			slog.Debug("Interrupt channel is busy")
		}
	}
}
