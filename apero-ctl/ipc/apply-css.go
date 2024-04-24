package ipc

func (client *IpcClient) ApplyCSS() (int, error) {
	defer client.Close()

	var reply int
	err := client.Call("AperoCtl.ApplyCSS", &emptyArgs{}, &reply)
	return reply, err
}
