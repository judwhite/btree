package btree

import (
	"bytes"
	"fmt"
	"strings"
)

const factor = 4
const half = factor / 2

type Tree struct {
	Root *Node
}

type Node struct {
	leaf     bool
	elems    []uint64
	children []*Node
	parent   *Node
}

func New() *Tree {
	return &Tree{
		Root: &Node{
			leaf:  true,
			elems: make([]uint64, 0, factor),
		},
	}
}

func (t *Tree) String() string {
	return t.Root.String()
}

func (t *Tree) Search(key uint64) (uint64, error) {
	n, idx := t.Root.search(key)
	if idx == -1 {
		return 0, fmt.Errorf("key %d not found", key)
	}
	return n.elems[idx], nil
}

func (n *Node) search(key uint64) (*Node, int) {
	//fmt.Printf("search: key %d, leaf: %v\n", key, n.leaf)
	i := 0
	for ; i < len(n.elems); i++ {
		if !n.leaf {
			if key < n.elems[i] {
				return n.children[i].search(key)
			}
		} else {
			if key == n.elems[i] {
				return n, i
			}
		}
	}
	if !n.leaf {
		return n.children[i].search(key)
	}
	return n, -1
}

func (t *Tree) Insert(key uint64) {
	n, idx := t.Root.search(key)
	if idx != -1 {
		return
	}
	n.insertLeaf(key, t)
}

func (n *Node) String() string {
	return n.string(0)
}

func (n *Node) string(depth int) string {
	var buf bytes.Buffer
	buf.WriteString(strings.Repeat("-", depth+1))
	buf.WriteString(" ")
	buf.WriteString(n.elemsString())
	buf.WriteString("\n")
	for _, child := range n.children {
		buf.WriteString(child.string(depth + 1))
	}
	return buf.String()
}

func (n *Node) elemsString() string {
	var buf bytes.Buffer
	for _, elem := range n.elems {
		if n.leaf {
			buf.WriteString(fmt.Sprintf("[%03d] ", elem))
		} else {
			buf.WriteString(fmt.Sprintf("%03d ", elem))
		}
	}
	return buf.String()
}

func (n *Node) split(t *Tree) {
	//fmt.Printf("Node.split: %v\n", n.elems)
	//fmt.Printf("tree before: %s\n", t)
	/*defer func() {
		fmt.Printf("tree after: %s\n", t)
	}()*/
	if n.parent == nil {
		t.split()
		return
	}

	left, right := n.createTwo()

	parent := n.parent
	parent.replaceNode(n, left)
	parent.insertNode(right, t)
}

func (n *Node) replaceNode(remove *Node, insert *Node) {
	//fmt.Printf("Node.replaceNode: in: %v replace: %v with: %v\n", n.elems, remove.elems, insert.elems)
	for i := 0; i < len(n.children); i++ {
		if n.children[i] == remove {
			n.children[i] = insert
			if i > 0 && n.elems[i-1] != insert.least() {
				panic(fmt.Sprintf("node: %s, want: %d got: %d", n, insert.least(), n.elems[i-1]))
			}
			return
		}
	}
	panic(fmt.Errorf("node not found: parent: %v, remove: %v\n", n.String(), remove.String()))
}

func (n *Node) insertNode(insert *Node, t *Tree) {
	//fmt.Printf("Node.insertNode: into: %v node: %v\n", n.elems, insert.elems)

	newKey := insert.least()

	found := false
	for i := 0; i < len(n.children); i++ {
		least := n.children[i].least()
		if least > newKey {
			n.children = append(n.children, nil)
			copy(n.children[i+1:], n.children[i:])
			n.children[i] = insert

			n.elems = append(n.elems, 0)
			n.recalculateElems()

			found = true
			break
		}
	}

	if !found {
		n.elems = append(n.elems, newKey)
		n.children = append(n.children, insert)
	}

	if len(n.elems) > factor {
		n.split(t)
	}
}

func (n *Node) recalculateElems() {
	for i := 1; i < len(n.children); i++ {
		n.elems[i-1] = n.children[i].least()
	}
}

func (t *Tree) split() {
	left, right := t.Root.createTwo()
	//fmt.Printf("Tree.split: %v -> ", t.Root.elems)

	root := &Node{
		leaf:     false,
		elems:    make([]uint64, 1, factor),
		children: make([]*Node, 2, factor),
	}

	root.children[0] = left
	root.children[1] = right
	root.elems[0] = right.least()

	left.parent, right.parent = root, root

	t.Root = root

	//fmt.Printf("%v -> %v, %v\n", t.Root.elems, left.elems, right.elems)
}

func (n *Node) least() uint64 {
	if n.leaf {
		return n.elems[0]
	}
	return n.children[0].least()
}

func (n *Node) createTwo() (*Node, *Node) {
	left := &Node{leaf: n.leaf, elems: make([]uint64, 0, factor), parent: n.parent}
	right := &Node{leaf: n.leaf, elems: make([]uint64, 0, factor), parent: n.parent}

	if n.leaf {
		for i := 0; i < len(n.elems); i++ {
			if i < half {
				left.elems = append(left.elems, n.elems[i])
			} else {
				right.elems = append(right.elems, n.elems[i])
			}
		}
	} else {
		left.children = make([]*Node, 0, factor)
		right.children = make([]*Node, 0, factor)

		left.children = append(left.children, n.children[0])
		for i := 0; i < len(n.elems); i++ {
			if i < half {
				left.children = append(left.children, n.children[i+1])
			} else {
				right.children = append(right.children, n.children[i+1])
			}
		}

		// set elems, parent
		for i, child := range left.children {
			if i > 0 {
				left.elems = append(left.elems, child.least())
			}
			child.parent = left
		}
		for i, child := range right.children {
			if i > 0 {
				right.elems = append(right.elems, child.least())
			}
			child.parent = right
		}
	}

	return left, right
}

func (n *Node) insertLeaf(key uint64, t *Tree) {
	//fmt.Printf("Node.insertLeaf: key: %d node: %v leaf: %v\n", key, n.elems, n.leaf)

	i := 0
	for ; i < len(n.elems); i++ {
		if n.elems[i] > key {
			break
		}
	}

	if n.leaf {
		if i == len(n.elems) {
			n.elems = append(n.elems, key)
		} else {
			n.elems = append(n.elems, 0)
			copy(n.elems[i+1:], n.elems[i:])
			n.elems[i] = key
		}
		if len(n.elems) > factor {
			n.split(t)
		}
		return
	}

	n.children[i].insertLeaf(key, t)
}
