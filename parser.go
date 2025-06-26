package main

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type OrderedMap[V interface{}] struct {
	Map   map[string]V
	Order []string
}

type Named[V interface{}] struct {
	Name string
	Val  V
}

type Config struct {
	Name    string  `yaml:"name"`
	Root    string  `yaml:"root"`
	Windows Windows `yaml:"windows"`
}

type Windows struct {
	Cmds    []string `yaml:"-"` // used if it's a list
	Named   []Named[WindowDefinition]
	IsShort bool `yaml:"-"` // true if just []string
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
	Named   []Named[PaneDefinition]
	IsShort bool `yaml:"-"` // true if just []string
}

func (w *WindowDefinition) UnmarshalYAML(value *yaml.Node) error {
	var list []string
	if err := value.Decode(&list); err == nil {
		w.Cmds = list
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
	var list []string
	if err := value.Decode(&list); err == nil {
		p.Cmds = list
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

	om, err := getOrderedMap[PaneDefinition](value)
	if err != nil {
		return err
	}
	ol, err := getOrderedList(om)
	if err != nil {
		return err
	}
	p.Named = *ol
	return nil
}

func (w *Windows) UnmarshalYAML(value *yaml.Node) error {
	var list []string
	if err := value.Decode(&list); err == nil {
		w.Cmds = list
		w.IsShort = true
		return nil
	}

	om, err := getOrderedMap[WindowDefinition](value)
	if err != nil {
		return err
	}
	ol, err := getOrderedList(om)
	if err != nil {
		return err
	}
	w.Named = *ol
	return nil
}

func (p Panes) MarshalJSON() ([]byte, error) {
	if p.IsShort {
		return json.Marshal(p.Cmds)
	}
	return json.Marshal(p.Named)
}

func getOrderedMap[V interface{}](node *yaml.Node) (om *OrderedMap[V], err error) {
	content := node.Content
	end := len(content)
	count := end / 2

	om = &OrderedMap[V]{
		Map:   make(map[string]V, count),
		Order: make([]string, 0, count),
	}

	for pos := 0; pos < end; pos += 2 {
		keyNode := content[pos]
		valueNode := content[pos+1]

		if keyNode.Tag != "!!str" {
			err = fmt.Errorf("expected a string key but got %s on line %d", keyNode.Tag, keyNode.Line)
			return
		}

		var k string
		if err = keyNode.Decode(&k); err != nil {
			return
		}

		var v V
		if err = valueNode.Decode(&v); err != nil {
			return
		}

		om.Map[k] = v
		om.Order = append(om.Order, k)
	}

	return
}

func getOrderedList[V interface{}](orderedMap *OrderedMap[V]) (ol *[]Named[V], err error) {
	ol = &[]Named[V]{}

	for _, k := range orderedMap.Order {
		v, ok := orderedMap.Map[k]
		if !ok {
			err = fmt.Errorf("key %s not found in map", k)
			return
		}
		*ol = append(*ol, Named[V]{Name: k, Val: v})
	}
	return ol, nil
}
