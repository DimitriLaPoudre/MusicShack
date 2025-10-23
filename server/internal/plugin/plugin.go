package plugin

type Plugin interface {
	Test() string
}

var registry = make(map[string]Plugin)

func Register(name string, p Plugin) {
	registry[name] = p
}

func Get(name string) (Plugin, bool) {
	p, ok := registry[name]
	return p, ok
}
