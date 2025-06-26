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

func formatWindows(windows Windows) []FWindowDefinition {
	if windows.IsShort {
		return []FWindowDefinition{{Panes: []FPaneDefinition{{Cmds: windows.Cmds}}}}
	}
	var formatted []FWindowDefinition
	for _, window := range windows.Named {
		nPanes := formatPanes(window.Val.Panes)
		if len(window.Val.Cmds) > 0 {
			// create a pane with all cmds
			nPanes = append(nPanes, FPaneDefinition{Cmds: window.Val.Cmds})
		}
		formatted = append(formatted, FWindowDefinition{Name: window.Name, Dir: window.Val.Dir, Pre: window.Val.Pre, Post: window.Val.Post, Layout: window.Val.Layout, Panes: nPanes})
		// formatted[name] = FWindowDefinition{Dir: window.Dir, Pre: window.Pre, Post: window.Post, Layout: window.Layout, Panes: nPanes}
	}
	return formatted
}

func formatPanes(panes Panes) []FPaneDefinition {

	var formatted []FPaneDefinition
	for _, pane := range panes.Named {
		formatted = append(formatted, FPaneDefinition{Dir: pane.Val.Dir, Cmds: pane.Val.Cmds, Pre: pane.Val.Pre, Post: pane.Val.Post})
	}
	for _, rawPane := range panes.Cmds {
		formatted = append(formatted, FPaneDefinition{Cmds: []string{rawPane}})
	}
	return formatted
}
