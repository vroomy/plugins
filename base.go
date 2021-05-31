package plugins

var _ Plugin = &BasePlugin{}

type BasePlugin struct{}

func (b *BasePlugin) Init(env map[string]string) error {
	return nil
}

func (b *BasePlugin) Load(p *Plugins) error {
	return nil
}

func (b *BasePlugin) Backend() interface{} {
	return nil
}

func (b *BasePlugin) Close() error {
	return nil
}
