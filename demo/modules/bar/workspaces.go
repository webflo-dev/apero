package bar

import (
	"fmt"
	"slices"
	"webflo-dev/apero/services/hyprland"
	"webflo-dev/apero/ui"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type workspacesHandler struct {
	// hyprland.Subscriber
	workspaces map[int]*gtk.Button
}

func newWorkspacesModule() *gtk.Box {
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

	hyprland.OnWorkspacev2("bar-workspace", func(payload hyprland.PayloadWorkspaceV2) {
		handler.WorkspaceV2(payload.WorkspaceId, payload.WorkspaceName)
	})
	// hyprland.Register(handler, hyprland.EventWorkspaceV2, hyprland.EventUrgent)

	return box
}

func newWorkspace(id int) *gtk.Button {
	button, _ := gtk.ButtonNew()
	image := ui.NewFontSizeImageFromIconName(fmt.Sprintf(Icon_Workspace_pattern, id))
	button.Add(image)

	ui.AddCSSClass(&button.Widget, "workspace")

	button.Connect("clicked", func() {
		hyprland.Dispatch("workspace", id)
	})

	return button
}

func (handler *workspacesHandler) WorkspaceV2(workspaceId int, name string) {
	glib.IdleAdd(func() {

		workspacesFromHyprland := hyprland.Workspaces()
		for id, workspace := range handler.workspaces {
			ui.ToggleCSSClassFromBool(&workspace.Widget, "active", id == workspaceId)

			if id == workspaceId {
				if ui.HasCSSClass(&workspace.Widget, "urgent") {
					ui.RemoveCSSClass(&workspace.Widget, "urgent")
				}
			} else {
				for _, whl := range workspacesFromHyprland {
					if whl.Id == id {
						ui.ToggleCSSClassFromBool(&workspace.Widget, "occupied", whl.Windows != 0)
					}
				}
			}
		}
	})
}

func (handler *workspacesHandler) Urgent(windowAddress string) {

	clients := hyprland.Clients()
	clientIndex := slices.IndexFunc(clients, func(client hyprland.Client) bool {
		return client.Address == windowAddress
	})

	if clientIndex == -1 {
		return
	}

	client := clients[clientIndex]

	for id, workspace := range handler.workspaces {
		ui.ToggleCSSClassFromBool(&workspace.Widget, "urgent", id == client.Workspace.Id)
	}
}
