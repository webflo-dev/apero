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
	pinned     *gtk.Button
	fullscreen *gtk.Button
	floating   *gtk.Button
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

func newWindowInfo(iconName string, className string) *gtk.Button {
	icon := ui.NewFontSizeImageFromIconName(iconName)
	ui.AddCSSClass(&icon.Widget, className)

	box, _ := gtk.ButtonNew()
	box.Add(icon)

	box.Connect("clicked", func() {
		switch className {
		case "pinned":
			hyprland.Dispatch("pin address:%s", hyprland.ActiveClient().Address)
		case "fullscreen":
			hyprland.Dispatch("fullscreen 1")
		case "floating":
			hyprland.Dispatch("togglefloating address:%s", hyprland.ActiveClient().Address)
		}
	})

	return box
}

func (handler *windowStateHandler) ActiveWindowV2(windowAddress string) {
	activeClient := hyprland.ActiveClient()
	handler.Pin(windowAddress, activeClient.Pinned)
	handler.Fullscreen(activeClient.Fullscreen)
	handler.ChangeFloatingMode(windowAddress, activeClient.Floating)

	log.Println("ActiveWindowV2", windowAddress, activeClient.Address, activeClient.XWayland)
	if windowAddress == activeClient.Address {
		ui.ToggleCSSClassFromBool(&handler.box.Widget, "xwayland", activeClient.XWayland)
	}
}

func (handler *windowStateHandler) Fullscreen(fullscreen bool) {
	ui.ToggleCSSClassFromBool(&handler.fullscreen.Widget, "active", fullscreen)
}

func (handler *windowStateHandler) ChangeFloatingMode(windowAddress string, floating bool) {
	if windowAddress == hyprland.ActiveClient().Address {
		ui.ToggleCSSClassFromBool(&handler.floating.Widget, "active", floating)
	}
}

func (handler *windowStateHandler) Pin(windowAddress string, pinned bool) {
	if windowAddress == hyprland.ActiveClient().Address {
		ui.ToggleCSSClassFromBool(&handler.pinned.Widget, "active", pinned)
	}
}
