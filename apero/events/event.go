package events

import (
	"context"
	"errors"
	"sync"
)

//
// in-memory event store
//

type ID uint32

var __id ID = 0

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

	__id++
	x := __id
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
	RegisterHandler(id string, f Handler[Payload]) error
	UnregisterHandler(id string)
	Publish(payload Payload) error
	PublishWithContext(ctx context.Context, payload Payload) error
}

//
// Event system
//

var ErrDuplicateID = errors.New("Duplicate handler ID")

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
	handlers map[string]Handler[Payload]
	mu       sync.RWMutex
}

func (e *event[Payload]) subscribe(payload Payload) {
	for _, handler := range e.handlers {
		go handler.Handle(payload)
	}
}

// Register registers an event handler
func (e *event[Payload]) RegisterHandler(id string, f Handler[Payload]) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.handlers == nil {
		e.handlers = map[string]Handler[Payload]{}
	}

	if _, ok := e.handlers[id]; ok {
		return ErrDuplicateID
	}

	e.handlers[id] = f

	return nil
}

// Register registers an event handler
func (c *event[Payload]) UnregisterHandler(id string) {
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
