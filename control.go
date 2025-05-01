package main

const (
	ConnectionStart = iota + 1 // 1
	ConnectionEnd              // 2
	Ack                        // 3
	Nack                       // 4
)

const ContentLengthHeaderSize = 8

const (
	OperationModeServer = "server"
	OperationModeClient = "client"
)
