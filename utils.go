package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func GetTmuxOptionInt(option string) (int, error) {
	cmd := exec.Command("tmux", "show-options", "-g", option)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("failed to run tmux: %w", err)
	}

	// Output format: "base-index 1"
	parts := strings.Fields(out.String())
	if len(parts) != 2 {
		return 0, fmt.Errorf("unexpected output: %s", out.String())
	}

	val, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("not an integer: %w", err)
	}

	return val, nil
}

func GetWindowBaseIndex() int {
	res, err := GetTmuxOptionInt("base-index")
	if err != nil {
		return 0
	}
	return res
}

func GetPaneBaseIndex() int {
	res, err := GetTmuxOptionInt("pane-base-index")
	if err != nil {
		return 0
	}
	return res
}

func SysExec(command ...string) (string, error) {
	var cmd *exec.Cmd
	cmdLine := strings.Join(command, " ")

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdLine)
	} else {
		cmd = exec.Command("sh", "-c", cmdLine)
	}

	needsTTY := false
	if len(command) > 0 && strings.HasPrefix(command[0], "tmux") && strings.Contains(cmdLine, "attach") {
		needsTTY = true
	}

	if needsTTY {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		return "", err
	}

	var outputBuf bytes.Buffer
	cmd.Stdout = &outputBuf
	cmd.Stderr = &outputBuf
	err := cmd.Run()
	return strings.TrimSpace(outputBuf.String()), err
}

func AttachToSession(name string, window *string, pane *int) {
	if IsInsideTmuxEnv() {
		Log("Switching to session", name)
		// SysExec(fmt.Sprintf("tmux switch-client -t=%s", name))
		if window != nil && pane != nil {
			SysExec("tmux", "switch-client", "-t", fmt.Sprintf("%s:%s.%d", name, *window, *pane))
		} else if window != nil {
			SysExec("tmux", "switch-client", "-t", fmt.Sprintf("%s:%s", name, *window))
		} else {
			SysExec("tmux", "switch-client", "-t", name)
		}
		Log("Switched to session", name)
	} else {
		Log("Attaching to session", name)

		if window != nil && pane != nil {
			SysExec("tmux", "attach", "-t", fmt.Sprintf("%s:%s.%d", name, *window, *pane))
		} else if window != nil {
			SysExec("tmux", "attach", "-t", fmt.Sprintf("%s:%s", name, *window))
		} else {
			SysExec("tmux", "attach", "-t", name)
		}
	}
}

func HasTmuxSession(name string) bool {
	_, err := SysExec(fmt.Sprintf("tmux has-session -t=%s", name))
	return err == nil
}

func IsInsideTmuxEnv() bool {
	res, _ := SysExec("echo $TMUX")
	return res != ""
}

func ExecTmux(session string, windowIndex int, paneIndex int, cmd string) {
	SysExec(fmt.Sprintf("tmux send-keys -t %s:%d.%d \"%s\" Enter", session, windowIndex, paneIndex, cmd))
}

func Log(msg any, others ...any) {
	if !Args.DebugMode {
		return
	}
	fmt.Println(append([]any{msg}, others...)...)
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
