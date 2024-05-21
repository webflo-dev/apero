package systray

var _service = newService()

func StartService() {
	_service.server.start()
}

func GetSystrayItems() ItemList {
	return _service.server.items
}
