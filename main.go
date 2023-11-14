package main

import (
	"bufio"
	"fmt"
	"os"
	"math/bits"
)

type Tree struct {
	id uint64
	parent *Tree
	children []*Tree
	receiver chan uint64
}

func (t *Tree) addChild(new *Tree) {
	t.children = append(t.children, new)
	new.parent = t
}

func (t *Tree) sixColoringRound() {
	t.sendIdToChildren()
	if t.parent == nil {
		t.id = 0 | (t.id & 1)
		return 
	}
	
	oldParentsId := <- t.receiver
	index := t.getSmallestDiffIndex(oldParentsId)

	// Finds value at LSB
	valueAtIndex := (t.id >> index) & 1
	t.id = ((index << 1) | valueAtIndex)
}

func (t * Tree) sendIdToChildren() {
	for _, child := range t.children {
		child.receiver <- t.id
	}
}

func (t *Tree) getSmallestDiffIndex(oldParentsId uint64) uint64 {
	// Sets 1 on bits that are different between oldParentsid and t.id
	xorIds := oldParentsId ^ t.id
	// Returns LSB
	return uint64(bits.TrailingZeros64(xorIds))
}

func main() {

	filePath := "example_tree.txt"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// line := scanner.Text()
		// fields := strings.Fields(line)
		// fmt.Println(fields)
	}
}
