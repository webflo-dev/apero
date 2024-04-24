package bar

import (
	"fmt"
	"math"
	"webflo-dev/apero/services/stats"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type systemInfo struct {
	box   *gtk.Box
	label *gtk.Label
	icon  *gtk.Image
}

const (
	defaultSystemInfoClass = "info"
)

func newSystemInfo() *gtk.Box {

	box := gtk.NewBox(gtk.OrientationHorizontal, 24)
	box.SetName("system-info")

	cpu := newSystemInfoBox("_processor-symbolic")
	box.Append(cpu.box)

	memory := newSystemInfoBox("_memory-symbolic")
	box.Append(memory.box)

	nvidia := newSystemInfoBox("_gpu-symbolic")
	box.Append(nvidia.box)

	statsChan := stats.WatchStats()
	go func() {
		for {
			value := <-statsChan
			glib.IdleAdd(func() {
				cpu.SetValue(value.Cpu.Usage)
				memory.SetValue(int(math.Floor((float64(value.Memory.Used) / float64(value.Memory.Total)) * 100)))
				nvidia.SetValue(value.Nvidia.GpuUsage)
			})
		}
	}()

	return box
}

func newSystemInfoBox(iconName string) *systemInfo {
	box := gtk.NewBox(gtk.OrientationHorizontal, 8)
	box.SetCSSClasses([]string{defaultSystemInfoClass})

	icon := gtk.NewImageFromIconName(iconName)
	box.Append(icon)

	label := gtk.NewLabel("")
	box.Append(label)

	return &systemInfo{box, label, icon}
}

func (s *systemInfo) SetValue(text interface{}) {
	s.label.SetText(fmt.Sprintf("%2d%%", text))
}

func (s *systemInfo) ToggleClass(className string) {
	if className == "" {
		s.box.SetCSSClasses([]string{defaultSystemInfoClass})
	} else {
		s.box.AddCSSClass(className)
	}
}
