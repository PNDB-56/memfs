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

// TODO: Mkdir creates a dir in current dir only, It should accept path to create a dir (abs path , relative path)
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

// TODO: Ls works for current dir, should accept path (abs path , relative path) to query children, may be supported flags as well ?
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

// TODO: Invalid paths should throw proper error, add validations for input path
// TODO: should support abs path as well
func (f *Node) Cd(path string) *Node {
	orginalPwd := f
	subPaths := strings.FieldsFunc(path, func(c rune) bool { return c == '/' })
	for _, p := range subPaths {
		if p == "." {
			continue
		} else if p == ".." {
			// go to parent
			if f.parent != nil {
				f = f.parent
			} else {
				// TODO: improve Error
				fmt.Printf("Parent Path: %s doesn't exist in %s\n", p, path)
				return orginalPwd
			}
		} else {
			if addr, ok := f.index[p]; ok {
				f = addr
			} else {
				// TODO: improve Error
				fmt.Printf("Child Path: %s doesn't exist in %s\n", p, path)
				return orginalPwd
			}
		}
	}
	return f
}
