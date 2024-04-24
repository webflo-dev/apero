package ipc

func (client *IpcClient) SendQuit() (int, error) {
	defer client.Close()

	var reply int
	err := client.Call("AperoCtl.Quit", &emptyArgs{}, &reply)
	return reply, err
}
