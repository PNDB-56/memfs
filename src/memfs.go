package main

import "fmt"

func main() {
	root := NewRoot()
	pwd := &root
	log(pwd.Pwd())
	pwd.Mkdir("test1")
	log(pwd.Ls())
	pwd.Mkdir("test2")
	// testSetup1(pwd)
	log(pwd.Ls())
}

func testSetup(ctx *Node) {
	ctx.children = append(ctx.children, &Node{isRoot: false, kind: "dir", name: "a", children: make([]*Node, 0), parent: ctx})
	ctx = ctx.children[0]
	ctx.children = append(ctx.children, &Node{isRoot: false, kind: "dir", name: "b", children: make([]*Node, 0), parent: ctx})
	ctx = ctx.children[0]
}

func testSetup1(ctx *Node) {
	ctx.children = append(ctx.children, &Node{isRoot: false, kind: "dir", name: "a", children: make([]*Node, 0), parent: ctx})
	ctx.children = append(ctx.children, &Node{isRoot: false, kind: "dir", name: "b", children: make([]*Node, 0), parent: ctx})
	ctx.children = append(ctx.children, &Node{isRoot: false, kind: "file", name: "c.txt", children: nil, parent: ctx})
	ctx = ctx.children[1]
	ctx.children = append(ctx.children, &Node{isRoot: false, kind: "file", name: "d.txt", children: nil, parent: ctx})
}

func log[T any](args ...T) {
	fmt.Println(args)
}
