/*
  The container/tree package implements tree container
  this package does not serve binary-tree but ordered-tree.
  Node struct can be used as XMLNode etc...

*/
package tree

import (
	"container"
	"errors"
	"fmt"
	"log"
)

//Error
var (
	ERROR_REMOVE_CHILD = errors.New("You tried to remove a Node that is not in Node's children")
	ERROR_NOT_CONTAINS = errors.New("NodeList does not contain the Node")
)

//NodeList hanle Node's children
type NodeList []*Node

//NodePath is slice of pointers from root to a node.
type NodePath []*Node

//Index(node) returns a index number of the node in the Nodelist
//if the node is not in the NodeList, this method returns error
func (nl NodeList) Index(node *Node) (index int, err error) {
	for i, n := range nl {
		if n == node {
			index = i
			return
		}
	}
	index, err = -1, ERROR_NOT_CONTAINS
	return
}

//Swap(i, j) swaps nodelist[i] and nodelist[j]
func (nl NodeList) Swap(i, j int) {
	nl[i], nl[j] = nl[j], nl[i]
}

//SwapNodes(node1, node2) swaps the indice of node1 and node2
//If eather node1 or node2 is not in the NodeList, this method returns error
func (nl NodeList) SwapNodes(n1, n2 *Node) (err error) {
	i1, err1 := nl.Index(n1)
	i2, err2 := nl.Index(n2)
	if err1 != nil {
		err = err1
	} else if err2 != nil {
		err = err2
	} else {
		nl.Swap(i1, i2)
	}
	return
}

//Tree Node have its parent's pointer ,NodeList of its children and its Value
type Node struct {
	parent   *Node
	children NodeList
	Value    Value
}

//NewNode function creates a Node and returns a pointer of the Node
func NewNode(contents Value) *Node {
	node := new(Node)
	node.Value = &contents
	return node
}

func (node *Node) Parent() *Node {
	return node.parent
}

func (node *Node) Children() NodeList {
	return node.children
}

//Path() returns a slice of pointers from the root node to the node
func (node *Node) Path() NodePath {
	path := make(NodePath, 0, 0)
	path = append(path, node)
	next := node
	for {
		next = next.Parent()
		log.Println(next)
		if next == nil {
			break
		}
		path = append(path, next)

	}
	length := len(path)
	result := make(NodePath, length, length)
	for i, n := range path {
		result[length-i-1] = n
	}
	return result
}

//Append a node to its children
//If the new child is in another node's children,
//it is removed from the previous parent's children
func (node *Node) AppendChild(child *Node) (err error) {
	node.children = append(node.children, child)
	if parent := node.Parent(); parent != nil {
		parent.RemoveChild(child)
	}
	child.parent = node
	return
}

//RemoveChild remove a child from the node's children
//If the child is not in the children, its returns error.
func (node *Node) RemoveChild(child *Node) (err error) {
	capa := len(node.children)
	newNodeList := make(NodeList, 0, capa)
	for _, childnode := range node.children {
		if childnode != child {
			newNodeList = append(newNodeList, child)
			child.parent = nil
		}
	}
	if len(newNodeList) == capa {
		err = ERROR_REMOVE_CHILD
	} else {
		node.children = newNodeList
	}
	return
}

/*
Walk() is a generator.
It returns all nodes under the node.
for node := range parent.Walk(){
  //do something
}
*/
func (node *Node) Walk() (ch chan *Node) {
	f := func(node *Node) bool { return true }
	return node.Filter(f)
}

/*
Filter(f) is a generator.
It returns nodes that funcion returns true 
*/
func (node *Node) Filter(f func(*Node) bool) (ch chan *Node) {
	ch = make(chan *Node)
	go func() {
		for _, child := range node.children {
			if f(child) {
				ch <- child
			}
			for gchild := range child.Walk() {
				if f(gchild) {
					ch <- gchild
				}
			}
		}
		close(ch)
	}()
	return
}

func (node *Node) String() string {
	return fmt.Sprintf("%v", node.Value)
}
