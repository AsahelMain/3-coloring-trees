package main

import (
	"bufio"
	"fmt"
	"os"
	"math/bits"
	"strconv"
	"strings"
	"sync"
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
	new.receiver = make(chan uint64)
}

func (t *Tree) sixColoringRound(wg *sync.WaitGroup) {
	defer wg.Done()
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

	//fmt.Println(t.id)
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
	root.receiver = make(chan uint64)
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

func treeToList(root *Tree, nodes *[]*Tree) {
	if root == nil {
		return
	}

	*nodes = append(*nodes, root)

	for _, child := range root.children {
		treeToList(child, nodes)
	}
}

func checkColorsRange(nodes []*Tree) bool{
	for _, node := range nodes{
		if node.id > 5{
			return false
		}
	}

	return true
}


func main() {
	filePath := "example_tree3.txt"

	tree, size, err := buildTreeFromFile(filePath)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	fmt.Println("Done creating the tree")
	fmt.Printf("Root ID: %d\n", tree.id)
	fmt.Printf("Size of the tree: %d\n", size)

	nodes := []*Tree{}
	treeToList(tree, &nodes)

	var wg sync.WaitGroup
	
	for{
		wg.Add(size)
		for _, node := range nodes{
			go node.sixColoringRound(&wg)
		}

		wg.Wait()

		for i, node := range nodes{
			fmt.Printf("%d", node.id)
			if i != len(nodes) - 1{
				fmt.Print(", ")
			}
		}

		fmt.Printf("\n")

		colorsInRange := checkColorsRange(nodes)

		if colorsInRange {
			fmt.Println("Six coloring stage completed!")
			break;
		}
		
	}
	
	
	
}
