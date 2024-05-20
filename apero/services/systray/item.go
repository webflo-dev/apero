package systray

import (
	"strings"

	"github.com/godbus/dbus/v5"
)

type ItemObserver interface {
	ItemUpdated(name string)
}

const (
	itemIface = "org.kde.StatusNotifierItem"
)

type itemStatus string

const (
	ItemStatusActive         itemStatus = "Active"
	ItemStatusPassive        itemStatus = "Passive"
	ItemStatusNeedsAttention itemStatus = "NeedsAttention"
)

type itemCategory string

const (
	ItemCategoryApplicationStatus itemCategory = "ApplicationStatus"
	ItemCategoryCommunications    itemCategory = "Communications"
	ItemCategorySystemServices    itemCategory = "SystemServices"
	ItemCategoryHardware          itemCategory = "Hardware"
)

type itemProp string

const (
	ItemPropAttentionIconName   itemProp = "AttentionIconName"
	ItemPropAttentionIconPixmap itemProp = "AttentionIconPixmap"
	ItemPropAttentionMovieName  itemProp = "AttentionMovieName"
	ItemPropCategory            itemProp = "Category"
	ItemPropIconPixmap          itemProp = "IconPixmap"
	ItemPropIconName            itemProp = "IconName"
	ItemPropId                  itemProp = "Id"
	ItemPropOverlayIconName     itemProp = "OverlayIconName"
	ItemPropOverlayIconPixmap   itemProp = "OverlayIconPixmap"
	ItemPropStatus              itemProp = "Status"
	ItemPropTitle               itemProp = "Title"
	ItemPropToolTip             itemProp = "ToolTip"
	ItemPropWindowId            itemProp = "WindowId"
)

const (
	signalItemNewIcon          = itemIface + ".NewIcon"
	signalItemNewTitle         = itemIface + ".NewTitle"
	signalItemNewAttentionIcon = itemIface + ".NewAttentionIcon"
	signalItemNewOverlayIcon   = itemIface + ".NewOverlayIcon"
	signalItemNewToolTip       = itemIface + ".NewToolTip"
	signalItemNewStatus        = itemIface + ".NewStatus"
)

type Pixmap struct {
	width  int32
	height int32
	data   []byte
}

func (p *Pixmap) Width() int32 {
	return p.width
}
func (p *Pixmap) Height() int32 {
	return p.height
}
func (p *Pixmap) Data() []byte {
	return p.data
}

type ItemTooltip struct {
	name        string
	data        Pixmap
	title       string
	description string
}

func (t *ItemTooltip) Name() string {
	return t.name
}
func (t *ItemTooltip) Data() Pixmap {
	return t.data
}
func (t *ItemTooltip) Title() string {
	return t.title
}
func (t *ItemTooltip) Description() string {
	return t.description
}

type item struct {
	attentionIconName   string
	attentionIconPixmap Pixmap
	attentionMovieName  string
	category            itemCategory
	iconName            string
	iconPixmap          Pixmap
	id                  string
	overlayIconName     string
	overlayIconPixmap   Pixmap
	status              itemStatus
	title               string
	tooltip             ItemTooltip
	windowId            uint32
}

func (p *Pixmap) update(input [][]any) {
	if len(input) != 1 {
		return
	}
	values := input[0]
	// log.Printf("newFromData >> %+v\n", values)

	if len(values) != 3 {
		return
	}

	p.width = values[0].(int32)
	p.height = values[1].(int32)
	p.data = values[2].([]byte)
	return
}

func (t *ItemTooltip) update(values []any) {
	if len(values) != 4 {
		return
	}

	t.name = values[0].(string)
	t.title = values[2].(string)
	t.description = values[3].(string)
	return
}

func (s *item) updateProps(values map[string]any) {
	if val, ok := values[string(ItemPropAttentionIconName)]; ok {
		s.attentionIconName = val.(string)
	}

	if val, ok := values[string(ItemPropAttentionIconPixmap)]; ok {
		arr := val.([][]any)
		if len(arr) == 1 {
			s.attentionIconPixmap.update(val.([][]any))
		}
	}

	if val, ok := values[string(ItemPropAttentionMovieName)]; ok {
		s.attentionMovieName = val.(string)
	}

	if val, ok := values[string(ItemPropCategory)]; ok {
		s.category = itemCategory(val.(string))
	}

	if val, ok := values[string(ItemPropIconName)]; ok {
		s.iconName = val.(string)
	}

	if val, ok := values[string(ItemPropIconPixmap)]; ok {
		s.iconPixmap.update(val.([][]any))
	}

	if val, ok := values[string(ItemPropId)]; ok {
		s.id = val.(string)
	}

	if val, ok := values[string(ItemPropOverlayIconName)]; ok {
		s.overlayIconName = val.(string)
	}

	if val, ok := values[string(ItemPropOverlayIconPixmap)]; ok {
		s.overlayIconPixmap.update(val.([][]any))
	}

	if val, ok := values[string(ItemPropStatus)]; ok {
		s.status = itemStatus(val.(string))
	}

	if val, ok := values[string(ItemPropTitle)]; ok {
		s.title = val.(string)
	}

	if val, ok := values[string(ItemPropToolTip)]; ok {
		s.tooltip.update(val.([]any))
	}

	if val, ok := values[string(ItemPropWindowId)]; ok {
		s.windowId = uint32(val.(int32))
	}
}

type SystrayItem struct {
	conn   *dbus.Conn
	obj    dbus.BusObject
	sender string
	path   dbus.ObjectPath
	item   *item
}

func newSysTrayItem(sender string, path string) *SystrayItem {
	correctPath := path
	if strings.HasPrefix(path, "/") == false {
		correctPath = "/StatusNotifierItem"
	}

	return &SystrayItem{
		sender: sender,
		path:   dbus.ObjectPath(correctPath),
	}
}

func (s *SystrayItem) Register(observer ItemObserver) (err error) {
	s.conn, err = dbus.ConnectSessionBus()

	s.updateProps()
	// log.Printf("GetAll (item) >> %+v\n", item)

	go func() {
		defer s.conn.Close()

		s.conn.AddMatchSignal(
			dbus.WithMatchObjectPath(s.path),
			dbus.WithMatchSender(s.sender),
			dbus.WithMatchInterface("org.kde.StatusNotifierItem"),
		)
		c := make(chan *dbus.Signal, 10)
		s.conn.Signal(c)
		for v := range c {
			switch v.Name {
			case signalItemNewAttentionIcon:
			case signalItemNewIcon:
			case signalItemNewOverlayIcon:
			case signalItemNewStatus:
			case signalItemNewTitle:
			case signalItemNewToolTip:
				s.updateProps()

				go func() {
					observer.ItemUpdated(s.sender)
				}()
				break
			default:
				logger.Printf("unhandled signal: name(%s) path(%s) sender(%s) body(%+v)\n", v.Name, v.Path, v.Sender, v.Body)
				break
			}
		}
	}()

	return
}

func (s *SystrayItem) Unregister() error {
	if s.conn == nil || s.conn.Connected() == false {
		return nil
	}

	s.conn.Close()

	return nil
}

func (s *SystrayItem) updateProps() {
	s.obj = s.conn.Object(s.sender, s.path)

	var response map[string]any
	s.obj.Call("org.freedesktop.DBus.Properties.GetAll", 0).Store(&response)

	s.item = &item{}
	s.item.updateProps(response)
}

func (s *SystrayItem) GetAttentionIconName() string {
	return s.item.attentionIconName
}

func (s *SystrayItem) GetAttentionIconpixmap() Pixmap {
	return s.item.attentionIconPixmap
}

func (s *SystrayItem) GetAttentionMovieName() string {
	return s.item.attentionMovieName
}

func (s *SystrayItem) GetCategory() itemCategory {
	return s.item.category
}

func (s *SystrayItem) GetIconName() string {
	return s.item.iconName
}

func (s *SystrayItem) GetIconPixmap() Pixmap {
	return s.item.iconPixmap
}

func (s *SystrayItem) GetId() string {
	return s.item.id
}

func (s *SystrayItem) GetOverlayIconName() string {
	return s.item.overlayIconName
}

func (s *SystrayItem) GetOverlayIconPixmap() Pixmap {
	return s.item.overlayIconPixmap
}

func (s *SystrayItem) GetStatus() itemStatus {
	return s.item.status
}

func (s *SystrayItem) GetTitle() string {
	return s.item.title
}

func (s *SystrayItem) GetToolTip() ItemTooltip {
	return s.item.tooltip
}
