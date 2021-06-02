package plugins

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gdbu/queue"
	"github.com/gdbu/scribe"

	"github.com/hatchify/errors"
)

const (
	// ErrExpectedEndParen is returned when an ending parenthesis is missing
	ErrExpectedEndParen = errors.Error("expected ending parenthesis")
	// ErrInvalidDir is returned when a directory is empty
	ErrInvalidDir = errors.Error("invalid directory, cannot be empty")
	// ErrPluginKeyExists is returned when a plugin cannot be added because it already exists
	ErrPluginKeyExists = errors.Error("plugin cannot be added, key already exists")
	// ErrPluginNotLoaded is returned when a plugin namespace is provided that has not been loaded
	ErrPluginNotLoaded = errors.Error("plugin with that key has not been loaded")
	// ErrNotAddressable is returned when a non-addressable value is provided
	ErrNotAddressable = errors.Error("provided backend must be addressable")
)

var p = newPlugins()

// Register will register a plugin with a given key
func Register(key string, pi Plugin) error {
	return p.Register(key, pi)
}

// Get will retrieve a plugin with a given key
func Get(key string) (Plugin, error) {
	return p.Get(key)
}

// Loaded will return the plugins which have been loaded
func Loaded() map[string]Plugin {
	return p.Loaded()
}

// Backend will associated the backend of the requested key
func Backend(key string, backend interface{}) error {
	return p.Backend(key, backend)
}

func newPlugins() *Plugins {
	var p Plugins
	p.out = scribe.New("Plugins")
	p.pm = make(map[string]Plugin)
	return &p
}

// Plugins manages loaded plugins
type Plugins struct {
	mu  sync.RWMutex
	out *scribe.Scribe

	pm map[string]Plugin

	closed bool
}

// New will load a new plugin by plugin key
// The following formats are accepted as keys:
//	- path/to/file/plugin.so
//	- github.com/username/repository/pluginDir
func (p *Plugins) Register(key string, pi Plugin) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		err = errors.ErrIsClosed
		return
	}

	if _, ok := p.pm[key]; ok {
		return fmt.Errorf("plugin with the key of <%s> has already been loaded", key)
	}

	p.pm[key] = pi
	return
}

// Get will get a plugin by it's key
func (p *Plugins) Get(key string) (pi Plugin, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		err = errors.ErrIsClosed
		return
	}

	var ok bool
	if pi, ok = p.pm[key]; !ok {
		err = fmt.Errorf("plugin with key of <%s> has not been registered", key)
		return
	}

	return
}

func (p *Plugins) Loaded() (pm map[string]Plugin) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	pm = make(map[string]Plugin, len(p.pm))
	for key, val := range p.pm {
		pm[key] = val
	}

	return
}

// Backend will associated the backend of the requested key
func (p *Plugins) Backend(key string, backend interface{}) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errors.ErrIsClosed
	}

	var (
		pi Plugin
		ok bool
	)

	if pi, ok = p.pm[key]; !ok {
		err = fmt.Errorf("plugin with key of <%s> has not been registered", key)
		return
	}

	refVal := reflect.ValueOf(backend)
	elem := refVal.Elem()
	if !elem.CanSet() {
		return ErrNotAddressable
	}

	beVal := reflect.ValueOf(pi.Backend())

	switch {
	// Check to see if the types match exactly
	case elem.Type() == beVal.Type():
	// Check to see if the backend type implements the provided interface
	case beVal.Type().Implements(elem.Type()):

	default:
		// The provided value isn't an exact match, nor does it match the provided interface
		return fmt.Errorf("invalid type, expected %v and received %v", elem.Type(), beVal.Type())
	}

	elem.Set(beVal)
	return
}

// Test will test all of the plugins
func (p *Plugins) Test() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	//for _, pi := range p.pm {
	// TODO: Resolve test stuff here
	//if err = pi.test(); err != nil {
	//	return
	//}
	//}

	return errors.Error("testing has not yet been implemented")

}

// TestAsync will test all of the plugins asynchronously
func (p *Plugins) TestAsync(q *queue.Queue) (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	//var wg sync.WaitGroup
	//wg.Add(len(p.pm))
	//
	//var errs errors.ErrorList
	//for _, pi := range p.pm {
	//	q.New(func(pi Plugin) func() {
	//		return func() {
	//			defer wg.Done()
	//			// Fix test stuff here
	//		}
	//	}(pi))
	//}
	//
	//wg.Wait()
	//
	//return errs.Err()
	return errors.Error("testing has not yet been implemented")
}

// Close will close plugins
func (p *Plugins) Close() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errors.ErrIsClosed
	}

	var errs errors.ErrorList
	p.out.Notification("Closing plugins")
	for key, pi := range p.pm {
		if err = pi.Close(); err != nil {
			errs.Push(fmt.Errorf("error closing %s: %v", key, err))
			continue
		}

		p.out.Successf("Closed %s", key)
	}

	p.closed = true
	return errs.Err()
}
