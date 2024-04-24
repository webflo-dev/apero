package ipc

import (
	"log"
	"net/rpc"
)

var socketPath = "/tmp/apero.sock"

type emptyArgs struct{}

type IpcClient struct {
	*rpc.Client
}

func NewClient() *IpcClient {
	client, err := rpc.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}

	// defer client.Close()

	return &IpcClient{client}
}
