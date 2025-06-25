package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name    string                      `yaml:"name"`
	Root    string                      `yaml:"root"`
	Windows map[string]WindowDefinition `yaml:"windows"`
}

type WindowDefinition struct {
	// Can be either a list of commands (shorthand) or a detailed config
	Cmds   []string `yaml:"-"` // used if it's a list
	Dir    string   `yaml:"dir,omitempty"`
	Pre    string   `yaml:"pre,omitempty"`
	Post   string   `yaml:"post,omitempty"`
	Layout string   `yaml:"layout,omitempty"`
	// Cmds    []string `yaml:"cmds,omitempty"`
	Panes   Panes `yaml:"panes,omitempty"`
	IsShort bool  `yaml:"-"` // true if this is just []string
}

type PaneDefinition struct {
	Dir     string   `yaml:"dir,omitempty"`
	Cmds    []string `yaml:"-"` // used if it's just a list
	Pre     string   `yaml:"pre,omitempty"`
	NCmds   []string `yaml:"cmds,omitempty"`
	Post    string   `yaml:"post,omitempty"`
	IsShort bool     `yaml:"-"` // true if just []string
}

type Panes struct {
	Cmds    []string `yaml:"-"` // used if it's just a list
	Named   map[string]PaneDefinition
	IsShort bool `yaml:"-"` // true if just []string
}

func (w *WindowDefinition) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal into []string
	var raw []string
	if err := value.Decode(&raw); err == nil {
		w.Cmds = raw
		w.IsShort = true
		return nil
	}

	// Try full struct
	type Alias WindowDefinition
	var tmp Alias
	if err := value.Decode(&tmp); err != nil {
		return err
	}
	*w = WindowDefinition(tmp)
	return nil
}

func (p *PaneDefinition) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal into []string
	var raw []string
	if err := value.Decode(&raw); err == nil {
		p.Cmds = raw
		p.IsShort = true
		return nil
	}

	// Try full struct
	type Alias PaneDefinition
	var tmp Alias
	if err := value.Decode(&tmp); err != nil {
		return err
	}
	*p = PaneDefinition(tmp)
	// fmt.Print(p)
	p.Cmds = p.NCmds
	p.NCmds = nil
	return nil
}

func (p *Panes) UnmarshalYAML(value *yaml.Node) error {
	var list []string
	if err := value.Decode(&list); err == nil {
		p.Cmds = list
		p.IsShort = true
		return nil
	}

	var mp map[string]PaneDefinition
	if err := value.Decode(&mp); err == nil {
		p.Named = mp
		p.IsShort = false
		return nil
	}

	return fmt.Errorf("invalid format for panes")
}

func (p Panes) MarshalJSON() ([]byte, error) {
	if p.IsShort {
		return json.Marshal(p.Cmds)
	}
	return json.Marshal(p.Named)
}
