package hyprland

import (
	"fmt"
	"strconv"
)

func toInt(value string) int {
	ret, _ := strconv.Atoi(value)
	return ret
}

func toBool(value string) bool {
	return value == "1"
}

func toAddress(value string) string {
	return fmt.Sprintf("0x%s", value)
}
