package ui

import (
	"log"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// type GtkWindow struct {
// 	*gtk.Window
// }

type Position int

const (
	PositionLeft   Position = 1
	PositionRight  Position = 2
	PositionTop    Position = 4
	PositionBottom Position = 8
)

var Positions = map[Position]gtk4layershell.Edge{
	PositionLeft:   gtk4layershell.LayerShellEdgeLeft,
	PositionRight:  gtk4layershell.LayerShellEdgeRight,
	PositionTop:    gtk4layershell.LayerShellEdgeTop,
	PositionBottom: gtk4layershell.LayerShellEdgeBottom,
}

type Layer string

const (
	LayerTop        Layer = "Top"
	LayerBottom     Layer = "Bottom"
	LayerBackground Layer = "Background"
	LayerOverlay    Layer = "Overlay"
)

var Layers = map[Layer]gtk4layershell.Layer{
	LayerTop:        gtk4layershell.LayerShellLayerTop,
	LayerBottom:     gtk4layershell.LayerShellLayerBottom,
	LayerBackground: gtk4layershell.LayerShellLayerBackground,
	LayerOverlay:    gtk4layershell.LayerShellLayerOverlay,
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
	if !gtk4layershell.IsSupported() {
		log.Fatalln("gtk-layer-shell not supported")
	}

	w := gtk.NewWindow()

	w.SetName(props.Name)

	gtk4layershell.InitForWindow(w)
	gtk4layershell.SetNamespace(w, props.Name)

	SetAnchor(w, props.Anchor)
	SetLayer(w, LayerTop)
	SetExclusiveZone(w, props.ExclusiveZone)

	return w
}

func SetAnchor(w *gtk.Window, anchor Position) {
	for position, edge := range Positions {
		match := position&anchor != 0
		gtk4layershell.SetAnchor(w, edge, match)
	}
}

func SetMargin(w *gtk.Window, position Position, value int) {
	gtk4layershell.SetMargin(w, Positions[position], value)
}

func SetLayer(w *gtk.Window, layer Layer) {
	if value, ok := Layers[layer]; ok {
		gtk4layershell.SetLayer(w, value)
	}
}

func SetExclusiveZone(w *gtk.Window, zone ExclusiveZone) {
	if zone == ExclusiveZoneNormal {
		gtk4layershell.SetExclusiveZone(w, 0)
		return
	}
	if zone == ExclusiveZoneIgnore {
		gtk4layershell.SetExclusiveZone(w, -1)
		return
	}

	gtk4layershell.AutoExclusiveZoneEnable(w)
}
