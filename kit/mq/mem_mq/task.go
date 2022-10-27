package mem_mq

import (
	"container/list"
	"sync"

	"github.com/iotexproject/Bumblebee/kit/mq"
)

func New() *TaskManager {
	return &TaskManager{l: list.New(), m: map[string]*list.Element{}}
}

type TaskManager struct {
	m map[string]*list.Element
	l *list.List

	rwm sync.RWMutex
}

var _ mq.TaskManager = (*TaskManager)(nil)

func (tm *TaskManager) Push(ch string, t mq.Task) error {
	tm.rwm.Lock()
	defer tm.rwm.Unlock()

	tm.m[key(ch, t.ID())] = tm.l.PushBack(t)
	return nil
}

func (tm *TaskManager) Pop(ch string) (mq.Task, error) {
	tm.rwm.Lock()
	defer tm.rwm.Unlock()

	elem := tm.l.Front()
	if elem == nil {
		return nil, nil
	}
	tm.l.Remove(elem)

	t, ok := elem.Value.(mq.Task)
	if !ok {
		return nil, nil
	}

	k := key(ch, t.ID())
	if _, ok = tm.m[k]; !ok {
		return nil, nil
	}
	return t, tm.remove(k)
}

func (tm *TaskManager) Remove(ch string, id string) error {
	tm.rwm.Lock()
	defer tm.rwm.Unlock()

	return tm.remove(key(ch, id))
}

func (tm *TaskManager) Clear(_ string) error {
	*tm = *New()
	return nil
}

func (tm *TaskManager) remove(key string) error {
	if t := tm.m[key]; t != nil {
		tm.l.Remove(t)
		delete(tm.m, key)
	}
	return nil
}

func key(ch, id string) string { return ch + "::" + id }
