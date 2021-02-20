package main

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type File struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Mode   string `yaml:"mode"`
}

func (f File) String() string {
	return fmt.Sprintf("%s:%s", f.Source, f.Target)
}

type Item struct {
	Name   string `yaml:"name"`
	Check  string `yaml:"check"`
	Action string `yaml:"action"`
}

func (i Item) String() string {
	return i.Name
}

type Files []File
type Items []Item

type Config struct {
	Version string `yaml:"version"`

	Files Files `yaml:"files"`
	Items Items `yaml:"items"`
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
