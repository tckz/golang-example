package util

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func LoadYaml(fn string) (map[interface{}]interface{}, error) {

	buf, err := func() ([]byte, error) {
		if fn == "-" {
			return ioutil.ReadAll(os.Stdin)
		} else {
			return ioutil.ReadFile(fn)
		}
	}()
	if err != nil {
		return nil, err
	}

	y := make(map[interface{}]interface{})
	err = yaml.Unmarshal(buf, &y)
	if err != nil {
		return nil, err
	}

	return y, nil
}
