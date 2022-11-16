package taskStatus

import "github.com/spidernet-io/spiderdoctor/pkg/lock"

type TaskStatus struct {
	l          lock.RWMutex
	taskStatus map[string]bool
}

func (s *TaskStatus) SetTask(taskName string, finished bool) {
	s.l.Lock()
	defer s.l.Unlock()
	s.taskStatus[taskName] = finished
}

func (s *TaskStatus) DeleteTask(taskName string) {
	s.l.Lock()
	defer s.l.Unlock()
	delete(s.taskStatus, taskName)
}

func (s *TaskStatus) CheckTask(taskName string) (finished bool, existed bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	finished, existed = s.taskStatus[taskName]
	return
}
