package taskStatus

import "github.com/spidernet-io/spiderdoctor/pkg/lock"

type taskStatus struct {
	l          lock.RWMutex
	taskStatus map[string]bool
}

type TaskStatus interface {
	SetTask(taskName string, finished bool)
	DeleteTask(taskName string)
	CheckTask(taskName string) (finished bool, existed bool)
}

func NewTaskStatus() TaskStatus {
	return &taskStatus{
		l:          lock.RWMutex{},
		taskStatus: map[string]bool{},
	}
}

func (s *taskStatus) SetTask(taskName string, finished bool) {
	s.l.Lock()
	defer s.l.Unlock()
	s.taskStatus[taskName] = finished
}

func (s *taskStatus) DeleteTask(taskName string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.taskStatus, taskName)
}

func (s *taskStatus) CheckTask(taskName string) (finished bool, existed bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	finished, existed = s.taskStatus[taskName]
	return
}
