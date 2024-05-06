package hyprland

type Workspace struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Monitor         string `json:"monitor"`
	MonitorId       string `json:"monitorID"`
	Windows         int    `json:"windows"`
	HasFullscreen   bool   `json:"hasfullscreen"`
	LastWindow      string `json:"lastwindow"`
	LastWindowTitle string `json:"lastwindowtitle"`
}

type ClientWorkspace struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Client struct {
	Address        string          `json:"address"`
	Mapped         bool            `json:"mapped"`
	Hidden         bool            `json:"hidden"`
	At             []int           `json:"at"`
	Size           []int           `json:"size"`
	Workspace      ClientWorkspace `json:"workspace"`
	Floating       bool            `json:"floating"`
	Monitor        int             `json:"monitor"`
	Class          string          `json:"class"`
	Title          string          `json:"title"`
	InitialClass   string          `json:"initialClass"`
	InitialTitle   string          `json:"initialTitle"`
	Pid            int             `json:"pid"`
	XWayland       bool            `json:"xwayland"`
	Pinned         bool            `json:"pinned"`
	Fullscreen     bool            `json:"fullscreen"`
	FullscreenMode int             `json:"fullscreenMode"`
	FakeFullscreen bool            `json:"fakeFullscreen"`
	Grouped        []string        `json:"grouped"`
	Swallowing     string          `json:"swallowing"`
	FocusHistoryId int             `json:"focusHistoryID"`
}

type Bind struct {
	Locked       bool   `json:"locked"`
	Mouse        bool   `json:"mouse"`
	Release      bool   `json:"release"`
	Repeat       bool   `json:"repeat"`
	NonConsuming bool   `json:"non_consuming"`
	ModMask      int    `json:"modmask"`
	Submap       string `json:"submap"`
	Key          string `json:"key"`
	Keycode      int    `json:"keycode"`
	CatchAll     bool   `json:"catch_all"`
	Dispatcher   string `json:"dispatcher"`
	Arg          string `json:"arg"`
}

type CursorPos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Devices struct {
	Mice      []DeviceMouse    `json:"mice"`
	Keyboards []DeviceKeyboard `json:"keyboards"`
	Tablets   []DeviceTablet   `json:"tablets"`
	Touch     []DeviceTouch    `json:"touch"`
	Switches  []DeviceSwitch   `json:"switches"`
}

type DeviceKeyboard struct {
	Address      string `json:"address"`
	Name         string `json:"name"`
	Rules        string `json:"rules"`
	Model        string `json:"model"`
	Layout       string `json:"layout"`
	Variant      string `json:"variant"`
	Options      string `json:"options"`
	ActiveKeymap string `json:"active_keymap"`
	Main         bool   `json:"main"`
}

type DeviceMouse struct {
	Address      string  `json:"address"`
	Name         string  `json:"name"`
	DefaultSpeed float32 `json:"defaultSpeed"`
}

type DeviceTablet struct{}

type DeviceTouch struct{}

type DeviceSwitch struct{}

type Instance struct {
	Instance string `json:"instance"`
	Time     int    `json:"time"`
	Pid      int    `json:"pid"`
	WlSocket string `json:"wl_socket"`
}

type LayerType string

const (
	LayerBackground LayerType = "0"
	LayerOverlay    LayerType = "3"
	LayerTop        LayerType = "2"
	LayerBottom     LayerType = "1"
)

type Layer struct {
	Address     string //"address": "0x328b670",
	X           int    // "x": 0,
	Y           int    //"y": 0,
	Width       int    //"w": 3840,
	Height      int    //"h": 2160,
	Namespace   string //"namespace": "wallpaper"
	MonitorName string //"monitorName": "DP-1",
	Layer       LayerType
}

type MonitorWorkspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Monitor struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Make             string           `json:"make"`
	Model            string           `json:"model"`
	Serial           string           `json:"serial"`
	Width            int              `json:"width"`
	Height           int              `json:"height"`
	RefreshRate      int              `json:"refreshRate"`
	X                int              `json:"x"`
	Y                int              `json:"y"`
	ActiveWorkspace  MonitorWorkspace `json:"activeWorkspace"`
	SpecialWorkspace MonitorWorkspace `json:"specialWorkspace"`
	Reserved         []int            `json:"reserved"`
	Scale            float32          `json:"scale"`
	Transform        string           `json:"transform"`
	Focused          bool             `json:"focused"`
	DpmsStatus       bool             `json:"dpmsStatus"`
	VRR              bool             `json:"vrr"`
	ActivelyTearing  bool             `json:"activelyTearing"`
	CurrentFormat    string           `json:"currentFormat"`
	AvailableModes   []string         `json:"availableModes"`
}
