package main

type FConfig struct {
	Name    string
	Root    string
	Windows []FWindowDefinition
}

type FWindowDefinition struct {
	Name   string
	Dir    string
	Pre    string
	Post   string
	Layout string
	Panes  []FPaneDefinition
}

type FPaneDefinition struct {
	Dir  string
	Cmds []string
	Pre  string
	Post string
}

func FormatConfig(cfg Config) FConfig {
	return FConfig{Name: cfg.Name, Root: cfg.Root, Windows: formatWindows(cfg.Windows)}
}

func formatWindows(windows map[string]WindowDefinition) []FWindowDefinition {
	var formatted []FWindowDefinition
	for name, window := range windows {
		nPanes := formatPanes(window.Panes)
		if len(window.Cmds) > 0 {
			// create a pane with all cmds
			nPanes = append(nPanes, FPaneDefinition{Cmds: window.Cmds})
		}
		formatted = append(formatted, FWindowDefinition{Name: name, Dir: window.Dir, Pre: window.Pre, Post: window.Post, Layout: window.Layout, Panes: nPanes})
		// formatted[name] = FWindowDefinition{Dir: window.Dir, Pre: window.Pre, Post: window.Post, Layout: window.Layout, Panes: nPanes}
	}
	return formatted
}

func formatPanes(panes Panes) []FPaneDefinition {
	var formatted []FPaneDefinition
	for _, pane := range panes.Named {
		formatted = append(formatted, FPaneDefinition{Dir: pane.Dir, Cmds: pane.Cmds, Pre: pane.Pre, Post: pane.Post})
	}
	for _, rawPane := range panes.Cmds {
		formatted = append(formatted, FPaneDefinition{Cmds: []string{rawPane}})
	}
	return formatted
}
