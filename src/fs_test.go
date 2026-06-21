package main

import (
	"io"
	"slices"
	"testing"
)

func TestPwd_Root(t *testing.T) {
	root := NewRoot()
	if root.Pwd() != "/root" {
		t.Error(t.Name(), ": Root folder initialization failed")
	}
}

func TestPwd_Nested(t *testing.T) {
	root := NewRoot()
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
	root := NewRoot()
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

func TestStat(t *testing.T) {
	root := NewRoot()
	root.Mkdir("A")
	root.Cd("./A").Mkdir("B")
	root.Cd("./A/B").Mkdir("C")
	exists, parent, curr, err := root.Stat(root, ".//A/B/C")
	if !exists || parent == nil || curr == nil || err != nil {
		t.Error(t.Name(), ": Stat failed - C exists")
	}

	exists, parent, curr, err = root.Stat(root, ".//A/B/D")
	if exists || parent != nil || curr != nil || err == nil {
		t.Error(t.Name(), ": Stat failed - D doens't exist")
	}
}

func TestTouch(t *testing.T) {
	root := NewRoot()
	err := root.Touch("a.txt")
	if err != nil {
		t.Error(t.Name(), ": Touch failed - Unable to create a file")
	}
	ls := root.Ls()
	if ls[0] != "a.txt" {
		t.Error(t.Name(), ": Touch failed - a.txt doesn't exist")
	}
	err = root.Touch("a.txt")
	if err == nil || err.Error() != "File already exists" {
		t.Error(t.Name(), ": Failed to prevent duplicate file creation")
	}
}

func TestWrite(t *testing.T) {
	root := NewRoot()
	err := root.Touch("a.txt")
	if err != nil {
		t.Error(t.Name(), err.Error())
	}
	fileContent := "abcd"
	bytesWritten, err := root.index["a.txt"].file.Write([]byte(fileContent))
	if err != nil {
		t.Error(t.Name(), err.Error())
	}
	if bytesWritten != len(fileContent) {
		t.Error(t.Name(), ": Written bytes size mismatch")
	}
	by := byte('a')
	// root.index["a.txt"].file.data[1] = 'e'
	// root.index["a.txt"].file.data[2] = 'e'
	// root.index["a.txt"].file.data[3] = 'e'
	for index, x := range root.index["a.txt"].file.data {
		if x != by {
			t.Errorf("%s : Expected byte %d at index %d, but found %d", t.Name(), by, index, x)
		}
		by += 1
	}
}

func TestRead(t *testing.T) {
	root := NewRoot()
	err := root.Touch("a.txt")
	if err != nil {
		t.Error(t.Name(), err.Error())
	}
	fileContent := "abcdefghijkl"
	bytesWritten, err := root.index["a.txt"].file.Write([]byte(fileContent))
	if err != nil {
		t.Error(t.Name(), err.Error())
	}
	if bytesWritten != len(fileContent) {
		t.Error(t.Name(), ": Mock write - Written bytes size mismatch")
	}
	read5 := make([]byte, 5)

	readBytes, err := root.index["a.txt"].file.Read(read5)
	if err != nil && err != io.EOF {
		t.Error(t.Name(), err.Error())
	}
	if readBytes != 5 {
		t.Error(t.Name(), ": Expected to read 5 bytes but read:", readBytes, "bytes")
	}

	by := byte('a')
	for index, x := range read5 {
		if x != by {
			t.Errorf("%s : Expected to read byte %d at index %d but found byte %d", t.Name(), by, index, x)
		}
		by += 1
	}

	read12 := make([]byte, 12)
	// root.index["a.txt"].file.data[1] = 'e'
	// root.index["a.txt"].file.data[2] = 'e'
	// root.index["a.txt"].file.data[3] = 'f'
	readBytes, err = root.index["a.txt"].file.Read(read12)
	// fmt.Println(readBytes, err)
	if err != nil && err != io.EOF {
		t.Error(t.Name(), err.Error())
	}
	if readBytes != 12 {
		t.Error(t.Name(), ": Expected to read 12 bytes but read:", readBytes, "bytes")
	}

	by = byte('a')
	for index, x := range read12 {
		if x != by {
			t.Errorf("%s : Expected to read byte %d at index %d but found byte %d", t.Name(), by, index, x)
		}
		by += 1
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "a", false},
		{"valid file", "a.txt", false},
		{"valid hyphen underscore", "hello_world-123", false},
		{"invalid empty", "", true},
		{"invalid dot", ".", true},
		{"invalid dotdot", "..", true},
		{"invalid slash", "a/b", true},
		{"invalid char tilde", "a~b", true},
		{"invalid char hash", "a#b", true},
		{"invalid char asterisk", "a*b", true},
		{"invalid char backslash", "a\\b", true},
		{"invalid char semicolon", "a;b", true},
		{"invalid char quote", "a\"b", true},
		{"invalid too long", string(make([]byte, 256)), true},
		{"invalid control char", "a\x01b", true},
		{"invalid del char", "a\x7fb", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateName(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid root slash", "/", false},
		{"valid relative simple", "a/b", false},
		{"valid absolute simple", "/a/b", false},
		{"valid trailing slash", "/a/b/", false},
		{"valid navigation dot", "a/./b", false},
		{"valid navigation dotdot", "a/../b", false},
		{"valid double slash", "a//b", false},
		{"invalid empty", "", true},
		{"invalid component hash", "a/b#c/d", true},
		{"invalid component semicolon", "/a//b/x/../y/z;1", true},
		{"invalid component asterisk", "/a/../b*c", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePath(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

