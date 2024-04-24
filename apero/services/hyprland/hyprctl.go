package hyprland

import (
	"encoding/json"
)

func Dispatch(command string) error {
	_, error := writeCmd(command)
	return error
}

func Workspaces() ([]Workspace, error) {
	var workspaces []Workspace
	return workspaces, request("j/workspaces", &workspaces)
}

func writeCmd(command string) ([]byte, error) {
	connection := createWritableConnection()

	message := []byte(command)
	_, err := connection.Write(message)
	if err != nil {
		return nil, err
	}

	reply := make([]byte, 102400)
	n, err := connection.Read(reply)
	if err != nil {
		return nil, err
	}

	defer connection.Close()

	return reply[:n], nil
}

func request(command string, target any) error {
	buf, err := writeCmd(command)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(buf, target); err != nil {
		return err
	}

	return nil
}
