package systray

type service struct{}

var _server SystrayServer
var _service = newService()

func newService() *service {
	return &service{}
}

func StartService() {
	StartNewServer(_service)
}

func GetSystrayItems() ItemList {
	return _server.GetItems()
}

func (s *service) NewItem(name string, item *SystrayItem) {
}

func (s *service) ItemRemoved(name string) {}

func (s *service) ItemUpdated(name string) {}
