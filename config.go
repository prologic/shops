package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
type Funcs map[string]string
type Items []Item

type Config struct {
	Version string `yaml:"version"`

	envList []*Env
	envMap  map[string]*Env
	Env     yaml.MapSlice `yaml:"env"`

	Files Files `yaml:"files"`
	Funcs Funcs `yaml:"funcs"`
	Items Items `yaml:"items"`
}

func (conf Config) SetEnvVars(env []string) {
	for _, e := range env {
		tokens := strings.Split(e, "=")

		var key, value string

		if len(tokens) == 2 {
			key, value = tokens[0], tokens[1]
		} else {
			key = tokens[0]
			value = os.Getenv(key)
		}

		conf.envMap[key].Value = value
	}
}

func (conf Config) Context(cmd string) Context {
	ctx := Context{
		Env:     conf.envList,
		Funcs:   conf.Funcs,
		Command: cmd,
	}

	return ctx
}

func readConfig(fn string) (conf Config, err error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		log.WithError(err).Error("error reading config file")
		err = fmt.Errorf("error reading config file %s: %w", fn, err)
		return
	}

	if err = yaml.Unmarshal([]byte(data), &conf); err != nil {
		log.WithError(err).Error("error parsing config file")
		err = fmt.Errorf("error parsing config file %s: %w", fn, err)
		return
	}

	conf.envMap = make(map[string]*Env)

	for _, item := range conf.Env {
		key, value := item.Key.(string), item.Value.(string)
		conf.envMap[key] = &Env{Key: key, Value: value}
		conf.envList = append(conf.envList, conf.envMap[key])
	}

	return
}
