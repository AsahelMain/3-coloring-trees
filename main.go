package main

import (
	"bufio"
	"fmt"
	"os"
	"math/bits"
	"strconv"
	"strings"
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

func buildTreeFromFile(filePath string) (*Tree, int, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, 0, err
	}

	scanner := bufio.NewScanner(file)
	treeMap := make(map[uint64]*Tree)

	scanner.Scan()
	rootID, err := strconv.ParseUint(scanner.Text(), 10, 64)

	if err != nil {
		return nil, 0, err
	}
	root := &Tree{id: rootID}
	treeMap[rootID] = root

	for scanner.Scan(){
		line := scanner.Text()
		fields := strings.Fields(line)
		parentID, err := strconv.ParseUint(fields[0], 10, 64)
		
		if err != nil {
			return nil, 0, err
		}

		parentNode, exists := treeMap[parentID]

		if !exists {
			parentNode = &Tree{id: parentID}
			treeMap[parentID] = parentNode
		}

		for i := 1; i < len(fields); i++ {
			childID, err := strconv.ParseUint(fields[i], 10, 64)

			if err != nil {
				return nil, 0, err
			}

			childNode, exists := treeMap[childID]

			if !exists {
				childNode = &Tree{id: childID}
				treeMap[childID] = childNode
			}
			parentNode.addChild(childNode)
		}
	}

	return root, len(treeMap), nil

}

func traverseTree(root *Tree) {
	if root == nil {
		return
	}

	fmt.Println(root.id)

	for _, child := range root.children {
		traverseTree(child)
	}
}



func main() {
	filePath := "example_tree.txt"

	tree, size, err := buildTreeFromFile(filePath)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	fmt.Println("Done creating the tree")
	fmt.Printf("Root ID: %d\n", tree.id)
	fmt.Printf("Size of the tree: %d\n", size)

	traverseTree(tree)

	/*var wg sync.WaitGroup
	wg.Add(size)

	for {
		traverse


	}*/
	
	
}
