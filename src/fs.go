package main

import (
	"errors"
	"fmt"
	"io"
	"path"
	"slices"
	"strings"
	"time"
)

type File struct {
	fileName   string
	createdAt  time.Time
	modifiedAt time.Time
	data       []byte
	open       bool
}

type Node struct {
	isRoot   bool
	kind     string
	name     string
	children []*Node
	index    map[string]*Node
	parent   *Node
	file     *File
}

func NewRoot() *Node {
	return &Node{
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
	slices.Reverse(fullPath)
	return strings.Join(fullPath, "/")
}

// TODO: Mkdir creates a dir in current dir only, It should accept path to create a dir (abs path , relative path)
// TODO: optimization of childer array to binary tree or trie for easy search
func (f *Node) Mkdir(name string) (bool, error) {
	if f == nil {
		return false, errors.New("current context is nil")
	}
	_, ok := f.index[name]
	if ok {
		return false, errors.New("Dir already exists")
	}
	node := Node{isRoot: false, kind: "dir", name: name, children: make([]*Node, 0), parent: f, index: make(map[string]*Node)}
	f.children = append(f.children, &node)
	f.index[name] = &node
	return true, nil
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
// TODO: if cd is given for file name instead of DIr throw err
func (f *Node) Cd(p string) *Node {
	orginalPwd := f
	cleanedPath := path.Clean(p)
	// fmt.Println(cleanedPath, p)
	subPaths := strings.FieldsFunc(cleanedPath, func(c rune) bool { return c == '/' })
	for _, p1 := range subPaths {
		if p1 == "." {
			continue
		} else if p1 == ".." {
			// go to parent
			if f.parent != nil {
				f = f.parent
			} else {
				// TODO: improve Error
				fmt.Printf("Parent Path: %s doesn't exist in %s\n", p1, p)
				return orginalPwd
			}
		} else {
			if addr, ok := f.index[p1]; ok {
				f = addr
			} else {
				// TODO: improve Error
				fmt.Printf("Child Path: %s doesn't exist in %s\n", p1, p)
				return orginalPwd
			}
		}
	}
	return f
}

func (f *Node) Stat(root *Node, path string) (exists bool, parent *Node, node *Node, err error) {
	exists = false
	parent = nil
	node = nil
	err = nil
	if f == nil {
		err = errors.New("current context is nil")
		return
	}
	pathArr := strings.FieldsFunc(path, func(c rune) bool { return c == '/' })
	// fmt.Println(pathArr)
	if len(pathArr) == 0 {
		err = errors.New("stat path can't be empty")
		return
	}
	startNode := f
	for _, x := range pathArr {
		switch x {
		case ".":
			startNode = f
		case "..":
			if f.parent == nil {
				err = fmt.Errorf("Invalid path at %s in %q", x, path)
				return
			}
			startNode = f.parent
		case "root":
			startNode = root
		default:
			if childNode, ok := startNode.index[x]; ok {
				startNode = childNode
			} else {
				err = fmt.Errorf("Invalid path at %s in %q", x, path)
				return
			}
		}

	}
	parent = startNode.parent
	node = startNode
	exists = true
	return
}

// TODO: add validations to fromPath, toPath
// func (f *Node) Move(root *Node, fromPath string, toPath string) (bool, error) {

// }

// TODO: add path validations for name, eg: ./a.txt = a.txt and support relative and abs paths in name
func (f *Node) Touch(name string) error {
	if f == nil {
		return errors.New("current context is nil")
	}
	_, ok := f.index[name]
	if ok {
		return errors.New("File already exists")
	}
	node := Node{
		isRoot:   false,
		kind:     "file",
		name:     name,
		children: nil,
		parent:   f,
		index:    nil,
		file: &File{
			fileName:   name,
			createdAt:  time.Now(),
			modifiedAt: time.Now(),
			data:       make([]byte, 0, 100),
			open:       false}}
	f.children = append(f.children, &node)
	f.index[name] = &node
	return nil
}

func (f *File) Open() error {
	if f.open {
		return errors.New("File is already open")
	}
	f.open = true
	return nil
}

func (f *File) Close() error {
	if f.open {
		f.open = false
		return nil
	}
	return errors.New("File is already closed")
}

// TODO: Come up with some ways to implement open close methods
// TODO: SOME mechanism to remember curson if full data is not read
func (f *File) Read(p []byte) (n int, err error) {
	destLength := len(p)
	sourceLength := len(f.data)
	n = 0
	for n < sourceLength {
		if n < destLength {
			p[n] = f.data[n]
		} else {
			return
		}
		n += 1
	}
	err = io.EOF
	return
}

// TODO: write modes, append / fresh write
// Write will override for now
func (f *File) Write(p []byte) (n int, err error) {
	n = len(p)
	f.data = make([]byte, n)
	err = nil
	// f.data = append(f.data, p...)
	copy(f.data, p)
	return
}

func (f *Node) Write(filePath string, data string) {

}
