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
	// ChangesMade is a global way to tell if a change was made to any list
	ChangesMade bool
)

// List is a double linked list
type List struct {
	first         *node
	last          *node
	nodes         map[string]*node
	lock          *sync.Mutex
	MaxTotal      int
	Total         int
	MaxTextLength int
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

func NewList(maxItems int) List {
	return List{
		nodes:         make(map[string]*node),
		MaxTotal:      maxItems,
		lock:          &sync.Mutex{},
		MaxTextLength: 1000,
	}
}

// AddItem will add an item to the bottom of the list
func (l *List) AddItem(entry Item) (string, error) {
	if l.Total >= l.MaxTotal {
		return "", fmt.Errorf("already at total number of items in list")
	}
	if len(entry.Text) > l.MaxTextLength {
		return "", fmt.Errorf("entry text is too long")
	}
	if entry.UUID == "" {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return "", fmt.Errorf("uuid generation failed: %w", err)
		}
		entry.UUID = uuid.String()
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Total++
	ChangesMade = true
	n := node{
		item: entry,
	}
	l.nodes[entry.UUID] = &n
	if l.first == nil {
		l.first = &n
		l.last = &n
		return entry.UUID, nil
	}
	n.previous = l.last
	l.last.next = &n
	l.last = &n
	return entry.UUID, nil
}

// EditItem is used to update the text inside an item
func (l *List) EditItem(name string, value string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if len(value) > l.MaxTextLength {
		ln("Can't make item length that long")
		return
	}
	n, ok := l.nodes[name]
	if !ok {
		fmt.Printf("cannot edit '%v' because it does not exist", name)
		return
	}
	n.item.Text = value
	ChangesMade = true
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
	l.Total--
	n = nil
	ChangesMade = true
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
	ChangesMade = true
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
	ChangesMade = true
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
	type item struct {
		Items         []Item
		MaxTotal      int
		MaxTextLength int
		Total         int
	}
	var values item
	node := l.first
	for node != nil {
		v := Item{
			UUID: node.item.UUID,
			Text: node.item.Text,
		}
		values.Items = append(values.Items, v)
		node = node.next
	}
	values.MaxTotal = l.MaxTotal
	values.MaxTextLength = l.MaxTextLength
	values.Total = l.Total
	if len(values.Items) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(values)
}

// UnmarshalJSON is used to generate a list from a json payload
func (l *List) UnmarshalJSON(b []byte) error {
	if l.lock == nil {
		l.lock = &sync.Mutex{}
		l.nodes = make(map[string]*node)
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	type item struct {
		Items         []Item
		MaxTotal      int
		MaxTextLength int
	}
	var values item
	l.MaxTotal = values.MaxTotal
	l.MaxTextLength = values.MaxTextLength
	l.Total = len(values.Items)
	err := json.Unmarshal(b, &values)
	if err != nil {
		return err
	}
	if len(values.Items) == 0 {
		return nil
	}
	n := &node{item: values.Items[0]}
	if n.item.UUID == "" {
		return ErrEmptyUUID
	}
	l.nodes[n.item.UUID] = n
	l.first = n
	nPrevious := n
	// skip the first entry
	for i := 1; i < len(values.Items); i++ {
		nNew := &node{item: values.Items[i]}
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
