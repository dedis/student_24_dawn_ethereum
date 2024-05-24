// Copyright EPFL DEDIS

package f3b

import (
	"os"
	"encoding/json"
)

type Params struct {
	Protocol string `toml:protocol`
}

type FullParams struct {
	Params
	SmcPath string
}

var cached *FullParams

func ReadParams() (*FullParams, error) {
	if cached != nil {
		return cached, nil
	}
	path, ok := os.LookupEnv("F3B_PARAMS")
	if !ok {
		path = "params.json"
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cached = new(FullParams)
	err = decoder.Decode(cached)
	if err != nil {
		return nil, err
	}
	return cached, nil
}
