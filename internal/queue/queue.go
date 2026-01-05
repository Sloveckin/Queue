package queue

import (
	"container/list"
	"errors"
)

var (
	NameAlreadyExists = errors.New("name already exists")
	IsEmpty           = errors.New("queue is empty")
)

type Queue struct {
	Names map[string]bool
	// TODO: write own list
	List list.List
}

func New() *Queue {
	return &Queue{
		Names: make(map[string]bool),
		List:  list.List{},
	}
}

func (q *Queue) Add(name *string) error {
	if _, ok := q.Names[*name]; ok {
		return NameAlreadyExists
	}

	q.Names[*name] = true
	q.List.PushBack(*name)

	return nil
}

func (q *Queue) Next() error {
	front := q.List.Front()
	if front == nil {
		return IsEmpty
	}

	delete(q.Names, front.Value.(string))
	q.List.Remove(front)

	return nil
}
