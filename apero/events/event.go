package events

import (
	"context"
	"sync"
)

//
// in-memory event store
//

type ID uint32

var __storeIds ID = 0

func newStore[Payload any]() *memoryStore[Payload] {
	return &memoryStore[Payload]{
		callbacks: make(map[ID]func(Payload)),
	}
}

func New[Payload any]() Event[Payload] {
	return NewWithContext[Payload](context.Background())
}

func NewWithContext[Payload any](ctx context.Context) Event[Payload] {
	store := newStore[Payload]()
	return NewEvent(ctx, store)
}

type memoryStore[Payload any] struct {
	mu        sync.RWMutex
	callbacks map[ID]func(Payload)
}

func (m *memoryStore[Payload]) Publish(payload Payload) error {
	return m.PublishWithContext(context.Background(), payload)
}

func (m *memoryStore[Payload]) PublishWithContext(ctx context.Context, payload Payload) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, c := range m.callbacks {
		go c(payload)
	}

	return nil
}

func (m *memoryStore[Payload]) Subscribe(f func(Payload)) error {
	return m.SubscribeWithContext(context.Background(), f)
}

func (m *memoryStore[Payload]) SubscribeWithContext(ctx context.Context, f func(Payload)) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	__storeIds++
	x := __storeIds
	m.callbacks[x] = f

	// close and cleanup the channel
	go func() {
		<-ctx.Done()

		m.mu.Lock()
		defer m.mu.Unlock()

		delete(m.callbacks, x)
	}()

	return nil
}

//
// Handler
//

type Handler[Payload any] interface {
	Handle(payload Payload)
}

type HandlerFunc[Payload any] func(payload Payload)

func (e HandlerFunc[Payload]) Handle(payload Payload) {
	e(payload)
}

type Event[Payload any] interface {
	RegisterHandler(f Handler[Payload]) (ID, error)
	UnregisterHandler(id ID)
	Publish(payload Payload) error
	PublishWithContext(ctx context.Context, payload Payload) error
}

//
// Event system
//

type Store[Payload any] interface {
	Publish(payload Payload) error
	PublishWithContext(ctx context.Context, payload Payload) error

	// Free any resources when the context is done
	Subscribe(f func(Payload)) error
	SubscribeWithContext(ctx context.Context, f func(Payload)) error
}

func NewEvent[Payload any](ctx context.Context, store Store[Payload]) Event[Payload] {
	e := &event[Payload]{
		store: store,
	}

	// Will exit when the subscription channel closes
	if err := store.SubscribeWithContext(ctx, e.subscribe); err != nil {
		return nil
	}

	return e
}

type event[Payload any] struct {
	store    Store[Payload]
	handlers map[ID]Handler[Payload]
	mu       sync.RWMutex
}

var __ids ID = 0

func (e *event[Payload]) subscribe(payload Payload) {
	for _, handler := range e.handlers {
		go handler.Handle(payload)
	}
}

// Register registers an event handler
func (e *event[Payload]) RegisterHandler(f Handler[Payload]) (ID, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.handlers == nil {
		e.handlers = map[ID]Handler[Payload]{}
	}

	__ids++
	e.handlers[__ids] = f

	return __ids, nil
}

// Register registers an event handler
func (c *event[Payload]) UnregisterHandler(id ID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.handlers, id)
}

func (e *event[Payload]) Publish(payload Payload) error {
	return e.PublishWithContext(context.Background(), payload)
}

func (e *event[Payload]) PublishWithContext(ctx context.Context, payload Payload) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.store.PublishWithContext(ctx, payload)
}
