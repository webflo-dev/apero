package services

type ServiceApi interface {
	Start() bool
	Stop()
	IsStarted() bool
}

type Service struct {
	started bool
}

type LoopBehavior bool

const (
	LoopBehaviorContinue LoopBehavior = true
	LoopBehaviorStop     LoopBehavior = false
)

type Looper func() LoopBehavior
type Deferer func()

func NewService() Service {
	return Service{
		started: false,
	}
}

func (s *Service) Start(deferer Deferer, loop Looper) bool {
	if s.started {
		return true
	}

	s.started = true

	go func() {
		if deferer != nil {
			defer deferer()
		}

		for {
			if s.started == false {
				return
			}
			if loop != nil && loop() == LoopBehaviorStop {
				return
			}
		}
	}()

	return true
}

func (s *Service) Stop() {
	s.started = false
}

func (s *Service) IsStarted() bool {
	return s.started
}
