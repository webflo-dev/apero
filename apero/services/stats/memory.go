package stats

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type MemoryStats struct {
	Total               uint64
	Used                uint64
	Cached              uint64
	Free                uint64
	Available           uint64
	Buffers             uint64
	MemAvailableEnabled bool
}

func GetMemoryStats() (*MemoryStats, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var memoryStats = MemoryStats{}
	statsMap := map[string]*uint64{
		"MemTotal":     &memoryStats.Total,
		"MemFree":      &memoryStats.Free,
		"MemAvailable": &memoryStats.Available,
		"Cached":       &memoryStats.Cached,
		"Buffers":      &memoryStats.Buffers,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		field := line[:i]

		if ptr := statsMap[field]; ptr != nil {
			val := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				*ptr = v / 1024
			}
			if field == "MemAvailable" {
				memoryStats.MemAvailableEnabled = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if memoryStats.MemAvailableEnabled {
		memoryStats.Used = memoryStats.Total - memoryStats.Available
	} else {
		memoryStats.Used = memoryStats.Total - memoryStats.Free - memoryStats.Buffers - memoryStats.Cached
	}
	return &memoryStats, nil
}
