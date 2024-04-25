package bar

import (
	"log"
	"webflo-dev/apero/services/hyprland"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/gtk"
)

type windowStateHandler struct {
	hyprland.HyprlandEventHandler
	box        *gtk.Box
	pinned     *gtk.Image
	fullscreen *gtk.Image
	floating   *gtk.Image
}

func newWindowInfoModule() *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	box.SetName("window-info")

	pinned := newWindowInfo("__window-state-pinned-symbolic", "pinned")
	fullscreen := newWindowInfo("__window-state-fullscreen-symbolic", "fullscreen")
	floating := newWindowInfo("__window-state-floating-symbolic", "floating")

	box.Add(fullscreen)
	box.Add(floating)
	box.Add(pinned)

	handler := &windowStateHandler{
		box:        box,
		pinned:     pinned,
		fullscreen: fullscreen,
		floating:   floating,
	}
	hyprland.WatchEvents(handler,
		hyprland.EventActiveWindowv2,
		hyprland.EventFullscreen,
		hyprland.EventChangeFloatingMode,
		hyprland.EventPin,
	)

	return box
}

func newWindowInfo(iconName string, className string) *gtk.Image {
	icon := ui.NewFontSizeImageFromIconName(iconName)

	ui.AddCSSClass(&icon.Widget, "indicator")
	if className != "" {
		ui.AddCSSClass(&icon.Widget, className)
	}

	return icon
}

func (handler *windowStateHandler) ActiveWindowV2(windowAddress string) {
	activeClient := hyprland.ActiveClient()
	handler.Pin(windowAddress, activeClient.Pinned)
	handler.Fullscreen(activeClient.Fullscreen)
	handler.ChangeFloatingMode(windowAddress, activeClient.Floating)

	log.Println("ActiveWindowV2", windowAddress, activeClient.Address, activeClient.XWayland)
	if windowAddress == activeClient.Address {
		ui.ToggleCSSClass(&handler.box.Widget, "xwayland", activeClient.XWayland)
	}
}

func (handler *windowStateHandler) Fullscreen(fullscreen bool) {
	ui.ToggleCSSClass(&handler.fullscreen.Widget, "active", fullscreen)
}

func (handler *windowStateHandler) ChangeFloatingMode(windowAddress string, floating bool) {
	if windowAddress == hyprland.ActiveClient().Address {
		ui.ToggleCSSClass(&handler.floating.Widget, "active", floating)
	}
}

func (handler *windowStateHandler) Pin(windowAddress string, pinned bool) {
	if windowAddress == hyprland.ActiveClient().Address {
		ui.ToggleCSSClass(&handler.pinned.Widget, "active", pinned)
	}
}
