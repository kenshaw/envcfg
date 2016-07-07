package envcfg

type Option func(*Envcfg)

// EnvVar sets the name of the environment variable to get the configuration data
// from.
func EnvVar(name string) Option {
	return func(ec *Envcfg) {
		ec.envVarName = name
	}
}

// ConfigFile is the file to read data from.
func ConfigFile(path string) Option {
	return func(ec *Envcfg) {
		ec.configFile = path
	}
}

// KeyFilter adds a key filter.
func KeyFilter(key string, f Filter) Option {
	return func(ec *Envcfg) {
		ec.filters[key] = f
	}
}
