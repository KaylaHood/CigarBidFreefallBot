package main

// This has no relation to the cigarbid project,
// I just included it because I was happy that I figured it out
// (from A Tour of Go: https://tour.golang.org/concurrency/8)

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func MinWalk(t *tree.Tree, ch chan int, orig bool) {
	//fmt.Println(t.String())
	if t.Left != nil {
		left := *(t.Left)
		//fmt.Println("Going left")
		MinWalk(&left, ch, false)
	}
	//fmt.Printf("Adding %d to channel\n", t.Value)
	ch <- t.Value
	if t.Right != nil {
		right := *(t.Right)
		//fmt.Println("Going right")
		MinWalk(&right, ch, false)
	}
	if orig {
		close(ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	var t1channel chan int = make(chan int)
	var t2channel chan int = make(chan int)
	go MinWalk(t1, t1channel, true)
	go MinWalk(t2, t2channel, true)
	for i := range t1channel {
		j := <-t2channel
		fmt.Printf("i: %d, j: %d\n", i, j)
		if i != j {
			return false
		}
	}
	return true
}

func main() {
	var myTree *tree.Tree = tree.New(1)
	var channel chan int = make(chan int)
	//fmt.Printf(myTree.Left.String())
	go MinWalk(myTree, channel, true)
	for x := 0; x < 10; x++ {
		v := <-channel
		fmt.Printf("Found %d in channel\n", v)
	}

	var myTree2 *tree.Tree = tree.New(2)
	var myTree3 *tree.Tree = tree.New(1)
	result := Same(myTree3, myTree)
	fmt.Printf("Should be true: %v\n", result)
	result2 := Same(myTree2, myTree3)
	fmt.Printf("Should be false: %v\n", result2)
}
