package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
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

func (t *Tree) shiftDown(wg *sync.WaitGroup) {
	defer wg.Done()
	t.sendIdToChildren()
	if t.parent == nil {
		t.id = (t.id + 1) % 3
		return
	}
	t.id = <-t.receiver
}

func (t *Tree) first_free(wg *sync.WaitGroup, round int) {
	defer wg.Done()
	var innerWg sync.WaitGroup
	innerWg.Add(2)
	go func(wg *sync.WaitGroup){
		defer wg.Done()
		t.sendIdToChildren()
		if (t.parent != nil) {
			t.parent.receiver <- t.id
		}
	}(&innerWg)
	go func(wg *sync.WaitGroup){
		defer wg.Done()
		to_receive := len(t.children)
		if (t.parent != nil) {
			to_receive++
		}
		
		ids_received := make(map[uint64]bool)
		for i := 0; i < to_receive; i++ {
			ids_received[<-t.receiver] = true
		}

		if (t.id != uint64(round)) {
			return
		}
		for i := uint64(0); i <= 2; i++ {
			if ids_received[i] == false {
				t.id = i
				break
			}
		}

	}(&innerWg)
	innerWg.Wait()
}
func six2threeRound(nodes[] *Tree, round int) {
	var wg sync.WaitGroup	
	wg.Add(len(nodes))
	for _, node := range nodes {
		go node.shiftDown(&wg)
	}
	wg.Wait()
	wg.Add(len(nodes))
	for _, node := range nodes {
		go node.first_free(&wg, round)
	}
	wg.Wait()
}
func  six2three(nodes []*Tree) {
	for round := 3; round <= 5; round++ {
		six2threeRound(nodes, round)
	}
}


func (t *Tree) sendIdToNeighbors() {
	t.sendIdToChildren()
	if (t.parent == nil) {
		return
	}
	t.parent.receiver <- t.id
}

func (t *Tree) sixColoringRound(wg *sync.WaitGroup) {
	defer wg.Done()
	t.sendIdToChildren()
	// If t is the root.
	if t.parent == nil {
		// Sets index equal to 0 and concatenates with
		// the value of t.id[0].
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

func (t* Tree) String() string {
	var s string
	s = fmt.Sprintf("%d: ", t.id)
	for i, child := range t.children {
		s = fmt.Sprintf("%s%d", s, child.id)
		if i != len(t.children) - 1{
			s += ", "
		}
	}
	return s
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
	// fmt.Println(root)
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

func processNode(t *Tree) []uint64{
	var list []uint64
	list = append(list,  t.id)
	for _, child:= range t.children{
		list = append(list, child.id)
	}
	return list
}

func treeToListFile(t *Tree) [][]uint64{
	var queue = make([]*Tree,0)
	queue = append(queue, t)
	var list  [][]uint64
	list = append(list,[]uint64{t.id})
	for len(queue) != 0 {
		node := queue[0]
		queue = queue[1:]
		if node.children != nil {
			list = append (list,processNode(node))
			for  _, child := range node.children{
				queue = append(queue, child)
			}
		}
	}
	return list
}

func WriteByLines(filename string, lines [][]uint64){
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	for _, line := range lines{
		for _, val := range line{
			_, err := fmt.Fprintf(file, "%s " , strconv.FormatUint(val, 10))
			if err != nil{
				fmt.Println(err.Error())
				return
			}
		}
		_, err:= fmt.Fprintf(file,"\n")
		if err != nil{
			fmt.Println(err.Error())
			return
		}
	}
}

func printTree(root *Tree, level int) {
	if root == nil {
		return
	}

	color := "\033[0m" 

	switch root.id {
	case 0:
		color = "\033[34m" 
	case 1:
		color = "\033[31m" 
	case 2:
		color = "\033[32m" 
	}

	fmt.Printf("%s%s- Node ID: %d\n", color, strings.Repeat("  ", level), root.id)

	
	for _, child := range root.children {
		printTree(child, level+1)
	}

	fmt.Print("\033[0m")
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
	fmt.Println("Original Tree:")
	for _, node := range nodes{
		fmt.Println(node)
	}
	fmt.Println()
	var wg sync.WaitGroup
	for{
		wg.Add(size)
		for _, node := range nodes{
			go node.sixColoringRound(&wg)
		}

		wg.Wait()
		colorsInRange := checkColorsRange(nodes)

		if colorsInRange {
			fmt.Println("Six coloring stage completed!")
			break;
		}
	
	}
	for _, node := range nodes{
		fmt.Println(node)
	}
	fmt.Println()
	six2three(nodes)
	fmt.Println("six2three stage completed!")
	for _, node := range nodes {
		fmt.Println(node)
	}

	WriteByLines("output.txt",treeToListFile(tree))

	fmt.Println("\nResulting tree:")
	printTree(tree, 0)
}
