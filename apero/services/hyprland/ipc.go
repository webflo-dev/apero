package hyprland

import (
	"fmt"
	"net"
	"os"
	"strings"
	"webflo-dev/apero/logger"
)

var hyprlandInstanceSignature = os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

func createEventsConnection() net.Conn {
	socketPath := fmt.Sprintf("/tmp/hypr/%s/.socket2.sock", hyprlandInstanceSignature)
	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		logger.AppLogger.Fatalln("Cannot connect to Hyprland service (.socket2.sock)")
	}

	return connection
}

func (s *hyprlandEventService) closeConnection() {
	if err := s.connection.Close(); err != nil {
		logger.AppLogger.Println("Could not close connection", err)
	}
}

func (s *hyprlandEventService) receive() ([]EventData, error) {
	buf := make([]byte, bufferSize)
	n, err := s.connection.Read(buf)
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

func createWritableConnection() net.Conn {
	socketPath := fmt.Sprintf("/tmp/hypr/%s/.socket.sock", hyprlandInstanceSignature)
	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		logger.AppLogger.Fatalln("Cannot connect to Hyprland service (.socket2.sock)")
	}

	return connection
}
