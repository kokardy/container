/*
  The container package implements tree container
  this package does not serve binary-tree but ordered-tree.
  Node struct can be used as XMLNode etc...

*/
package container

import (
	"errors"
	"fmt"
	"log"
)

//Error
var (
	ERROR_REMOVE_CHILD = errors.New("You tried to remove a Node that is not in Node's children")
	ERROR_NOT_CONTAINS = errors.New("NodeList does not contain the Node")
	ERROR_LAST_CHILD   = errors.New("No younger sibling node")
)

//NodePath is slice of pointers from root to a node.
type NodePath []*Node

//Children hanle Node's children
type Children struct {
	Nodes      []*Node
	node_index map[*Node]int
}

func NewChildren() *Children {
	children := new(Children)
	children.node_index = make(map[*Node]int)
	return children
}

//
func (c *Children) Append(node *Node) {
	_, ok := c.node_index[node]
	if ok {
		c.Remove(node)
	}
	c.Nodes = append(c.Nodes, node)
	c.node_index[node] = len(c.Nodes) - 1
}

//
func (c *Children) Remove(node *Node) (err error) {
	index, ok := c.node_index[node]
	if !ok {
		return ERROR_REMOVE_CHILD
	}
	new_list := make([]*Node, 0, 0)
	for _, n := range c.Nodes[0:index] {
		new_list = append(new_list, n)
	}
	for _, n := range c.Nodes[index+1:] {
		new_list = append(new_list, n)
	}
	c.Nodes = new_list
	delete(c.node_index, node)
	return
}

//Index(node) returns a index number of the node in the children
//if the node is not in the children, this method returns error
func (c *Children) Index(node *Node) (index int, err error) {
	log.Println("index:", c.node_index)
	index, ok := c.node_index[node]
	if !ok {
		err = ERROR_NOT_CONTAINS
	}
	return
}

//Swap(i, j) swaps nodelist[i] and nodelist[j]
func (c *Children) Swap(i, j int) {
	node1, node2 := c.Nodes[i], c.Nodes[j]
	c.node_index[node1], c.node_index[node2] = c.node_index[node2], c.node_index[node1]
	c.Nodes[i], c.Nodes[j] = c.Nodes[j], c.Nodes[i]

}

//SwapNodes(node1, node2) swaps the indice of node1 and node2
//If eather node1 or node2 is not in the NodeList, this method returns error
func (c *Children) SwapNodes(n1, n2 *Node) (err error) {
	i1, err1 := c.Index(n1)
	i2, err2 := c.Index(n2)
	if err1 != nil {
		err = err1
	} else if err2 != nil {
		err = err2
	} else {
		c.Swap(i1, i2)
	}
	return
}

//Tree Node have its parent's pointer ,NodeList of its children and its Value
type Node struct {
	parent   *Node
	children *Children
	Value    Value
}

//NewNode function creates a Node and returns a pointer of the Node
func NewNode(contents Value) *Node {
	node := new(Node)
	node.children = NewChildren()
	node.Value = contents
	return node
}

func (node *Node) Parent() *Node {
	return node.parent
}

func (node *Node) Children() []*Node {
	return node.children.Nodes
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
func (node *Node) AppendChild(child *Node) {
	node.children.Append(child)
	child.parent = node
}

//RemoveChild remove a child from the node's children
//If the child is not in the children, its returns error.
func (node *Node) RemoveChild(child *Node) (err error) {
	return node.children.Remove(child)
}

//Sibling returns a younger sibling Node and error.
//When youngest node calls Sibling(),
//Sibling returns nil and not-nil error.
func (node *Node) Sibling() (sibling *Node, err error) {
	parent := node.Parent()
	index, _ := parent.children.Index(node)
	if index+1 < len(parent.children.Nodes) {
		sibling = parent.children.Nodes[index+1]
	} else {
		err = ERROR_LAST_CHILD
	}
	return
}

/*Walk() is a generator.
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
		for _, child := range node.children.Nodes {
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
