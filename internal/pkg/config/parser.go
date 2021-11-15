package config

import (
	"io/ioutil"

	"olympos.io/encoding/edn"
)

type ConfigSet struct {
	Binary string
	Token  string
}

// Parse parses the configuration given at path.
func Parse(path string) (error, ConfigSet) {
	if len(path) == 0 {
		return nil, ConfigSet{}
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err, ConfigSet{}
	}
	var ret ConfigSet
	err = edn.Unmarshal(data, &ret)
	return nil, ret
}
