package main

import (
	"fmt"
)

func main() {

	fmt.Printf("Ya\n")

	var ch1 chan bool
	ch1 = make(chan bool)
	var ch2 chan bool
	close(ch1)
	d, ok := <-ch1
	fmt.Printf("%v, %v\n", d, ok)
	defer func() {
		//if ch1 != nil {
			//fmt.Printf("close1\n")
			close(ch1)
			close(ch1)
		//}
	}()
	defer func() {
		if ch2 != nil {
			fmt.Printf("close2\n")
			close(ch2)
		}
	}()
}
