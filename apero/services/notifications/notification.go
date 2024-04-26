package notifications

import (
	"errors"
	"os"
	"reflect"
)

type hints map[string]any

type notification struct {
	id        int
	appName   string
	appIcon   string
	summary   string
	body      string
	actions   []string
	hints     hints
	timestamp int64
}

type urgency uint

const (
	urgencyLow      urgency = 0
	urgencyCritical urgency = 2
	urgencyNormal   urgency = 1
)

type category string

const (
	categoryNone                category = ""
	categoryDevice              category = "device"
	categoryDeviceAdded         category = "device.added"
	categoryDeviceError         category = "device.error"
	categoryDeviceRemoved       category = "device.removed"
	categoryEmail               category = "email"
	categoryEmailArrived        category = "email.arrived"
	categoryEmailBounced        category = "email.bounced"
	categoryIm                  category = "im"
	categoryImError             category = "im.error"
	categoryImReceived          category = "im.received"
	categoryNetwork             category = "network"
	categoryNetworkConnected    category = "network.connected"
	categoryNetworkDisconnected category = "network.disconnected"
	categoryNetworkError        category = "network.error"
	categoryPresence            category = "presence"
	categoryPresenceOffline     category = "presence.offline"
	categoryPresenceOnline      category = "presence.online"
	categoryTransfer            category = "transfer"
	categoryTransferComplete    category = "transfer.complete"
	categoryTransferError       category = "transfer.error"
)

func (n notification) urgency() urgency {
	value := n.hints["urgency"]
	if value == nil {
		return urgencyNormal
	}
	return urgency(reflect.ValueOf(value).Uint())
}

func (n notification) resident() bool {
	value := n.hints["resident"]
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Bool()
}

func (n notification) senderPID() (uint64, bool) {
	value := n.hints["sender-pid"]
	if value == nil {
		return 0, false
	}
	return reflect.ValueOf(value).Uint(), true
}

func (n notification) category() (category, bool) {
	value := n.hints["category"]
	if value == nil {
		return "", false
	}
	return category(reflect.ValueOf(value).String()), true
}

func (n notification) actionIcons() bool {
	value := n.hints["action-icons"]
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Bool()
}

func (n notification) desktopEntry() (string, bool) {
	value := n.hints["desktop-entry"]
	if value == nil {
		return "", false
	}
	return reflect.ValueOf(value).String(), true
}

func (n notification) transient() bool {
	value := n.hints["transient"]
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Bool()
}

func (n notification) x() (int32, bool) {
	value := n.hints["x"]
	if value == nil {
		return 0, false
	}
	return int32(reflect.ValueOf(value).Int()), true
}

func (n notification) y() (int32, bool) {
	value := n.hints["y"]
	if value == nil {
		return 0, false
	}
	return int32(reflect.ValueOf(value).Int()), true
}

func (n notification) image() (string, bool) {
	if n.appIcon != "" {
		if _, err := os.Stat(n.appIcon); errors.Is(err, os.ErrNotExist) {
			return n.appIcon, false
		}
		return n.appIcon, true
	}

	value := n.hints["image-path"]
	if value == nil {
		return "", false
	}
	return reflect.ValueOf(value).String(), true
}
