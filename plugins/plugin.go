package plugins

type PluginInterface interface {
		Register()
		PluginName() string
}

var (
		_instances = map[string]PluginInterface{}
)

// 插件
func Plugin(name string, plugins ...PluginInterface) PluginInterface {
		if len(plugins) == 0 {
				return _instances[name]
		}
		_instances[name] = plugins[0]
		return plugins[0]
}
