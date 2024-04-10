package ipc

import (
	"log"
	"net"
	"net/rpc"
	"os"
)

const (
	SocketPath = "/tmp/apero.sock"
)

func StartIPC() {
	aperoCtl := new(AperoCtl)
	rpc.Register(aperoCtl)

	os.RemoveAll(SocketPath)

	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		log.Fatalf("unable to listen at %s: %s", SocketPath, err)
	}
	go rpc.Accept(listener)
}

// func wait() {
// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
// 	<-signals
// }
