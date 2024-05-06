package bar

import (
	"demo/utils"
	"fmt"
	"math"
	systemStats "webflo-dev/apero/services/system-stats"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/gtk"
)

type systemInfo struct {
	box   *gtk.Box
	label *gtk.Label
	icon  *gtk.Image
}

type statsHandler struct {
	systemStats.Subscriber
	cpu    *systemInfo
	memory *systemInfo
	nvidia *systemInfo
}

func newSystemInfoModule() *gtk.Box {

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 24)
	box.SetName("system-info")

	cpu := newSystemInfoBox(Icon_SystemStats_Cpu)
	box.Add(cpu.box)

	memory := newSystemInfoBox(Icon_SystemStats_Memory)
	box.Add(memory.box)

	nvidia := newSystemInfoBox(Icon_SystemStats_Gpu)
	box.Add(nvidia.box)

	systemStats.Register(&statsHandler{
		cpu:    cpu,
		memory: memory,
		nvidia: nvidia,
	}, systemStats.EventAll)

	return box
}

func (h *statsHandler) UpdateAll(stats *systemStats.SystemStats) {
	h.cpu.SetValue(stats.Cpu.Usage)
	h.memory.SetValue(int(math.Floor((float64(stats.Memory.Used) / float64(stats.Memory.Total)) * 100)))
	h.nvidia.SetValue(stats.Nvidia.GpuUsage)
}

func newSystemInfoBox(iconName string) *systemInfo {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	ui.AddCSSClass(&box.Widget, "info")

	icon := ui.NewFontSizeImageFromIconName(iconName)
	box.Add(icon)

	label, _ := gtk.LabelNew("")
	box.Add(label)

	return &systemInfo{box, label, icon}
}

var statsThresholds = utils.Threshold[int, string]{
	90: "critical",
	70: "warning",
}

func (s *systemInfo) SetValue(value int) {
	s.label.SetText(fmt.Sprintf("%2d%%", value))

	for _, text := range statsThresholds {
		ui.RemoveCSSClass(&s.box.Widget, text)
	}

	threshold, _ := statsThresholds.GetThreshold(value)
	ui.AddCSSClass(&s.box.Widget, threshold)
}
