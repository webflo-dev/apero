package systray

type systrayService struct {
	server *server
}

func newService() *systrayService {
	return &systrayService{
		server: newServer(),
	}
}
