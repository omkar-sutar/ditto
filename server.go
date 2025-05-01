package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func StartServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("Server listening on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	if n == 0 || buffer[0] != ConnectionStart {
		fmt.Println("Invalid start byte, closing connection")
		return
	}

	for {

		//1. Read file name
		fileNameBytes, err := ReadData(conn)
		if err != nil {
			fmt.Println("Error reading file name:", err)
			return
		}

		fileName := string(fileNameBytes)
		fileName = filepath.Base(fileName)

		if fileName == "" {
			fmt.Println("Empty file name, closing connection")
			return
		}
		fmt.Println("Received file name:", fileName)

		//2. Read file size
		fileSizeBytes, err := ReadData(conn)
		if err != nil {
			fmt.Println("Error reading file size:", err)
			return
		}
		fileSize := binary.BigEndian.Uint64(fileSizeBytes)
		if fileSize == 0 {
			fmt.Println("Invalid file size, closing connection")
			return
		}
		fmt.Printf("Received file size: %d bytes, %f MB\n", fileSize, float64(fileSize)/1024/1024)

		if fileExists(filepath.Join(destinationPath, fileName)) {
			fmt.Printf("File %s already exists, skipping download\n", fileName)
			// Send ACK for file existence
			_, err = conn.Write([]byte{Ack})
			if err != nil {
				fmt.Println("Error sending ACK:", err)
				return
			}
			continue
		}
		conn.Write([]byte{Nack}) // Send Nack for file non-existence

		// Read file data
		file, err := os.Create(filepath.Join(destinationPath, fileName))
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		_, err = io.CopyN(file, conn, int64(fileSize))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			file.Close()
			os.Remove(filepath.Join(destinationPath, fileName)) // cleanup
			return
		}
		file.Close()
		fmt.Printf("File %s downloaded successfully\n", fileName)

		stopByte := make([]byte, 1)
		n, err = conn.Read(stopByte)
		if err != nil || n != 1 {
			fmt.Println("Read error:", err)
			return
		}
		if stopByte[0] == ConnectionEnd {
			fmt.Println("Connection end byte received, closing connection")
			return
		}
	}
}

func ReadData(conn net.Conn) ([]byte, error) {
	sizeInfoBuffer := make([]byte, ContentLengthHeaderSize)
	n, err := io.ReadFull(conn, sizeInfoBuffer)
	if err != nil {
		return nil, err
	}
	if n != ContentLengthHeaderSize {
		return nil, fmt.Errorf("expected %d bytes, got %d", ContentLengthHeaderSize, n)
	}
	contentLength := binary.BigEndian.Uint64(sizeInfoBuffer)
	dataBuffer := make([]byte, contentLength)
	_, err = io.ReadFull(conn, dataBuffer)
	return dataBuffer, err
}

// Function to check if given filepath already exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
