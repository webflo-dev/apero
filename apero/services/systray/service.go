package systray

type service struct{}

var _server SystrayServer
var _service = newService()

func newService() *service {
	return &service{}
}

func StartService() {
	StartNewServer()
}
