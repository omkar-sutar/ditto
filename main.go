package main

import (
	"fmt"
	"os"
	"strings"
)

var destinationPath string
var opMode string
var address string

func main() {
	// Check for minimum required arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	opMode = os.Args[1]

	if opMode == OperationModeServer {
		// Server mode: go run main.go server [port] [destination_path]
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go server [port] [destination_path]")
			os.Exit(1)
		}
		address = ":" + os.Args[2] // Format port as ":8080"
		destinationPath = os.Args[3]
		if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
			// Create destination path if it doesn't exist
			err := os.MkdirAll(destinationPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error creating destination path: %v\n", err)
				os.Exit(1)
			}
		}
		StartServer()
	} else if opMode == OperationModeClient {
		// Client mode: go run main.go client [address:port] [file_path_1] [file_path_2] ...
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go client [address:port] [file_path_1] [file_path_2] ...")
			os.Exit(1)
		}
		address = os.Args[2]
		filePaths := os.Args[3:]

		if len(filePaths) == 0 {
			fmt.Println("No files specified for sending")
			os.Exit(1)
		}

		if strings.Contains(filePaths[0], "*") {
			fmt.Printf("Warning: unmatched pattern '%s'\n", filePaths[0])
		}

		SendFiles(filePaths)
	} else {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  Server mode: go run main.go server [port] [destination_path]")
	fmt.Println("  Client mode: go run main.go client [address:port] [file_path_1] [file_path_2] ...")
}
