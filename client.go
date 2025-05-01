package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func SendFiles(filePaths []string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Process each file
	for i, filePath := range filePaths {
		fmt.Printf("[%d/%d]\n", i+1, len(filePaths))
		err := sendFile(conn, filePath)
		if err != nil {
			fmt.Printf("Error sending file %s: %v\n", filePath, err)
			continue
		}
	}

	// Send connection end byte
	_, err = conn.Write([]byte{ConnectionEnd})
	if err != nil {
		fmt.Println("Error sending connection end byte:", err)
		return
	}

	fmt.Println("All files sent successfully")
}

func sendFile(conn net.Conn, filePath string) error {
	fmt.Printf("Sending file: %s\n", filePath)

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("unable to stat file: %v", err)
	}

	// Open file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	// Send connection start byte for file
	_, err = conn.Write([]byte{ConnectionStart})
	if err != nil {
		return fmt.Errorf("error sending connection start byte for next file: %v", err)
	}

	// Send file name
	err = sendData(conn, []byte(filePath))
	if err != nil {
		return fmt.Errorf("error sending file name: %v", err)
	}

	// Send file size
	sizeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sizeBytes, uint64(fileInfo.Size()))
	err = sendData(conn, sizeBytes)
	if err != nil {
		return fmt.Errorf("error sending file size: %v", err)
	}

	// Check if file already exists on server
	ackBuffer := make([]byte, 1)
	_, err = io.ReadFull(conn, ackBuffer)
	if err != nil {
		return fmt.Errorf("error reading acknowledgment: %v", err)
	}

	if ackBuffer[0] == Ack {
		fmt.Printf("File %s already exists on server, skipping\n", filepath.Base(filePath))
		return nil
	}

	// Send file data
	fmt.Printf("Sending file data (%d bytes, %f MB)...\n", fileInfo.Size(), float64(fileInfo.Size())/1024/1024)

	n, err := io.Copy(conn, file)
	if err != nil {
		return fmt.Errorf("error sending file data: %v", err)
	}

	if n != fileInfo.Size() {
		return fmt.Errorf("sent %d bytes but file size is %d bytes", n, fileInfo.Size())
	}

	fmt.Printf("File %s sent successfully\n", filepath.Base(filePath))

	return nil
}

func sendData(conn net.Conn, data []byte) error {
	// Send content length header
	contentLengthBytes := make([]byte, ContentLengthHeaderSize)
	binary.BigEndian.PutUint64(contentLengthBytes, uint64(len(data)))

	_, err := conn.Write(contentLengthBytes)
	if err != nil {
		return err
	}

	// Send actual data
	_, err = conn.Write(data)
	return err
}
