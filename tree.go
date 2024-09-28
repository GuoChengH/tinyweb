package tinyweb

import (
	"strings"
)

type treeNode struct {
	name       string
	children   []*treeNode
	routerName string
	isEnd      bool
}

func (t *treeNode) Put(path string) {
	root := t
	strs := strings.Split(path, "/")
	for index, str := range strs {
		if str == "" {
			continue
		}
		isMatch := false
		for _, node := range t.children {
			if node.name == str {
				t = node
				isMatch = true
				break
			}
		}
		if !isMatch {
			isEnd := false
			if index == len(strs)-1 {
				isEnd = true
			}
			newNode := &treeNode{name: str, children: []*treeNode{}, isEnd: isEnd}
			t.children = append(t.children, newNode)
			t = newNode
		}
	}
	t = root
}
func (t *treeNode) Get(path string) *treeNode {
	strs := strings.Split(path, "/")
	routerName := ""
	for index, str := range strs {
		if str == "" {
			continue
		}
		isMatch := false
		for _, node := range t.children {
			// 匹配精确值，* 通配符 或 : 参数
			if node.name == str || node.name == "*" || (len(node.name) > 0 && node.name[0] == ':') {
				t = node
				isMatch = true
				routerName += "/" + node.name
				node.routerName = routerName
				if index == len(strs)-1 {
					return node
				}
				break
			}
		}
		// 没有匹配时处理 **
		if !isMatch {
			for _, node := range t.children {
				if node.name == "**" {
					routerName += "/" + node.name
					node.routerName = routerName
					return node
				}
			}
			return nil
		}
	}
	return nil
}
