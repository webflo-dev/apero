package ipc

import (
	"os"
)

type EmptyArgs struct{}

type AperoCtl int

func (a *AperoCtl) Quit(args *EmptyArgs, reply *int) error {
	os.Exit(0)
	return nil
}
