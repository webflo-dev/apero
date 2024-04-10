package ipc

import (
	"log"
	"net/rpc"
)

func NewClient() *rpc.Client {
	client, err := rpc.Dial("unix", SocketPath)
	if err != nil {
		log.Fatalf("failed: %s", err)
	}
	return client;
}

func SendQuit() (int,error) {
	client := NewClient()

	defer client.Close()

	args := &EmptyArgs{}
	var reply int
	err := client.Call("AperoCtl.Quit", args, &reply)
	return reply, err
}