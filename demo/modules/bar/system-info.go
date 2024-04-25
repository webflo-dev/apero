package bar

import (
	"fmt"
	"math"
	systemStats "webflo-dev/apero/services/system-stats"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type systemInfo struct {
	box   *gtk.Box
	label *gtk.Label
	icon  *gtk.Image
}

type systemStatsHandler struct {
	systemStats.SystemStatsEventHandler
	cpu    *systemInfo
	memory *systemInfo
	nvidia *systemInfo
}

func newSystemInfoModule() *gtk.Box {

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 24)
	box.SetName("system-info")

	cpu := newSystemInfoBox("_processor-symbolic")
	box.Add(cpu.box)

	memory := newSystemInfoBox("_memory-symbolic")
	box.Add(memory.box)

	nvidia := newSystemInfoBox("_gpu-symbolic")
	box.Add(nvidia.box)

	systemStatsHandler := &systemStatsHandler{
		cpu:    cpu,
		memory: memory,
		nvidia: nvidia,
	}
	systemStats.WatchSystemStats(systemStatsHandler)

	return box
}

func (handler *systemStatsHandler) Notify(value *systemStats.SystemStats) {
	glib.IdleAdd(func() {
		handler.cpu.SetValue(value.Cpu.Usage)
		handler.memory.SetValue(int(math.Floor((float64(value.Memory.Used) / float64(value.Memory.Total)) * 100)))
		handler.nvidia.SetValue(value.Nvidia.GpuUsage)
	})
}

func newSystemInfoBox(iconName string) *systemInfo {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	ui.AddCSSClass(&box.Widget, "info")

	icon, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
	box.Add(icon)

	label, _ := gtk.LabelNew("")
	box.Add(label)

	return &systemInfo{box, label, icon}
}

type threshold struct {
	value int
	level string
}

var thresholds = []threshold{
	{90, "critical"},
	{70, "warning"},
}

func getLevel(value int) string {
	for _, threshold := range thresholds {
		if value >= threshold.value {
			return threshold.level
		}
	}
	return ""
}

func (s *systemInfo) SetValue(value int) {
	s.label.SetText(fmt.Sprintf("%2d%%", value))

	for _, threshold := range thresholds {
		ui.RemoveCSSClass(&s.box.Widget, threshold.level)
	}
	ui.AddCSSClass(&s.box.Widget, getLevel(value))
}
