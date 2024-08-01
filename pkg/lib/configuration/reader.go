package configuration

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func Reader(path string) (*Configuration, error) {
	result := Configuration{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
