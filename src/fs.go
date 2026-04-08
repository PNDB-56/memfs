package main

import (
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	isRoot   bool
	kind     string
	name     string
	children []*Node
	index    map[string]*Node
	parent   *Node
}

func NewRoot() Node {
	return Node{
		isRoot:   true,
		kind:     "dir",
		name:     "root",
		children: make([]*Node, 0),
		index:    make(map[string]*Node),
		parent:   nil,
	}
}

func (f *Node) Pwd() string {
	fullPath := make([]string, 0, 10)
	f1 := f
	for !f1.isRoot {
		fullPath = append(fullPath, f1.name)
		f1 = f1.parent
	}
	fullPath = append(fullPath, "/root")
	// fmt.Println(fullPath)
	reverseSlice(&fullPath)
	// fmt.Println(fullPath)
	return strings.Join(fullPath, "/")
}

func reverseSlice[T any](s *[]T) {
	n := len(*s)
	for i := n - 1; i >= n/2; i -= 1 {
		t := (*s)[i]
		(*s)[i] = (*s)[n-1-i]
		(*s)[n-1-i] = t
	}
}

func (f *Node) Mkdir(name string) (bool, error) {
	if f != nil {
		_, ok := f.index[name]
		if !ok {
			node := Node{isRoot: false, kind: "dir", name: name, children: make([]*Node, 0), parent: f, index: make(map[string]*Node)}
			f.children = append(f.children, &node)
			f.index[name] = &node
			return true, nil
		} else {
			return false, errors.New("Dir already exists")
		}
	} else {
		return false, errors.New("current context is nil")
	}
}

func (f *Node) Ls() []string {
	dirs := make([]string, 0, len(f.children))
	for _, x := range f.children {
		if x.kind == "dir" {
			dirs = append(dirs, fmt.Sprintf("/%s", x.name))
		} else {
			dirs = append(dirs, x.name)

		}
	}
	return dirs
}
