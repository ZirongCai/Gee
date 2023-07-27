package gee

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang ,只在最后一个node有值
	part     string  // 该node的值，路由中的一部分，例如 :lang
	children []*node //子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) Insert(pattern string) {
	parts := ParsePattern(pattern)
	n.insertRecursive(pattern, parts, 0)
}

// param如果只有pattern和height会不会更好？
func (n *node) insertRecursive(pattern string, parts []string, height int) {
	// parts忽略了根节点，从第一个路由地址开始
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)

	// 每个node只代表一个part，所以就算append了一个node也没完，还得往下走
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insertRecursive(pattern, parts, height+1)
}

// 找到node并返回
func (n *node) Search(path string) *node {
	return n.searchRecursive(ParsePattern(path), 0)
}

func (n *node) searchRecursive(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part) //因为有pattern的存在，所以可能存在多条路径，但match的只能有一个

	for _, child := range children {
		if result := child.searchRecursive(parts, height+1); result != nil {
			return result
		}
	}

	return nil
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
