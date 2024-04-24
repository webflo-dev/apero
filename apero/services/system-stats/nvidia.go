package systemStats

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type NvidiaStats struct {
	Name          string
	DriverVersion string
	GpuUsage      int
	MemoryUsage   int
	MemoryTotal   int
	MemoryUsed    int
	MemoryFree    int
	Temperature   int
	FanSpeed      int
}

type Field struct {
	name string
	ptr  interface{}
}

func GetNvidiaStats() (*NvidiaStats, error) {
	nvidiaStats := NvidiaStats{}

	fields := []Field{
		{"name", &nvidiaStats.Name},
		{"driver_version", &nvidiaStats.DriverVersion},
		{"utilization.gpu", &nvidiaStats.GpuUsage},
		{"utilization.memory", &nvidiaStats.MemoryUsage},
		{"memory.total", &nvidiaStats.MemoryTotal},
		{"memory.used", &nvidiaStats.MemoryUsed},
		{"memory.free", &nvidiaStats.MemoryFree},
		{"temperature.gpu", &nvidiaStats.Temperature},
		{"fan.speed", &nvidiaStats.FanSpeed},
	}

	keys := make([]string, 0, len(fields))
	for _, field := range fields {
		keys = append(keys, field.name)
	}
	cmdFields := strings.Join(keys, ",")

	out, err := exec.Command(
		"nvidia-smi",
		"--query-gpu="+cmdFields,
		"--format=csv,noheader,nounits").Output()

	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	record, err := csvReader.Read()

	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}

	for index, field := range fields {
		value := record[index]
		switch field.ptr.(type) {
		case *int:
			if num, err := strconv.Atoi(value); err == nil {
				*field.ptr.(*int) = num
			}
		default:
			*field.ptr.(*string) = value
		}
	}
	return &nvidiaStats, nil
}
