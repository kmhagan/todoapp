package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"
	"todoapp/pkg/linkedlist"

	"github.com/google/uuid"
)

var (
	// Localstore is used to access in memory
	Localstore = local{
		m:    make(map[string]*List),
		lock: &sync.RWMutex{},
	}
)

type local struct {
	m    map[string]*List
	lock *sync.RWMutex
}

// List is a todo list
type List struct {
	UUID    string
	Name    string
	Created time.Time
	List    *linkedlist.List `json:"List,omitempty"`
}

func (l *local) AllLists() []*List {
	l.lock.Lock()
	defer l.lock.Unlock()
	var lists []*List
	for _, v := range l.m {
		lists = append(lists, v)
	}
	// sort by creation date
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Created.Before(lists[j].Created)
	})
	return lists
}

func (l *local) NewList(name string) (*List, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	uuid, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("failed to generate uuid: %w", err)
		return nil, err
	}
	newLinkedList := linkedlist.NewList()
	newList := &List{
		UUID:    uuid.String(),
		Name:    name,
		List:    &newLinkedList,
		Created: time.Now().UTC().Round(time.Millisecond),
	}
	l.m[uuid.String()] = newList
	return newList, nil
}

func (l *local) GetList(uuid string) *List {
	l.lock.Lock()
	defer l.lock.Unlock()
	item := l.m[uuid]
	return item
}

func (l *local) DeleteList(uuid string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.m, uuid)
}
