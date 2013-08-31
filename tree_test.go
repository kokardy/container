package tree

import (
	"log"
	"testing"
)

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

}
