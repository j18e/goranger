package main

import (
	"os"
	"os/exec"
)

func TmuxNavigate(cmd Command) error {
	var arg string
	switch cmd {
	case TmuxUp:
		arg = "-U"
	case TmuxDown:
		arg = "-D"
	case TmuxLeft:
		arg = "-L"
	case TmuxRight:
		arg = "-R"
	}
	return exec.Command("tmux", "select-pane", arg, "-t", os.Getenv("TMUX_PANE")).Start()
}
