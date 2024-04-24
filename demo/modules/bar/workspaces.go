package bar

import (
	"fmt"
	"webflo-dev/apero/services/hyprland"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type workspacesHandler struct {
	hyprland.HyprlandEventHandler
	workspaces map[int]*gtk.Button
}

func newWorkspaces() *gtk.Box {
	ids := 5

	box, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 8)
	box.SetName("workspaces")
	ui.AddCSSClass(&box.Widget, "workspaces")

	var workspaces = make(map[int]*gtk.Button, 5)
	for i := range ids {
		id := i + 1
		workspace := newWorkspace(id)
		box.Add(workspace)
		workspaces[id] = workspace
	}

	handler := &workspacesHandler{
		workspaces: workspaces,
	}

	handler.WorkspaceV2(hyprland.ActiveWorkspace().Id, hyprland.ActiveWorkspace().Name)
	hyprland.WatchEvents(handler, hyprland.EventWorkspacev2)

	return box
}

func newWorkspace(id int) *gtk.Button {
	button, _ := gtk.ButtonNew()
	image := ui.NewFontSizeImageFromIconName(fmt.Sprintf("__workspace-%d-symbolic", id))
	button.Add(image)

	ui.AddCSSClass(&button.Widget, "workspace")

	button.Connect("clicked", func() {
		hyprland.Dispatch("workspace %d", id)
	})

	return button
}

func (handler *workspacesHandler) WorkspaceV2(activeId int, name string) {
	glib.IdleAdd(func() {

		workspacesFromHyprland := hyprland.Workspaces()

		for id, workspace := range handler.workspaces {
			if id == activeId {
				ui.AddCSSClass(&workspace.Widget, "active")
			} else {
				ui.RemoveCSSClass(&workspace.Widget, "active")
				for _, whl := range workspacesFromHyprland {
					if whl.Id == id {
						if whl.Windows != 0 {
							ui.AddCSSClass(&workspace.Widget, "occupied")
						} else {
							ui.RemoveCSSClass(&workspace.Widget, "occupied")
						}
					}
				}
			}
		}
	})
}
