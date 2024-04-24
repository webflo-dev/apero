package watcher

import (
	"context"
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

// Signal is a common interface for all signals.
type Signal interface {
	Name() string
	Interface() string
	Sender() string

	path() dbus.ObjectPath
	values() []interface{}
}

// Emit sends the given signal to the bus.
func Emit(conn *dbus.Conn, s Signal) error {
	return conn.Emit(s.path(), s.Interface()+"."+s.Name(), s.values()...)
}

// ErrUnknownSignal is returned by LookupSignal when a signal cannot be resolved.
var ErrUnknownSignal = errors.New("unknown signal")

// LookupSignal converts the given raw D-Bus signal with variable body
// into one with typed structured body or returns ErrUnknownSignal error.
func LookupSignal(signal *dbus.Signal) (Signal, error) {
	switch signal.Name {
	case InterfaceStatusNotifierWatcher + "." + "StatusNotifierItemRegistered":
		if len(signal.Body) < 1 {
			return nil, fmt.Errorf("signal has %v args rather than the expected 1", len(signal.Body))
		}
		v0, ok := signal.Body[0].(string)
		if !ok {
			return nil, fmt.Errorf("prop .V0 is %T, not string", signal.Body[0])
		}
		return &StatusNotifierWatcher_StatusNotifierItemRegisteredSignal{
			sender: signal.Sender,
			Path:   signal.Path,
			Body: &StatusNotifierWatcher_StatusNotifierItemRegisteredSignalBody{
				V0: v0,
			},
		}, nil
	case InterfaceStatusNotifierWatcher + "." + "StatusNotifierItemUnregistered":
		if len(signal.Body) < 1 {
			return nil, fmt.Errorf("signal has %v args rather than the expected 1", len(signal.Body))
		}
		v0, ok := signal.Body[0].(string)
		if !ok {
			return nil, fmt.Errorf("prop .V0 is %T, not string", signal.Body[0])
		}
		return &StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal{
			sender: signal.Sender,
			Path:   signal.Path,
			Body: &StatusNotifierWatcher_StatusNotifierItemUnregisteredSignalBody{
				V0: v0,
			},
		}, nil
	case InterfaceStatusNotifierWatcher + "." + "StatusNotifierHostRegistered":
		if len(signal.Body) < 0 {
			return nil, fmt.Errorf("signal has %v args rather than the expected 0", len(signal.Body))
		}
		return &StatusNotifierWatcher_StatusNotifierHostRegisteredSignal{
			sender: signal.Sender,
			Path:   signal.Path,
			Body:   &StatusNotifierWatcher_StatusNotifierHostRegisteredSignalBody{},
		}, nil
	case InterfaceStatusNotifierWatcher + "." + "StatusNotifierHostUnregistered":
		if len(signal.Body) < 0 {
			return nil, fmt.Errorf("signal has %v args rather than the expected 0", len(signal.Body))
		}
		return &StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal{
			sender: signal.Sender,
			Path:   signal.Path,
			Body:   &StatusNotifierWatcher_StatusNotifierHostUnregisteredSignalBody{},
		}, nil
	default:
		return nil, ErrUnknownSignal
	}
}

// AddMatchSignal registers a match rule for the given signal,
// opts are appended to the automatically generated signal's rules.
func AddMatchSignal(conn *dbus.Conn, s Signal, opts ...dbus.MatchOption) error {
	return conn.AddMatchSignal(append([]dbus.MatchOption{
		dbus.WithMatchInterface(s.Interface()),
		dbus.WithMatchMember(s.Name()),
	}, opts...)...)
}

// RemoveMatchSignal unregisters the previously registered subscription.
func RemoveMatchSignal(conn *dbus.Conn, s Signal, opts ...dbus.MatchOption) error {
	return conn.RemoveMatchSignal(append([]dbus.MatchOption{
		dbus.WithMatchInterface(s.Interface()),
		dbus.WithMatchMember(s.Name()),
	}, opts...)...)
}

// Interface name constants.
const (
	InterfaceStatusNotifierWatcher = "org.kde.StatusNotifierWatcher"
)

// StatusNotifierWatcherer is org.kde.StatusNotifierWatcher interface.
type StatusNotifierWatcherer interface {
	// RegisterStatusNotifierItem is org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem method.
	RegisterStatusNotifierItem(service string) (err *dbus.Error)
	// RegisterStatusNotifierHost is org.kde.StatusNotifierWatcher.RegisterStatusNotifierHost method.
	RegisterStatusNotifierHost(service string) (err *dbus.Error)
}

// ExportStatusNotifierWatcher exports the given object that implements org.kde.StatusNotifierWatcher on the bus.
func ExportStatusNotifierWatcher(conn *dbus.Conn, path dbus.ObjectPath, v StatusNotifierWatcherer) error {
	return conn.ExportSubtreeMethodTable(map[string]interface{}{
		"RegisterStatusNotifierItem": v.RegisterStatusNotifierItem,
		"RegisterStatusNotifierHost": v.RegisterStatusNotifierHost,
	}, path, InterfaceStatusNotifierWatcher)
}

// UnexportStatusNotifierWatcher unexports org.kde.StatusNotifierWatcher interface on the named path.
func UnexportStatusNotifierWatcher(conn *dbus.Conn, path dbus.ObjectPath) error {
	return conn.Export(nil, path, InterfaceStatusNotifierWatcher)
}

// UnimplementedStatusNotifierWatcher can be embedded to have forward compatible server implementations.
type UnimplementedStatusNotifierWatcher struct{}

func (*UnimplementedStatusNotifierWatcher) iface() string {
	return InterfaceStatusNotifierWatcher
}

func (*UnimplementedStatusNotifierWatcher) RegisterStatusNotifierItem(service string) (err *dbus.Error) {
	err = &dbus.ErrMsgUnknownMethod
	return
}

func (*UnimplementedStatusNotifierWatcher) RegisterStatusNotifierHost(service string) (err *dbus.Error) {
	err = &dbus.ErrMsgUnknownMethod
	return
}

// NewStatusNotifierWatcher creates and allocates org.kde.StatusNotifierWatcher.
func NewStatusNotifierWatcher(object dbus.BusObject) *StatusNotifierWatcher {
	return &StatusNotifierWatcher{object}
}

// StatusNotifierWatcher implements org.kde.StatusNotifierWatcher D-Bus interface.
type StatusNotifierWatcher struct {
	object dbus.BusObject
}

// RegisterStatusNotifierItem calls org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem method.
func (o *StatusNotifierWatcher) RegisterStatusNotifierItem(ctx context.Context, service string) (err error) {
	err = o.object.CallWithContext(ctx, InterfaceStatusNotifierWatcher+".RegisterStatusNotifierItem", 0, service).Store()
	return
}

// RegisterStatusNotifierHost calls org.kde.StatusNotifierWatcher.RegisterStatusNotifierHost method.
func (o *StatusNotifierWatcher) RegisterStatusNotifierHost(ctx context.Context, service string) (err error) {
	err = o.object.CallWithContext(ctx, InterfaceStatusNotifierWatcher+".RegisterStatusNotifierHost", 0, service).Store()
	return
}

// GetRegisteredStatusNotifierItems gets org.kde.StatusNotifierWatcher.RegisteredStatusNotifierItems property.
//
// Annotations:
//
//	@org.qtproject.QtDBus.QtTypeName.Out0 = QStringList
func (o *StatusNotifierWatcher) GetRegisteredStatusNotifierItems(ctx context.Context) (registeredStatusNotifierItems []string, err error) {
	err = o.object.CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, InterfaceStatusNotifierWatcher, "RegisteredStatusNotifierItems").Store(&registeredStatusNotifierItems)
	return
}

// GetIsStatusNotifierHostRegistered gets org.kde.StatusNotifierWatcher.IsStatusNotifierHostRegistered property.
func (o *StatusNotifierWatcher) GetIsStatusNotifierHostRegistered(ctx context.Context) (isStatusNotifierHostRegistered bool, err error) {
	err = o.object.CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, InterfaceStatusNotifierWatcher, "IsStatusNotifierHostRegistered").Store(&isStatusNotifierHostRegistered)
	return
}

// GetProtocolVersion gets org.kde.StatusNotifierWatcher.ProtocolVersion property.
func (o *StatusNotifierWatcher) GetProtocolVersion(ctx context.Context) (protocolVersion int32, err error) {
	err = o.object.CallWithContext(ctx, "org.freedesktop.DBus.Properties.Get", 0, InterfaceStatusNotifierWatcher, "ProtocolVersion").Store(&protocolVersion)
	return
}

// StatusNotifierWatcher_StatusNotifierItemRegisteredSignal represents org.kde.StatusNotifierWatcher.StatusNotifierItemRegistered signal.
type StatusNotifierWatcher_StatusNotifierItemRegisteredSignal struct {
	sender string
	Path   dbus.ObjectPath
	Body   *StatusNotifierWatcher_StatusNotifierItemRegisteredSignalBody
}

// Name returns the signal's name.
func (s *StatusNotifierWatcher_StatusNotifierItemRegisteredSignal) Name() string {
	return "StatusNotifierItemRegistered"
}

// Interface returns the signal's interface.
func (s *StatusNotifierWatcher_StatusNotifierItemRegisteredSignal) Interface() string {
	return InterfaceStatusNotifierWatcher
}

// Sender returns the signal's sender unique name.
func (s *StatusNotifierWatcher_StatusNotifierItemRegisteredSignal) Sender() string {
	return s.sender
}

func (s *StatusNotifierWatcher_StatusNotifierItemRegisteredSignal) path() dbus.ObjectPath {
	return s.Path
}

func (s *StatusNotifierWatcher_StatusNotifierItemRegisteredSignal) values() []interface{} {
	return []interface{}{s.Body.V0}
}

// StatusNotifierWatcher_StatusNotifierItemRegisteredSignalBody is body container.
type StatusNotifierWatcher_StatusNotifierItemRegisteredSignalBody struct {
	V0 string
}

// StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal represents org.kde.StatusNotifierWatcher.StatusNotifierItemUnregistered signal.
type StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal struct {
	sender string
	Path   dbus.ObjectPath
	Body   *StatusNotifierWatcher_StatusNotifierItemUnregisteredSignalBody
}

// Name returns the signal's name.
func (s *StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal) Name() string {
	return "StatusNotifierItemUnregistered"
}

// Interface returns the signal's interface.
func (s *StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal) Interface() string {
	return InterfaceStatusNotifierWatcher
}

// Sender returns the signal's sender unique name.
func (s *StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal) Sender() string {
	return s.sender
}

func (s *StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal) path() dbus.ObjectPath {
	return s.Path
}

func (s *StatusNotifierWatcher_StatusNotifierItemUnregisteredSignal) values() []interface{} {
	return []interface{}{s.Body.V0}
}

// StatusNotifierWatcher_StatusNotifierItemUnregisteredSignalBody is body container.
type StatusNotifierWatcher_StatusNotifierItemUnregisteredSignalBody struct {
	V0 string
}

// StatusNotifierWatcher_StatusNotifierHostRegisteredSignal represents org.kde.StatusNotifierWatcher.StatusNotifierHostRegistered signal.
type StatusNotifierWatcher_StatusNotifierHostRegisteredSignal struct {
	sender string
	Path   dbus.ObjectPath
	Body   *StatusNotifierWatcher_StatusNotifierHostRegisteredSignalBody
}

// Name returns the signal's name.
func (s *StatusNotifierWatcher_StatusNotifierHostRegisteredSignal) Name() string {
	return "StatusNotifierHostRegistered"
}

// Interface returns the signal's interface.
func (s *StatusNotifierWatcher_StatusNotifierHostRegisteredSignal) Interface() string {
	return InterfaceStatusNotifierWatcher
}

// Sender returns the signal's sender unique name.
func (s *StatusNotifierWatcher_StatusNotifierHostRegisteredSignal) Sender() string {
	return s.sender
}

func (s *StatusNotifierWatcher_StatusNotifierHostRegisteredSignal) path() dbus.ObjectPath {
	return s.Path
}

func (s *StatusNotifierWatcher_StatusNotifierHostRegisteredSignal) values() []interface{} {
	return []interface{}{}
}

// StatusNotifierWatcher_StatusNotifierHostRegisteredSignalBody is body container.
type StatusNotifierWatcher_StatusNotifierHostRegisteredSignalBody struct {
}

// StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal represents org.kde.StatusNotifierWatcher.StatusNotifierHostUnregistered signal.
type StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal struct {
	sender string
	Path   dbus.ObjectPath
	Body   *StatusNotifierWatcher_StatusNotifierHostUnregisteredSignalBody
}

// Name returns the signal's name.
func (s *StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal) Name() string {
	return "StatusNotifierHostUnregistered"
}

// Interface returns the signal's interface.
func (s *StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal) Interface() string {
	return InterfaceStatusNotifierWatcher
}

// Sender returns the signal's sender unique name.
func (s *StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal) Sender() string {
	return s.sender
}

func (s *StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal) path() dbus.ObjectPath {
	return s.Path
}

func (s *StatusNotifierWatcher_StatusNotifierHostUnregisteredSignal) values() []interface{} {
	return []interface{}{}
}

// StatusNotifierWatcher_StatusNotifierHostUnregisteredSignalBody is body container.
type StatusNotifierWatcher_StatusNotifierHostUnregisteredSignalBody struct {
}
