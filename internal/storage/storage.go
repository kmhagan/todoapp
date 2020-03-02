package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"
	"todoapp/pkg/linkedlist"

	"github.com/google/uuid"
)

var (
	// Localstore is used to access in memory
	Localstore = local{
		M:         make(map[string]*List),
		lock:      &sync.RWMutex{},
		MaxTotal:  1000,
		Filestore: "data.json",
	}
)

type local struct {
	M                map[string]*List
	lock             *sync.RWMutex
	changesSinceSave bool
	MaxTotal         int
	Total            int
	Filestore        string
}

// List is a todo list
type List struct {
	UUID    string
	Name    string
	Created time.Time
	List    *linkedlist.List `json:"List,omitempty"`
}

func (l *local) Load() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	exists := fileExists(l.Filestore)
	if exists == false {
		return nil
	}
	data, err := ioutil.ReadFile(l.Filestore)
	if err != nil {
		return err
	}
	var list local
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err
	}
	l.M = list.M
	l.Total = list.Total
	return nil
}

func (l *local) Save() error {
	l.lock.Lock()
	defer l.lock.Unlock()
	if !l.changesSinceSave && !linkedlist.ChangesMade {
		return nil
	}
	b, err := json.Marshal(l)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(l.Filestore, b, 0644)
	if err != nil {
		return err
	}
	fmt.Println("File saved")
	l.changesSinceSave = false
	linkedlist.ChangesMade = false
	return nil
}

func (l *local) AllLists() []*List {
	l.lock.Lock()
	defer l.lock.Unlock()
	var lists []*List
	for _, v := range l.M {
		lists = append(lists, v)
	}
	// sort by creation date
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Created.Before(lists[j].Created)
	})
	return lists
}

func (l *local) NewList(name string, maxItemsInList int) (*List, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.Total >= l.MaxTotal {
		return nil, fmt.Errorf("already at maximum number of lists allowed")
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("failed to generate uuid: %w", err)
		return nil, err
	}
	newLinkedList := linkedlist.NewList(maxItemsInList)
	newList := &List{
		UUID:    uuid.String(),
		Name:    name,
		List:    &newLinkedList,
		Created: time.Now().UTC().Round(time.Millisecond),
	}
	l.M[uuid.String()] = newList
	l.Total++
	l.changesSinceSave = true
	return newList, nil
}

func (l *local) GetList(uuid string) *List {
	l.lock.Lock()
	defer l.lock.Unlock()
	item := l.M[uuid]
	return item
}

func (l *local) DeleteList(uuid string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.M[uuid]; ok {
		delete(l.M, uuid)
		l.Total--
		l.changesSinceSave = true
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
