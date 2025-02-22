package scheduler

import (
	"log"
	"scheduler/internal/task"
	"scheduler/internal/utils"
	"sync"
	"time"
)

const (
	WaitTaskTime  time.Duration = time.Duration(1000)
	MaxReadyTasks int           = 10
	ReadyLimit    int           = MaxReadyTasks - 2
)

type Scheduler struct {
	currentTask     *task.Task
	readyQueues     [4][]*task.Task
	suspendedQueues [4][]*task.Task
	waitingQueues   [4][]*task.Task
	sMu             sync.Mutex
	rMu             sync.Mutex
	stopChan        chan struct{}
	interruptChan   chan struct{}
}

func (s *Scheduler) Run() {
	s.stopChan = make(chan struct{})
	s.interruptChan = make(chan struct{})
	go s.manageQueues()
	go s.processTasks()
}

func (s *Scheduler) processTasks() {
	for {
		if s.currentTask == nil {
			nextTask := utils.PopNextFromQueues(s.readyQueues[:], &s.rMu)
			if nextTask == nil {
				time.Sleep(WaitTaskTime)
				continue
			}
			nextTask.SetState(task.Running)
			s.currentTask = nextTask
			go s.currentTask.Do()
		}
		select {
		case <-s.interruptChan:
			s.currentTask.Interrupt()
			s.currentTask.SetState(task.Ready)
			s.prependToReady(s.currentTask)
			s.currentTask = nil
			continue
		case <-s.currentTask.DoneChan:
			s.currentTask.SetState(task.Suspended)
			s.currentTask = nil
			continue
		default:
			continue
		}
	}
}

func (s *Scheduler) manageQueues() {
	for {
		//log.Printf("%d\n", utils.LMDA(s.suspendedQueues[:]))
		if utils.LMDA(s.readyQueues[:]) < ReadyLimit && utils.LMDA(s.suspendedQueues[:]) > 0 {
			t := utils.PopNextFromQueues(s.suspendedQueues[:], &s.sMu)
			t.SetState(task.Ready)
			s.appendToReady(t)
		}
	}
}

func (s *Scheduler) AddNewTask(t *task.Task) {
	t.SetState(task.Suspended)
	s.sMu.Lock()
	s.suspendedQueues[t.GetPriority()] = append(s.suspendedQueues[t.GetPriority()], t)
	s.sMu.Unlock()
}

func (s *Scheduler) appendToReady(t *task.Task) {
	s.rMu.Lock()
	s.readyQueues[t.GetPriority()] = append(s.readyQueues[t.GetPriority()], t)
	s.rMu.Unlock()
	s.checkInterruption(t)
}

func (s *Scheduler) prependToReady(t *task.Task) {
	s.rMu.Lock()
	s.readyQueues[t.GetPriority()] = append([]*task.Task{t}, s.readyQueues[t.GetPriority()]...)
	s.rMu.Unlock()
	s.checkInterruption(t)
}

func (s *Scheduler) appendToWaiting(t *task.Task) {
	s.waitingQueues[t.GetPriority()] = append(s.waitingQueues[t.GetPriority()], t)
}

func (s *Scheduler) checkInterruption(t *task.Task) {
	if s.currentTask != nil && s.currentTask.GetPriority() < t.GetPriority() {
		log.Printf("INTERRUPTING")
		s.interruptChan <- struct{}{}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopChan)
}
