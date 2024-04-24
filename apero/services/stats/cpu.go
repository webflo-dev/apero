package stats

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CpuStats struct {
	Usage int
}

var cpuLastSum = 0
var cpuLastIdle = 0

func GetCpuStats() (*CpuStats, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to scan /proc/stat")
	}

	line := scanner.Text()

	values := strings.Fields(line)[1:]

	sum := sumStrings(values)
	delta := sum - cpuLastSum
	idle, _ := strconv.Atoi(values[3])
	idleDelta := idle - cpuLastIdle
	used := delta - idleDelta
	usage := 100 * used / delta

	cpuLastSum = sum
	cpuLastIdle = idle

	cpuStats := &CpuStats{
		Usage: usage,
	}

	return cpuStats, nil
}

func sumStrings(source []string) int {
	n := 0
	for _, v := range source {
		i, _ := strconv.Atoi(v)
		n += i
	}
	return n
}
