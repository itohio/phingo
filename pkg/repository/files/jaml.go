package repository

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// FIXME
func yaml2json(in []byte) ([]byte, error) {
	data := map[string]interface{}{}
	err := yaml.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

// FIXME
func json2yaml(in []byte) ([]byte, error) {
	data := map[string]interface{}{}

	err := json.Unmarshal(in, &data)
	if err != nil {
		return nil, err
	}
	return yaml.Marshal(data)
}
