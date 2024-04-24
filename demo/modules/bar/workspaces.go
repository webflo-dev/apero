package bar

import (
	"fmt"
	"webflo-dev/apero/services/hyprland"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type handler struct {
	hyprland.EventHandler
	workspaces map[int]*gtk.Button
}

func newWorkspaces() *gtk.Box {
	ids := 5

	box := gtk.NewBox(gtk.OrientationHorizontal, 8)
	box.SetName("workspaces")
	box.SetCSSClasses([]string{"workspaces"})

	var workspaces = make(map[int]*gtk.Button, 5)
	for i := range ids {
		id := i + 1
		workspace := newWorkspace(id)
		box.Append(workspace)
		workspaces[id] = workspace
	}

	hyprland.WatchEvents(&handler{
		workspaces: workspaces,
	}, hyprland.EventWorkspacev2)

	return box
}

func newWorkspace(id int) *gtk.Button {
	button := gtk.NewButton()
	button.AddCSSClass("workspace")

	// button.SetIconName(fmt.Sprintf("workspace-%d-empty", id))
	label := gtk.NewLabel(fmt.Sprintf("%d", id))
	label.AddCSSClass("circular")
	button.SetChild(label)
	button.SetVAlign(gtk.AlignCenter)

	button.ConnectClicked(func() {
		hyprland.Dispatch(fmt.Sprintf("workspace %d", id))
	})

	return button
}

func (h *handler) WorkspaceV2(activeId int, name string) {

	workspacesFromHyprland, _ := hyprland.Workspaces()

	for id, workspace := range h.workspaces {

		if id == activeId {
			workspace.AddCSSClass("active")
			workspace.SetIconName(fmt.Sprintf("workspace-%d-active", id))
		} else {
			workspace.RemoveCSSClass("active")
			workspace.SetIconName(fmt.Sprintf("workspace-%d-empty", id))
			for _, whl := range workspacesFromHyprland {
				if whl.Id == id {
					if whl.Windows != 0 {
						workspace.SetIconName(fmt.Sprintf("workspace-%d-occupied", id))
						workspace.AddCSSClass("occupied")
					} else {
						workspace.RemoveCSSClass("occupied")
					}
				}
			}
		}
	}
}
