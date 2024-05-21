package dbus

import (
	"log"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

type Bus struct {
	conn   *dbus.Conn
	logger *log.Logger
}

func ConnectToSessionBus(logger *log.Logger) (b *Bus, err error) {
	b = &Bus{
		logger: logger,
	}

	if b.conn, err = dbus.ConnectSessionBus(); err != nil {
		b.logger.Println("failed to connect to session bus.", err)
	}

	return
}

func (b *Bus) Conn() *dbus.Conn {
	return b.conn
}

func (b *Bus) Close() {
	if b.conn != nil && b.conn.Connected() {
		b.conn.Close()
	}
}

func (b *Bus) ExportAll(v interface{}, path dbus.ObjectPath, iface string) bool {
	err := b.conn.ExportAll(v, path, iface)
	if err != nil {
		b.logger.Println("failed to export to dbus.", err)
		return false
	}

	return true
}

func (b *Bus) ExportProps(path string, iface string, props map[string]*prop.Prop) (p *prop.Properties, ok bool) {
	propsSpec := map[string]map[string]*prop.Prop{}
	propsSpec[iface] = props

	p, err := prop.Export(b.conn, dbus.ObjectPath(path), propsSpec)
	if err != nil {
		b.logger.Println("failed to export properties to dbus.", err)
		return nil, false
	}

	return p, true
}

func (b *Bus) RequestName(iface string) bool {
	reply, err := b.conn.RequestName(iface, dbus.NameFlagDoNotQueue)
	if err != nil {
		b.logger.Println("failed to request name on session dbus.", err)
		return false
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		b.logger.Println("name already taken.", err)
		return false
	}

	return true
}

func (b *Bus) ExportIntrospectable(path string, data introspect.Interface) bool {
	node := &introspect.Node{
		Name: path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			data,
		},
	}

	err := b.conn.Export(introspect.NewIntrospectable(node), dbus.ObjectPath(path), "org.freedesktop.DBus.Introspectable")
	if err != nil {
		b.logger.Println("failed to export introspectable data", err)
		return false
	}

	return true
}

func (b *Bus) EmitSignal(path string, signal string, values ...interface{}) bool {
	err := b.conn.Emit(dbus.ObjectPath(path), signal, values...)
	if err != nil {
		b.logger.Printf("Failed to emit signal <%s> %v\n", signal, err)
		return false
	}
	b.logger.Printf("Emitted signal <%s> %v\n", signal, values)
	return true
}

type WatchBehavior bool
type SignalWatcher func(signal *dbus.Signal) WatchBehavior

const (
	WatchBehaviorContinue WatchBehavior = true
	WatchBehaviorStop     WatchBehavior = false
)

func (b *Bus) WatchSignal(callback SignalWatcher, options ...dbus.MatchOption) bool {
	if !b.addMatchSignal(options...) {
		return false
	}

	return b.watchSignal(callback)
}

func (b *Bus) WatchSignals(callback SignalWatcher, options map[string][]dbus.MatchOption) bool {
	for _, opt := range options {
		if !b.addMatchSignal(opt...) {
			return false
		}
	}

	return b.watchSignal(callback)
}

func (b *Bus) addMatchSignal(options ...dbus.MatchOption) bool {
	err := b.conn.AddMatchSignal(options...)
	if err != nil {
		b.logger.Println("failed to watch signal", err)
		return false
	}
	return true
}

func (b *Bus) watchSignal(callback SignalWatcher) bool {
	c := make(chan *dbus.Signal, 10)
	b.conn.Signal(c)

	go func() {
		for v := range c {
			if callback(v) == WatchBehaviorStop {
				break
			}
		}
	}()

	return true
}

func (b *Bus) Call(dest string, path dbus.ObjectPath, method string, store interface{}, args ...interface{}) bool {
	obj := b.conn.Object(dest, path)
	return b.CallWithObj(obj, method, store, args...)
}

func (b *Bus) CallWithObj(obj dbus.BusObject, method string, store interface{}, args ...interface{}) bool {
	call := obj.Call(method, 0, args...)
	if call.Err != nil {
		b.logger.Println("failed to call method", method, call.Err)
		return false
	}

	if err := call.Store(store); err != nil {
		b.logger.Println("failed to store result", err)
		return false
	}

	return true
}
