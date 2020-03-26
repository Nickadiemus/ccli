package main

import (
	"flag"
	"fmt"
)

func main() {

	flaginptc := flag.String("c", "create", "-c creates a new contact")
	flaginptf := flag.String("f", "first_name", "-f creates field for the first name of a new contact")
	flaginptl := flag.String("l", "last_name", "-l creates field for the last name of a new contact")

	flag.Parse()

	fmt.Println("-c: ", *flaginptc)
	fmt.Println("-f: ", *flaginptf)
	fmt.Println("-l: ", *flaginptl)
	fmt.Println("tail: ", flag.Args())

}
