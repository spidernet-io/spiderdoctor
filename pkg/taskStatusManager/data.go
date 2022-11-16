package taskStatusManager

import "github.com/spidernet-io/spiderdoctor/pkg/lock"

type taskStatus struct {
	l          lock.RWMutex
	taskStatus map[string]RoundStatus
}

type RoundStatus string

const (
	RoundStatusOngoing   = "started"
	RoundStatusSucceeded = "succeed"
	RoundStatusFail      = "fail"
)

type TaskStatus interface {
	SetTask(taskName string, status RoundStatus)
	DeleteTask(taskName string)
	CheckTask(taskName string) (status RoundStatus, existed bool)
}

func NewTaskStatus() TaskStatus {
	return &taskStatus{
		l:          lock.RWMutex{},
		taskStatus: map[string]RoundStatus{},
	}
}

func (s *taskStatus) SetTask(taskName string, status RoundStatus) {
	s.l.Lock()
	defer s.l.Unlock()
	s.taskStatus[taskName] = status
}

func (s *taskStatus) DeleteTask(taskName string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.taskStatus, taskName)
}

func (s *taskStatus) CheckTask(taskName string) (status RoundStatus, existed bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	status, existed = s.taskStatus[taskName]
	return
}
