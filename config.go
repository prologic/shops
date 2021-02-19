package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Item struct {
	Name   string `yaml:"name"`
	Check  string `yaml:"check"`
	Action string `yaml"action"`
}

type Items []Item

type Config struct {
	Version string `yaml:"version"`

	Items Items `yaml"items"`
}

func readConfig(fn string) (config Config, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.WithError(err).Error("error reading config file")
		err = fmt.Errorf("error reading config file %s: %w", fn, err)
		return
	}

	if err = yaml.Unmarshal([]byte(data), &config); err != nil {
		log.WithError(err).Error("error parsing config file")
		err = fmt.Errorf("error parsing config file %s: %w", fn, err)
		return
	}

	return
}
