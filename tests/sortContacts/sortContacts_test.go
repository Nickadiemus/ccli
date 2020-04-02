package main

import (
	"os"
	"testing"
)

func TestSortContactsAscendingOrder(*testing.T) {
	workDir := os.Getenv("PWD")
	testPath := "/tests/sortContacts/input/"
	// struct for test cases
	tables := []struct {
		i []Person
		o []Person
	}{
		loadFile(workDir + testPath + "input1.json"),
		loadFile(workDir + testPath + "input2.json"),
		loadFile(workDir + testPath + "input3.json"),
	}

}
