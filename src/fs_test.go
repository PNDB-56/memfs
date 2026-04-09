package main

import (
	"slices"
	"testing"
)

func setupRoot() *Node {
	r := NewRoot()
	return &r
}
func TestPwd_Root(t *testing.T) {
	root := setupRoot()
	if root.Pwd() != "/root" {
		t.Error(t.Name(), ": Root folder initialization failed")
	}
}

func TestPwd_Nested(t *testing.T) {
	root := setupRoot()
	dirA := &Node{isRoot: false, kind: "dir", name: "a", children: make([]*Node, 0), parent: root, index: make(map[string]*Node)}
	root.children = append(root.children, dirA)
	root.index["a"] = dirA
	root = root.children[0] // now at /root/a
	if root.Pwd() != "/root/a" {
		t.Error(t.Name(), ": pwd is not at dir \"a\"")
	}
	dirB := &Node{isRoot: false, kind: "dir", name: "b", children: make([]*Node, 0), parent: root, index: make(map[string]*Node)} // order of Node is important due to parent node
	root.children = append(root.children, dirB)
	root.index["b"] = dirB
	root.children = append(root.children, dirB)
	root = root.children[0] // now at /root/b
	if root.Pwd() != "/root/a/b" {
		t.Error(t.Name(), ": pwd is not at dir \"root/a/b\"")
	}
}

func TestLs(t *testing.T) {
	root := setupRoot()
	dirA := &Node{isRoot: false, kind: "dir", name: "a", children: make([]*Node, 0), parent: root, index: make(map[string]*Node)}
	dirB := &Node{isRoot: false, kind: "dir", name: "b", children: make([]*Node, 0), parent: root, index: make(map[string]*Node)}
	dirC := &Node{isRoot: false, kind: "dir", name: "c", children: make([]*Node, 0), parent: dirA, index: make(map[string]*Node)}
	dirA.children = append(dirA.children, dirC)
	dirA.index["c"] = dirC
	dirD := &Node{isRoot: false, kind: "dir", name: "d", children: make([]*Node, 0), parent: dirB, index: make(map[string]*Node)}
	dirB.children = append(dirB.children, dirD)
	dirB.index["d"] = dirD
	root.children = append(root.children, dirA)
	root.children = append(root.children, dirB)
	root.index["a"] = dirA
	root.index["b"] = dirB
	lsResult := root.Ls()
	if !slices.Contains(lsResult, "/a") {
		t.Error(t.Name(), ": dir \"a\" is missing in /root")
	}
	if !slices.Contains(lsResult, "/b") {
		t.Error(t.Name(), ": dir \"b\" is missing in /root")
	}
	lsResult = dirA.Ls()
	if !slices.Contains(lsResult, "/c") {
		t.Error(t.Name(), ": dir \"c\" is missing in /root/a")
	}
	lsResult = dirB.Ls()
	if !slices.Contains(lsResult, "/d") {
		t.Error(t.Name(), ": dir \"d\" is missing in /root/b")
	}
}
