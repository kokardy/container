package container

import (
	"log"
	"testing"
)

func String(node *Node) string {
	s, _ := node.Value.(string)
	return s
}

func TestNode(t *testing.T) {
	log.Println("test")
	root := NewNode("root")
	log.Println("root:", root)
	a := NewNode("a")
	root.AppendChild(a)
	b := NewNode("b")
	c := NewNode("c")
	d := NewNode("d")
	e := NewNode("e")
	f := NewNode("f")
	g := NewNode("g")

	a.AppendChild(b)
	a.AppendChild(c)
	a.AppendChild(d)
	a.AppendChild(e)

	c.AppendChild(f)
	c.AppendChild(g)

	log.Println("root children:", root.Children())

	log.Println("root parent:", root.Parent())

	log.Println("path e:", e.Path())

	for child := range root.Walk() {
		log.Println("walk:", child)
	}

	log.Println("node_index:", a.children.node_index)
	log.Println("nodes:", a.children.Nodes)
	log.Println("parent:", b.Parent())
	var n *Node
	n, _ = b.Sibling()
	log.Println("b sibling:", String(n))
	n, _ = c.Sibling()
	log.Println("c sibling:", String(n))
	n, _ = d.Sibling()
	log.Println("d sibling:", String(n))

	for n := b; ; {
		log.Println("sibling:", String(n))
		next, err := n.Sibling()
		if err != nil {
			break
		}
		n = next

	}

	if a.Children()[0] != b {
		t.Fatal("a's first child is not b")
	}

	if c.Children()[1] != g {
		t.Fatal("c's second child is not g")
	}

	a.AppendChild(g)

	if g.Parent() != a {
		t.Fatal("g's parent is not a")
	}

	if a.Children()[4] != g {
		t.Fatal("a's fourth children is not g")
	}

}
