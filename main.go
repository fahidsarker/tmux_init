package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func getConfig() Config {
	yamlBytes, err := os.ReadFile(Args.ConfigFile)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(yamlBytes, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func main() {
	session := FormatConfig(getConfig())
	paneIdx := GetPaneBaseIndex()

	if HasTmuxSession(session.Name) {
		Log("Session already exists")
		AttachToSession(session.Name, &session.Windows[0].Name, &paneIdx)
	} else {
		Log("Creating session")
		SysExec(fmt.Sprintf("tmux new-session -d -s %s -n %s", session.Name, session.Name))
		windowBaseIdx := GetWindowBaseIndex()
		Log("Window base index", windowBaseIdx)
		for index, window := range session.Windows {
			BuildWindow(session.Name, session.Root, index+windowBaseIdx, window, index != 0)
		}
		AttachToSession(session.Name, &session.Windows[0].Name, &paneIdx)
	}
}

func BuildWindow(sessionName string, rootDir string, index int, window FWindowDefinition, isNew bool) {
	Log("Building window", window.Name, "index", index)
	basePaneIdx := GetPaneBaseIndex()
	Log("Pane base index", basePaneIdx)
	if window.Pre != "" {
		SysExec(window.Pre)
	}

	if isNew {
		SysExec(fmt.Sprintf("tmux new-window -t %s:%d -n %s", sessionName, index, window.Name))
	} else {
		SysExec(fmt.Sprintf("tmux rename-window -t %s:%d %s", sessionName, index, window.Name))
	}

	for paneIdx, pane := range window.Panes {
		BuildPane(sessionName, filepath.Join(rootDir, window.Dir), index, pane, paneIdx+basePaneIdx)
	}

	if window.Layout != "" {
		SysExec(fmt.Sprintf("tmux select-layout -t %s:%d %s", sessionName, index, window.Layout))
	}

	if window.Post != "" {
		SysExec(window.Post)
	}
	Log("Done building window", window.Name)
	Log(".........")
}

func BuildPane(sessionName string, rootDir string, windowIndex int, pane FPaneDefinition, paneIndex int) {
	Log("Building pane", sessionName, windowIndex, paneIndex)
	if pane.Pre != "" {
		SysExec(pane.Pre)
	}

	if paneIndex != 1 {
		SysExec(fmt.Sprintf("tmux splitw -h -t %s:%d", sessionName, windowIndex))
	}
	paneDir := filepath.Join(rootDir, pane.Dir)
	if paneDir != "" {
		ExecTmux(sessionName, windowIndex, paneIndex, fmt.Sprintf("cd %s", paneDir))
		ExecTmux(sessionName, windowIndex, paneIndex, "clear")
	}

	for _, cmd := range pane.Cmds {
		ExecTmux(sessionName, windowIndex, paneIndex, cmd)
	}

	if pane.Post != "" {
		SysExec(pane.Post)
	}
}
