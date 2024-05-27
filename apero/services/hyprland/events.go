package hyprland

import (
	"errors"
	"webflo-dev/apero/events"
)

//
// Event type
//

type Event string

const (
	EventWorkspace          Event = "workspace"
	EventWorkspaceV2        Event = "workspacev2"
	EventFocusedMonitor     Event = "focusedmonitor"
	EventActiveWindow       Event = "activewindow"
	EventActiveWindowV2     Event = "activewindowv2"
	EventFullscreen         Event = "fullscreen"
	EventMonitorRemoved     Event = "monitorremoved"
	EventMonitorAdded       Event = "monitoradded"
	EventMonitorAddedV2     Event = "monitoraddedv2"
	EventCreateWorkspace    Event = "createworkspace"
	EventCreateWorkspaceV2  Event = "createworkspacev2"
	EventDestroyWorkspace   Event = "destroyworkspace"
	EventDestroyWorkspaceV2 Event = "destroyworkspacev2"
	EventMoveWorkspace      Event = "moveworkspace"
	EventMoveWorkspaceV2    Event = "moveworkspacev2"
	EventRenameWorkspace    Event = "renameworkspace"
	EventActiveSpecial      Event = "activespecial"
	EventActiveLayout       Event = "activelayout"
	EventOpenWindow         Event = "openwindow"
	EventCloseWindow        Event = "closewindow"
	EventMoveWindow         Event = "movewindow"
	EventMoveWindowV2       Event = "movewindowv2"
	EventOpenLayer          Event = "openayer"
	EventCloseLayer         Event = "closelayer"
	EventSubMap             Event = "submap"
	EventChangeFloatingMode Event = "changefloatingmode"
	EventUrgent             Event = "urgent"
	EventMinimize           Event = "minimize"
	EventScreencast         Event = "screencast"
	EventWindowTitle        Event = "windowtitle"
	EventIgnoreGroupLock    Event = "ignoregrouplock"
	EventLockGroups         Event = "lockgroups"
	EventConfigReloaded     Event = "configreloaded"
	EventPin                Event = "pin"
)

var evts = map[Event]any{
	EventWorkspace:          newEvt[PayloadWorkspace](),
	EventWorkspaceV2:        newEvt[PayloadWorkspaceV2](),
	EventFocusedMonitor:     newEvt[PayloadFocusedMonitor](),
	EventActiveWindow:       newEvt[PayloadActiveWindow](),
	EventActiveWindowV2:     newEvt[PayloadActiveWindowV2](),
	EventFullscreen:         newEvt[PayloadFullscreen](),
	EventMonitorRemoved:     newEvt[PayloadMonitorRemoved](),
	EventMonitorAdded:       newEvt[PayloadMonitorAdded](),
	EventMonitorAddedV2:     newEvt[PayloadMonitorAddedV2](),
	EventCreateWorkspace:    newEvt[PayloadCreateWorkspace](),
	EventCreateWorkspaceV2:  newEvt[PayloadCreateWorkspaceV2](),
	EventDestroyWorkspace:   newEvt[PayloadDestroyWorkspace](),
	EventDestroyWorkspaceV2: newEvt[PayloadDestroyWorkspaceV2](),
	EventMoveWorkspace:      newEvt[PayloadMoveWorkspace](),
	EventMoveWorkspaceV2:    newEvt[PayloadMoveWorkspaceV2](),
	EventRenameWorkspace:    newEvt[PayloadRenameWorkspace](),
	EventActiveSpecial:      newEvt[PayloadActiveSpecial](),
	EventActiveLayout:       newEvt[PayloadActiveLayout](),
	EventOpenWindow:         newEvt[PayloadOpenWindow](),
	EventCloseWindow:        newEvt[PayloadCloseWindow](),
	EventMoveWindow:         newEvt[PayloadMoveWindow](),
	EventMoveWindowV2:       newEvt[PayloadMoveWindowV2](),
	EventOpenLayer:          newEvt[PayloadOpenLayer](),
	EventCloseLayer:         newEvt[PayloadCloseLayer](),
	EventSubMap:             newEvt[PayloadSubMap](),
	EventChangeFloatingMode: newEvt[PayloadChangeFloatingMode](),
	EventUrgent:             newEvt[PayloadUrgent](),
	EventMinimize:           newEvt[PayloadMinimize](),
	EventScreencast:         newEvt[PayloadScreencast](),
	EventWindowTitle:        newEvt[PayloadWindowTitle](),
	EventIgnoreGroupLock:    newEvt[PayloadIgnoreGroupLock](),
	EventLockGroups:         newEvt[PayloadLockGroups](),
	EventConfigReloaded:     newEvt[PayloadEmpty](),
	EventPin:                newEvt[PayloadPin](),
}

type HyprlandPayload interface {
	from(MsgValues)
}
type HyprlandEvtEmitter interface {
	publish(MsgValues)
}

type HyprlandEvt[P any] struct {
	events.Event[P]
}

func (e *HyprlandEvt[P]) publish(v MsgValues) {
	p := new(P)
	any(p).(HyprlandPayload).from(v)
	e.Publish(*p)
}

func newEvt[P any]() *HyprlandEvt[P] {
	return &HyprlandEvt[P]{
		Event: events.New[P](),
	}
}

func registerEvt[P any](event Event, handler func(P)) (events.ID, error) {
	if e, ok := evts[event].(*HyprlandEvt[P]); ok {
		return e.RegisterHandler(events.HandlerFunc[P](handler))
	}

	return 0, errors.New("invalid event type")
}

func publishEvt(event Event, values MsgValues) {
	if e, ok := evts[event].(HyprlandEvtEmitter); ok {
		e.publish(values)
	}
}

//
// Event payloads
//

type PayloadEmpty struct{}

type MsgValues []string

type PayloadWorkspace struct {
	WorkspaceName string
}

func (p *PayloadWorkspace) from(v MsgValues) {
	p.WorkspaceName = v[0]
}

type PayloadWorkspaceV2 struct {
	WorkspaceId   int
	WorkspaceName string
}

func (p *PayloadWorkspaceV2) from(v MsgValues) {
	p.WorkspaceId = toInt(v[0])
	p.WorkspaceName = v[1]
}

type PayloadFocusedMonitor struct {
	MonitorName   string
	WorkspaceName string
}

func (p *PayloadFocusedMonitor) from(v MsgValues) {
	p.MonitorName = v[0]
	p.WorkspaceName = v[1]
}

type PayloadActiveWindow struct {
	WindowClass string
	WindowTitle string
}

func (p *PayloadActiveWindow) from(v MsgValues) {
	p.WindowClass = v[0]
	p.WindowTitle = v[1]
}

type PayloadActiveWindowV2 struct {
	WindowAddress string
}

func (p *PayloadActiveWindowV2) from(v MsgValues) {
	p.WindowAddress = v[0]
}

type PayloadFullscreen struct {
	Fullscreen bool
}

func (p *PayloadFullscreen) from(v MsgValues) {
	p.Fullscreen = toBool(v[0])
}

type PayloadMonitorRemoved struct {
	MonitorName string
}

func (p *PayloadMonitorRemoved) from(v MsgValues) {
	p.MonitorName = v[0]
}

type PayloadMonitorAdded struct {
	MonitorName string
}

func (p *PayloadMonitorAdded) from(v MsgValues) {
	p.MonitorName = v[0]
}

type PayloadMonitorAddedV2 struct {
	MonitorId          int
	MonitorName        string
	MonitorDescription string
}

func (p *PayloadMonitorAddedV2) from(v MsgValues) {
	p.MonitorId = toInt(v[0])
	p.MonitorName = v[1]
	p.MonitorDescription = v[2]
}

type PayloadCreateWorkspace struct {
	WorkspaceName string
}

func (p *PayloadCreateWorkspace) from(v MsgValues) {
	p.WorkspaceName = v[0]
}

type PayloadCreateWorkspaceV2 struct {
	WorkspaceId   int
	WorkspaceName string
}

func (p *PayloadCreateWorkspaceV2) from(v MsgValues) {
	p.WorkspaceId = toInt(v[0])
	p.WorkspaceName = v[1]
}

type PayloadDestroyWorkspace struct {
	WorkspaceName string
}

func (p *PayloadDestroyWorkspace) from(v MsgValues) {
	p.WorkspaceName = v[0]
}

type PayloadDestroyWorkspaceV2 struct {
	WorkspaceId   int
	WorkspaceName string
}

func (p *PayloadDestroyWorkspaceV2) from(v MsgValues) {
	p.WorkspaceId = toInt(v[0])
	p.WorkspaceName = v[1]
}

type PayloadMoveWorkspace struct {
	WorkspaceName string
	MonitorName   string
}

func (p *PayloadMoveWorkspace) from(v MsgValues) {
	p.WorkspaceName = v[0]
	p.MonitorName = v[1]
}

type PayloadMoveWorkspaceV2 struct {
	WorkspaceId   int
	WorkspaceName string
	MonitorName   string
}

func (p *PayloadMoveWorkspaceV2) from(v MsgValues) {
	p.WorkspaceId = toInt(v[0])
	p.WorkspaceName = v[1]
	p.MonitorName = v[2]
}

type PayloadRenameWorkspace struct {
	WorkspaceId      int
	NewWorkspaceName string
}

func (p *PayloadRenameWorkspace) from(v MsgValues) {
	p.WorkspaceId = toInt(v[0])
	p.NewWorkspaceName = v[1]
}

type PayloadActiveSpecial struct {
	WorkspaceName string
	MonitorName   string
}

func (p *PayloadActiveSpecial) from(v MsgValues) {
	p.WorkspaceName = v[0]
	p.MonitorName = v[1]
}

type PayloadActiveLayout struct {
	KeyboardName string
	LayoutName   string
}

func (p *PayloadActiveLayout) from(v MsgValues) {
	p.KeyboardName = v[0]
	p.LayoutName = v[1]
}

type PayloadOpenWindow struct {
	WindowAddress string
	WorkspaceName string
	WindowClass   string
	WindowTitle   string
}

func (p *PayloadOpenWindow) from(v MsgValues) {
	p.WindowAddress = v[0]
	p.WorkspaceName = v[1]
	p.WindowClass = v[2]
	p.WindowTitle = v[3]
}

type PayloadCloseWindow struct {
	WindowAddress string
}

func (p *PayloadCloseWindow) from(v MsgValues) {
	p.WindowAddress = v[0]
}

type PayloadMoveWindow struct {
	WindowAddress string
	WorkspaceName string
}

func (p *PayloadMoveWindow) from(v MsgValues) {
	p.WindowAddress = v[0]
	p.WorkspaceName = v[1]
}

type PayloadMoveWindowV2 struct {
	WindowAddress string
	WorkspaceId   int
	WorkspaceName string
}

func (p *PayloadMoveWindowV2) from(v MsgValues) {
	p.WindowAddress = v[0]
	p.WorkspaceId = toInt(v[1])
	p.WorkspaceName = v[2]
}

type PayloadOpenLayer struct {
	Namespace string
}

func (p *PayloadOpenLayer) from(v MsgValues) {
	p.Namespace = v[0]
}

type PayloadCloseLayer struct {
	Namespace string
}

func (p *PayloadCloseLayer) from(v MsgValues) {
	p.Namespace = v[0]
}

type PayloadSubMap struct {
	SubmapName string
}

func (p *PayloadSubMap) from(v MsgValues) {
	p.SubmapName = v[0]
}

type PayloadChangeFloatingMode struct {
	WindowAddress string
	Floating      bool
}

func (p *PayloadChangeFloatingMode) from(v MsgValues) {
	p.WindowAddress = v[0]
	p.Floating = toBool(v[1])
}

type PayloadUrgent struct {
	WindowAddress string
}

func (p *PayloadUrgent) from(v MsgValues) {
	p.WindowAddress = v[0]
}

type PayloadMinimize struct {
	WindowAddress string
	Minimized     bool
}

func (p *PayloadMinimize) from(v MsgValues) {
	p.WindowAddress = v[0]
	p.Minimized = toBool(v[1])
}

type ScreencastOwner uint16

const (
	ScreencastOwnerMonitor ScreencastOwner = iota
	ScreencastOwnerWindow
)

type PayloadScreencast struct {
	Screencasting bool
	Owner         ScreencastOwner
}

func (p *PayloadScreencast) from(v MsgValues) {
	p.Screencasting = toBool(v[0])
	p.Owner = ScreencastOwner(toInt(v[1]))
}

type PayloadWindowTitle struct {
	WindowAddress string
}

func (p *PayloadWindowTitle) from(v MsgValues) {
	p.WindowAddress = v[0]
}

type PayloadIgnoreGroupLock struct {
	IgnoringGrouplock bool
}

func (p *PayloadIgnoreGroupLock) from(v MsgValues) {
	p.IgnoringGrouplock = toBool(v[0])
}

type PayloadLockGroups struct {
	LockingGroups bool
}

func (p *PayloadLockGroups) from(v MsgValues) {
	p.LockingGroups = toBool(v[0])
}

type PayloadPin struct {
	WindowAddress string
	Pinned        bool
}

func (p *PayloadPin) from(v MsgValues) {
	p.WindowAddress = v[0]
	p.Pinned = toBool(v[1])
}
