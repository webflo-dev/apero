package notifications

import (
	"log"
	"os"
	"reflect"
	"time"
)

type hints map[string]any
type actions map[string]string

type Notification struct {
	id           uint32
	appName      string
	appIcon      string
	summary      string
	body         string
	actions      actions
	hints        hints
	timestamp    int64
	urgency      urgency
	resident     bool
	desktopEntry string
	category     category
	imagePath    string
}

type imageData struct {
	width         int32
	height        int32
	rowStride     int32
	hasAlpha      bool
	bitsPerSample int32
	channels      int32
	data          []byte
}

type urgency uint16

const (
	urgencyLow      urgency = 0
	urgencyCritical urgency = 2
	urgencyNormal   urgency = 1
)

type category string

const (
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

func newNotification(id uint32, appName string, appIcon string, summary string, body string, actions []string, hints hints) (n Notification) {
	n = Notification{
		id:        id,
		appName:   appName,
		appIcon:   appIcon,
		summary:   summary,
		body:      body,
		hints:     hints,
		timestamp: time.Now().Unix(),
		urgency:   hints.urgency(),
		resident:  hints.resident(),
	}

	n.actions = getActions(actions)
	n.desktopEntry, _ = hints.desktopEntry()
	n.category, _ = hints.category()
	n.imagePath = getImagePath(n)

	return
}

// func (n Notification) InvokeAction(actionKey string) {
// 	if action, ok := n.actions[actionKey]; ok {
// 		server.InvokeAction(n.id, action)
// 	}
// }

func getActions(rawActions []string) actions {
	pair := make(actions)
	for i := 0; i < len(rawActions); i += 2 {
		if rawActions[i+1] == "" {
			pair[rawActions[i]] = rawActions[i+1]
		}
	}
	return pair
}

func getImagePath(n Notification) string {
	if n.appIcon != "" {
		return ""
	}

	if _, err := os.Stat(n.appIcon); err == nil {
		return n.appIcon
	} else {
		imageData, ok := n.hints.imageData()
		if ok {
			out, err := os.CreateTemp("", "apero-notification-*")
			if err != nil {
				log.Println("error creating temp file", err)
				return ""
			}
			_, err = out.Write(imageData.data)
			defer out.Close()
			return out.Name()
		} else {
			return n.hints.imagePath()
		}
	}
}

func (h hints) urgency() urgency {
	value := h["urgency"]
	if value == nil {
		return urgencyNormal
	}
	return urgency(reflect.ValueOf(value).Uint())
}

func (h hints) resident() bool {
	value := h["resident"]
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Bool()
}

func (h hints) desktopEntry() (string, bool) {
	value := h["desktop-entry"]
	if value == nil {
		return "", false
	}
	return reflect.ValueOf(value).String(), true
}

func (h hints) category() (category, bool) {
	value := h["category"]
	if value == nil {
		return "", false
	}
	return category(reflect.ValueOf(value).String()), true
}

func (h hints) senderPID() (uint64, bool) {
	value := h["sender-pid"]
	if value == nil {
		return 0, false
	}
	return reflect.ValueOf(value).Uint(), true
}

func (h hints) actionIcons() bool {
	value := h["action-icons"]
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Bool()
}

func (h hints) transient() bool {
	value := h["transient"]
	if value == nil {
		return false
	}
	return reflect.ValueOf(value).Bool()
}

func (h hints) x() (int32, bool) {
	value := h["x"]
	if value == nil {
		return 0, false
	}
	return int32(reflect.ValueOf(value).Int()), true
}

func (h hints) y() (int32, bool) {
	value := h["y"]
	if value == nil {
		return 0, false
	}
	return int32(reflect.ValueOf(value).Int()), true
}

func (h hints) imagePath() string {
	value := h["image-path"]
	if value == nil {
		return ""
	}
	return reflect.ValueOf(value).String()
}

func (h hints) imageData() (imageData, bool) {
	value := h["image-data"]
	imageData := imageData{}
	if value == nil {
		return imageData, false
	}
	values := reflect.ValueOf(value).Interface().([]any)
	imageData.width = values[0].(int32)
	imageData.height = values[1].(int32)
	imageData.rowStride = values[2].(int32)
	imageData.hasAlpha = values[3].(bool)
	imageData.bitsPerSample = values[4].(int32)
	imageData.channels = values[5].(int32)
	imageData.data = values[6].([]byte)

	// bits_per_sample (i): Must always be 8
	if imageData.bitsPerSample != 8 {
		return imageData, false
	}

	// channels (i): If has_alpha is TRUE, must be 4, otherwise 3
	if imageData.hasAlpha && imageData.channels == 4 {
		return imageData, true
	}
	if imageData.channels != 3 {
		return imageData, false
	}

	return imageData, true
}
