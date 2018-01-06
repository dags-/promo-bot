package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type History struct {
	records  map[string]bool
	first    *node
	last     *node
	size     int
	capacity int
}

type node struct {
	id       string
	previous *node
	next     *node
}

func NewHistory(capacity int) History {
	h := History{
		size:     0,
		capacity: capacity,
		records:  make(map[string]bool, capacity),
	}
	h.Read()
	return h
}

func (h *History) Contains(id string) bool {
	_, ok := h.records[id]
	return ok
}

func (h *History) Add(id string) {
	defer h.Write()
	n := node{id: id}
	h.records[id] = true
	if h.first == nil {
		h.first = &n
		h.last = h.first
	} else {
		n.previous = h.last
		h.last.next = &n
		h.last = &n
		h.size++
		if h.size >= h.capacity {
			remove := *h.first
			h.first = remove.next
			delete(h.records, remove.id)
		}
	}
}

func (h *History) Read() {
	data, err := ioutil.ReadFile("history.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	var array []string
	err = json.Unmarshal(data, &array)
	if err != nil {
		fmt.Println(err)
		return
	}

	h.first = nil
	h.last = nil
	for _, id := range array {
		h.Add(id)
	}
}

func (h *History) Write() {
	array := make([]string, h.capacity)

	i := 0
	node := h.first
	for node != nil {
		array[i] = node.id
		node = node.next
		i++
	}

	data, err := json.MarshalIndent(array[:i], "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("history.json", data, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
