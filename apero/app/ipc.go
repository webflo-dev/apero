package app

import (
	"net"
	"net/rpc"
	"os"
	"webflo-dev/apero/logger"
)

const (
	SocketPath = "/tmp/apero.sock"
)

func startIPC() {
	aperoCtl := new(AperoCtl)
	rpc.Register(aperoCtl)

	os.RemoveAll(SocketPath)

	listener, err := net.Listen("unix", SocketPath)
	if err != nil {
		logger.IpcLogger.Fatalf("Unable to listen at %s. %w", SocketPath, err)
	}

	go rpc.Accept(listener)
}

// func wait() {
// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
// 	<-signals
// }

type EmptyArgs struct{}

type AperoCtl int

func (a *AperoCtl) Quit(args *EmptyArgs, reply *int) error {
	logger.IpcLogger.Println("Quit")
	os.Exit(0)
	return nil
}

func (a *AperoCtl) ApplyCSS(args *EmptyArgs, reply *int) error {
	logger.IpcLogger.Println("ApplyCSS")
	ApplyCSS("")
	return nil
}
