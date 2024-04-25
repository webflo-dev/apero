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
