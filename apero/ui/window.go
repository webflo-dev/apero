package ui

import (
	"github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/gtk"
)

type Position int

const (
	PositionLeft   Position = 1
	PositionRight  Position = 2
	PositionTop    Position = 4
	PositionBottom Position = 8
)

var Positions = map[Position]layershell.LayerShellEdgeFlags{
	PositionLeft:   layershell.LAYER_SHELL_EDGE_LEFT,
	PositionRight:  layershell.LAYER_SHELL_EDGE_RIGHT,
	PositionTop:    layershell.LAYER_SHELL_EDGE_TOP,
	PositionBottom: layershell.LAYER_SHELL_EDGE_BOTTOM,
}

type Layer string

const (
	LayerTop        Layer = "Top"
	LayerBottom     Layer = "Bottom"
	LayerBackground Layer = "Background"
	LayerOverlay    Layer = "Overlay"
)

var Layers = map[Layer]layershell.LayerShellLayerFlags{
	LayerTop:        layershell.LAYER_SHELL_LAYER_TOP,
	LayerBottom:     layershell.LAYER_SHELL_LAYER_BOTTOM,
	LayerBackground: layershell.LAYER_SHELL_LAYER_BACKGROUND,
	LayerOverlay:    layershell.LAYER_SHELL_LAYER_OVERLAY,
}

type ExclusiveZone string

const (
	ExclusiveZoneNormal ExclusiveZone = "normal"
	ExclusiveZoneIgnore ExclusiveZone = "ignore"
	ExclusiveZoneAuto   ExclusiveZone = "auto"
)

type WindowProps struct {
	Name          string
	Anchor        Position
	Layer         Layer
	ExclusiveZone ExclusiveZone
}

func NewWindow(props WindowProps) *gtk.Window {
	// if !layershell.IsSupported() {
	// 	log.Fatalln("gtk-layer-shell not supported")
	// }

	w, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	w.SetName(props.Name)

	layershell.InitForWindow(w)
	// layershell.SetNamespace(w, props.Name)

	SetAnchor(w, props.Anchor)
	SetLayer(w, LayerTop)
	SetExclusiveZone(w, props.ExclusiveZone)

	return w
}

func SetAnchor(w *gtk.Window, anchor Position) {
	for position, edge := range Positions {
		match := position&anchor != 0
		layershell.SetAnchor(w, edge, match)
	}
}

func SetMargin(w *gtk.Window, position Position, value int) {
	layershell.SetMargin(w, Positions[position], value)
}

func SetLayer(w *gtk.Window, layer Layer) {
	if value, ok := Layers[layer]; ok {
		layershell.SetLayer(w, value)
	}
}

func SetExclusiveZone(w *gtk.Window, zone ExclusiveZone) {
	if zone == ExclusiveZoneNormal {
		layershell.SetExclusiveZone(w, 0)
		return
	}
	if zone == ExclusiveZoneIgnore {
		layershell.SetExclusiveZone(w, -1)
		return
	}

	layershell.AutoExclusiveZoneEnable(w)
}
