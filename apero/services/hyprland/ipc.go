package hyprland

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	bufferSize     = 8192
	eventSeperator = ">>"
)

var hyprlandInstanceSignature = os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
var eventSocketPath = fmt.Sprintf("/tmp/hypr/%s/.socket2.sock", hyprlandInstanceSignature)
var writableSocketPath = fmt.Sprintf("/tmp/hypr/%s/.socket.sock", hyprlandInstanceSignature)

func createEventsConnection() net.Conn {
	connection, err := net.Dial("unix", eventSocketPath)
	if err != nil {
		logger.Fatalln("Cannot connect to Hyprland service (.socket2.sock)")
	}
	return connection
}

func closeConnection(connection net.Conn) {
	if err := connection.Close(); err != nil {
		logger.Println("Could not close connection", err)
	}
}

func readEvent(connection net.Conn) ([]EventData, error) {

	buf := make([]byte, bufferSize)
	n, err := connection.Read(buf)
	if err != nil {
		return nil, err
	}

	buf = buf[:n]
	rawEvents := strings.Split(string(buf), "\n")
	var eventData []EventData
	for _, event := range rawEvents {
		if event == "" {
			continue
		}

		split := strings.Split(event, eventSeperator)
		if split[0] == "" || split[1] == "" || split[1] == "," {
			continue
		}

		eventData = append(eventData, EventData{
			Type: EventType(split[0]),
			Data: split[1],
		})
	}

	return eventData, nil
}

func writeCmd(command string, target any) error {
	connection, err := net.Dial("unix", writableSocketPath)
	if err != nil {
		logger.Fatalln("Cannot connect to Hyprland service (.socket.sock)")
	}

	message := []byte(command)
	_, err = connection.Write(message)
	if err != nil {
		return err
	}

	reply := make([]byte, 102400)
	n, err := connection.Read(reply)
	if err != nil {
		return err
	}

	defer connection.Close()

	buf := reply[:n]

	if target != nil {
		if err := json.Unmarshal(buf, target); err != nil {
			return err
		}
	}

	return nil
}
