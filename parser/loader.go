package parser

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadYAML(filename *string, v interface{}) error {
	raw, err := ioutil.ReadFile(*filename)
	if err != nil {
		return err
	} else {
		err = yaml.Unmarshal(raw, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadJSON(filename *string, v interface{}) error {
	raw, err := ioutil.ReadFile(*filename)
	if err != nil {
		return err
	} else {
		err = json.Unmarshal(raw, v)
		if err != nil {
			return err
		}
	}
	return nil
}
