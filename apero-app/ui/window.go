package ui

import (
	"log"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type GtkWindow struct {
	*gtk.Window
}

type Position int
const (
	PositionLeft Position = 1
	PositionRight Position = 2
	PositionTop Position = 4
	PositionBottom Position = 8
)
var Positions = map[Position]gtk4layershell.Edge{
	PositionLeft: gtk4layershell.LayerShellEdgeLeft,
	PositionRight: gtk4layershell.LayerShellEdgeRight,
	PositionTop: gtk4layershell.LayerShellEdgeTop,
	PositionBottom: gtk4layershell.LayerShellEdgeBottom,
}


type Layer string
const  (
	LayerTop Layer = "Top"
	LayerBottom Layer = "Bottom"
	LayerBackground Layer = "Background"
	LayerOverlay Layer = "Overlay"
)
var Layers = map[Layer]gtk4layershell.Layer{
	LayerTop: gtk4layershell.LayerShellLayerTop,
	LayerBottom: gtk4layershell.LayerShellLayerBottom,
	LayerBackground: gtk4layershell.LayerShellLayerBackground,
	LayerOverlay: gtk4layershell.LayerShellLayerOverlay,
}

type ExclusiveZone string
const (
	ExclusiveZoneNormal ExclusiveZone = "normal"
	ExclusiveZoneIgnore ExclusiveZone = "ignore"
	ExclusiveZoneAuto ExclusiveZone = "auto"
)


type WindowProps struct {
	name string
	anchor Position
	layer Layer
	exclusiveZone ExclusiveZone
}


func NewWindow(props WindowProps) GtkWindow {
	if !gtk4layershell.IsSupported() {
		log.Fatalln("gtk-layer-shell not supported")
	}

	w := GtkWindow{gtk.NewWindow()}

	w.SetName(props.name)
	
	gtk4layershell.InitForWindow(w.Window)
	gtk4layershell.SetNamespace(w.Window, props.name)

	w.SetAnchor(props.anchor)
	w.SetLayer(LayerTop)
	w.SetExclusiveZone(props.exclusiveZone)

	return w
}

func (w GtkWindow) SetAnchor(anchor Position) {
	for position, edge := range Positions {
		match := position & anchor != 0
		gtk4layershell.SetAnchor(w.Window, edge, match)
	}
}

func (w GtkWindow) SetMargin(position Position, value int) {
	gtk4layershell.SetMargin(w.Window, Positions[position], value)
}


func (w GtkWindow) SetLayer(layer Layer) {
	if value, ok := Layers[layer]; ok {
		gtk4layershell.SetLayer(w.Window, value)
	}
}

func (w GtkWindow) SetExclusiveZone(zone ExclusiveZone) {
	if (zone == ExclusiveZoneNormal) {
		gtk4layershell.SetExclusiveZone(w.Window, 0)
		return
	}
	if (zone == ExclusiveZoneIgnore) {
		gtk4layershell.SetExclusiveZone(w.Window, -1)
		return
	}

	gtk4layershell.AutoExclusiveZoneEnable(w.Window)
}