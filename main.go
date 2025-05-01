package main

import (
	"fmt"
	"os"
	"strings"
)

var destinationPath string
var opMode string

func main() {

	// Get arg0 as destination path
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <destination_path>")
		os.Exit(1)
	}

	opMode = os.Args[1]
	if opMode == OperationModeServer {
		if len(os.Args) < 3 {
			fmt.Println("Us.3age: go run main.go server <destination_path>")
			os.Exit(1)
		}
		destinationPath = os.Args[2]
		StartServer()
	}
	// Client mode
	filePaths := os.Args[2:]
	if len(filePaths) == 0 {
		fmt.Println("Usage: go run main.go client <file_path_1> <file_path_2> ...")
		os.Exit(1)
	}
	if strings.Contains(filePaths[0], "*") {
		fmt.Printf("Warning: unmatched pattern '%s'\n", filePaths[0])
	}
	SendFiles(filePaths)
}
