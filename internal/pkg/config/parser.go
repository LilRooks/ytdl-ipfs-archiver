package config

type ConfigSet struct {
	Binary string
}

// Parse parses the configuration given at path.
func Parse(path string) (error, ConfigSet) {
	return nil, ConfigSet{Binary: path}
}
