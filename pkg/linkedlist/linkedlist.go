package linkedlist

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrEmptyUUID = fmt.Errorf("empty UUID is not supported")
	ErrDupUUID   = fmt.Errorf("duplicate UUIDs cannot be in the same list")
)

// List is a double linked list
type List struct {
	first *node
	last  *node
	nodes map[string]*node
	lock  *sync.Mutex
}

type Item struct {
	UUID string `json:"UUID"`
	Text string `json:"value"`
}

type node struct {
	item     Item
	previous *node
	next     *node
}

func NewList() List {
	return List{
		nodes: make(map[string]*node),
		lock:  &sync.Mutex{},
	}
}

// AddItem will add an item to the bottom of the list
func (l *List) AddItem(entry Item) string {
	if entry.UUID == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			fmt.Println("uuid generation failed!?!")
			return ""
		}
		entry.UUID = uuid.String()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	n := node{
		item: entry,
	}
	l.nodes[entry.UUID] = &n
	if l.first == nil {
		l.first = &n
		l.last = &n
		return entry.UUID
	}
	n.previous = l.last
	l.last.next = &n
	l.last = &n
	return entry.UUID
}

// EditItem is used to update the text inside an item
func (l *List) EditItem(name string, value string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	n, ok := l.nodes[name]
	if !ok {
		fmt.Printf("cannot edit '%v' because it does not exist", name)
		return
	}
	n.item.Text = value
}

// DeleteItem is used to delete an item from the list
func (l *List) DeleteItem(name string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	n, ok := l.nodes[name]
	if !ok {
		return
	}
	n1Prev := n.previous
	n1Next := n.next
	if n1Prev == nil {
		l.first = n1Next
	} else {
		n1Prev.next = n1Next
	}
	delete(l.nodes, name)
	n = nil
}

// MoveItemAfter takes in the name of the item wanting to be moved and the location where it will be moved after
func (l *List) MoveItemAfter(name string, location string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if name == location {
		return
	}
	n1, ok := l.nodes[name]
	if !ok {
		fmt.Printf("name '%v' is not in the list!\n", name)
		return
	}
	n2, ok := l.nodes[location]
	if !ok {
		fmt.Printf("location '%v' is not in the list!\n", name)
		return
	}
	// remove n1 completely from the list
	if n1.previous == nil {
		l.first = n1.next
	} else {
		n1.previous.next = n1.next
	}
	if n1.next == nil {
		l.last = n1.previous
	} else {
		n1.next.previous = n1.previous
	}
	// add n1 after n2
	if n2.next == nil {
		l.last = n1
		n1.next = nil
	} else {
		n2.next.previous, n1.next = n1, n2.next
	}
	n1.previous = n2
	n2.next = n1
}

// MoveItemBefore takes in the name of the item wanting to be moved and the location where it will be moved in front of
func (l *List) MoveItemBefore(name string, location string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if name == location {
		return
	}
	n1, ok := l.nodes[name]
	if !ok {
		fmt.Printf("name '%v' is not in the list!\n", name)
		return
	}
	n2, ok := l.nodes[location]
	if !ok {
		fmt.Printf("location '%v' is not in the list!\n", name)
		return
	}
	// remove n1 completely from the list
	if n1.previous == nil {
		l.first = n1.next
	} else {
		n1.previous.next = n1.next
	}
	if n1.next == nil {
		l.last = n1.previous
	} else {
		n1.next.previous = n1.previous
	}
	// add n1 before n2
	if n2.previous == nil {
		l.first = n1
		n1.previous = nil
	} else {
		n2.previous.next, n1.previous = n1, n2.previous
	}
	n1.next = n2
	n2.previous = n1
}

// PrintAll is used to print out a list of items
func (l *List) PrintAll() {
	l.lock.Lock()
	defer l.lock.Unlock()
	node := l.first
	for node != nil {
		fmt.Println(node.item.Text)
		node = node.next
	}
}

// ListAll is used to retrieve a list of all items
func (l *List) ListAll() []Item {
	var items []Item
	l.lock.Lock()
	defer l.lock.Unlock()
	node := l.first
	for node != nil {
		items = append(items, node.item)
		node = node.next
	}
	return items
}

// ListAllReverse is used to retrieve a list of all items in the reverse order
func (l *List) ListAllReverse() []Item {
	var items []Item
	l.lock.Lock()
	defer l.lock.Unlock()
	node := l.last
	for node != nil {
		items = append(items, node.item)
		node = node.previous
	}
	return items
}

// MarshalJSON is used to create a custom JSON payload for list
func (l List) MarshalJSON() ([]byte, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	var values []Item
	node := l.first
	for node != nil {
		v := Item{
			UUID: node.item.UUID,
			Text: node.item.Text,
		}
		values = append(values, v)
		node = node.next
	}
	if len(values) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(values)
}

// UnmarshalJSON is used to generate a list from a json payload
func (l *List) UnmarshalJSON(b []byte) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	var values []Item
	err := json.Unmarshal(b, &values)
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return nil
	}
	n := &node{item: values[0]}
	if n.item.UUID == "" {
		return ErrEmptyUUID
	}

	l.nodes[n.item.UUID] = n
	l.first = n
	nPrevious := n
	// skip the first entry
	for i := 1; i < len(values); i++ {
		nNew := &node{item: values[i]}
		nNew.previous = nPrevious
		if nNew.item.UUID == "" {
			return ErrEmptyUUID
		}
		if _, ok := l.nodes[nNew.item.UUID]; ok {
			return ErrDupUUID
		}
		l.nodes[nNew.item.UUID] = nNew
		n.next = nNew
		n = nNew
		nPrevious = nNew
		l.last = nNew
	}
	return nil
}
