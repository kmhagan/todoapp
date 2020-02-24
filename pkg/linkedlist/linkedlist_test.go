package linkedlist

import (
	"encoding/json"
	"fmt"
	"testing"
)

func createTestList(count int) List {
	l := NewList()
	for i := 0; i < count; i++ {
		entry := Item{
			UUID: fmt.Sprintf("%d", i),
			Text: fmt.Sprintf("%d", i),
		}
		l.AddItem(entry)
	}
	return l
}

func TestMoveItemAfter(t *testing.T) {
	l := createTestList(7)
	l.MoveItemAfter("4", "2")
	list := l.ListAll()
	expected := []string{
		"0",
		"1",
		"2",
		"4",
		"3",
		"5",
		"6",
	}
	if len(list) != len(expected) {
		t.Fatalf("Expected a total of %d items but had %d", len(expected), len(list))
	}
	for i := range list {
		if list[i].Text != expected[i] {
			t.Errorf("Expected '%v' but got '%v'", expected[i], list[i].Text)
		}
	}
	list = l.ListAllReverse()
	if len(list) != len(expected) {
		t.Fatalf("Expected a total of %d items but had %d", len(expected), len(list))
	}
	for i := range list {
		j := len(expected) - i - 1
		if list[i].Text != expected[j] {
			t.Errorf("Reverse search failed expected '%v' but got '%v'", expected[j], list[i].Text)
		}
	}
}

func TestMoveItemBefore(t *testing.T) {
	l := createTestList(7)
	l.MoveItemBefore("4", "2")
	list := l.ListAll()
	expected := []string{
		"0",
		"1",
		"4",
		"2",
		"3",
		"5",
		"6",
	}
	if len(list) != len(expected) {
		t.Fatalf("Expected a total of %d items but had %d. List: %v", len(expected), len(list), list)
	}
	for i := range list {
		if list[i].Text != expected[i] {
			t.Errorf("Expected '%v' but got '%v'", expected[i], list[i].Text)
		}
	}
	list = l.ListAllReverse()
	if len(list) != len(expected) {
		t.Fatalf("Expected a total of %d items but had %d. List: %v", len(expected), len(list), list)
	}
	for i := range list {
		j := len(expected) - i - 1
		if list[i].Text != expected[j] {
			t.Errorf("Reverse search failed expected '%v' but got '%v'", expected[j], list[i].Text)
		}
	}
}

func TestMoveItemAfterLastItem(t *testing.T) {
	l := createTestList(7)
	l.MoveItemAfter("4", "6")
	list := l.ListAll()
	expected := []string{
		"0",
		"1",
		"2",
		"3",
		"5",
		"6",
		"4",
	}
	if len(list) != len(expected) {
		t.Fatalf("Expected a total of %d items but had %d. List: %v", len(expected), len(list), list)
	}
	for i := range list {
		if list[i].Text != expected[i] {
			t.Errorf("Expected '%v' but got '%v'", expected[i], list[i].Text)
		}
	}
	list = l.ListAllReverse()
	for i := range list {
		j := len(expected) - i - 1
		if list[i].Text != expected[j] {
			t.Errorf("Reverse search failed expected '%v' but got '%v'", expected[j], list[i].Text)
		}
	}
}

func TestMoveItemBeforeFirstItem(t *testing.T) {
	l := createTestList(7)
	l.MoveItemBefore("4", "0")
	list := l.ListAll()
	expected := []string{
		"4",
		"0",
		"1",
		"2",
		"3",
		"5",
		"6",
	}
	if len(list) != len(expected) {
		t.Fatalf("Expected a total of %d items but had %d. List: %v", len(expected), len(list), list)
	}
	for i := range list {
		if list[i].Text != expected[i] {
			t.Errorf("Expected '%v' but got '%v'", expected[i], list[i].Text)
		}
	}
	list = l.ListAllReverse()
	for i := range list {
		j := len(expected) - i - 1
		if list[i].Text != expected[j] {
			t.Errorf("Reverse search failed expected '%v' but got '%v'", expected[j], list[i].Text)
		}
	}
}

func TestJSON(t *testing.T) {
	l := createTestList(4)
	b, err := json.Marshal(l)
	if err != nil {
		t.Fatalf("failed to marshal list")
	}
	expectedB := `[{"UUID":"0","value":"0"},{"UUID":"1","value":"1"},{"UUID":"2","value":"2"},{"UUID":"3","value":"3"}]`
	if string(b) != expectedB {
		t.Fatalf("expected json `%s` instead got %s", expectedB, b)
	}

	l = NewList()
	err = json.Unmarshal(b, &l)
	all := l.ListAll()
	if len(all) != 4 {
		t.Errorf("expected 4 in ListAll but got %d", len(all))
	}
	for i, item := range all {
		if item.Text != fmt.Sprintf("%d", i) || item.UUID != fmt.Sprintf("%d", i) {
			t.Errorf("%v is wrong item for %d", item, i)
		}
	}

	all = l.ListAllReverse()
	for i, item := range all {
		j := len(all) - i - 1
		if item.Text != fmt.Sprintf("%d", j) || item.UUID != fmt.Sprintf("%d", j) {
			t.Errorf("%v is wrong item for %d", item, j)
		}
	}
}
