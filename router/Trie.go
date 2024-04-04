package router

import (
	"net/http"
	"strings"
)

type handler func(res http.ResponseWriter, req *http.Request, params map[string]string)

type Node struct {
	part       string
	children   map[string]*Node
	isEnd      bool
	handler    handler
	params     map[string]string
	isVariable bool
	method     string
}

type Trie struct {
	root *Node
}

func NewNode(part string, isVariable bool) *Node {
	return &Node{
		part:       part,
		children:   make(map[string]*Node),
		isEnd:      false,
		isVariable: isVariable,
	}
}

func NewRouter() *Trie {
	return &Trie{
		root: NewNode("", false),
	}
}

func (t *Trie) Insert(route string, handler handler, method string) {
	if route == "/" {
		route = "1st"
	}

	node := t.root
	parts := strings.Split(route, "/")[1:]

	for _, part := range parts {
		if _, ok := node.children[part]; !ok {
			var variable bool = strings.HasPrefix(part, ":")
			if variable {
				part = part[1:]
			}
			node.children[part] = NewNode(part, variable)
		}
		node = node.children[part]
	}
	node.isEnd = true
	node.handler = handler
	node.params = make(map[string]string)
	node.method = method
}

func (t *Trie) Search(route string) (*Node, bool) {
	if route == "/" {
		route = "1st"
	}

	node := t.root
	parts := strings.Split(route, "/")[1:]
	for _, part := range parts {
		if _, ok := node.children[part]; !ok {

			var found bool

			for _, child := range node.children {
				if child.isVariable {
					node = child
					node.params[node.part] = part
					found = true
					break
				}
			}
			if !found {
				return nil, false
			}
		} else {
			node = node.children[part]
		}
	}
	return node, node.isEnd
}
