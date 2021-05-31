package plugins

type Plugin interface {
	Init(env map[string]string) error
	Load(*Plugins) error
	Backend() interface{}
	Close() error
}
