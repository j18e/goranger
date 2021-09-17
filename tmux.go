package main

import (
	"os"
	"os/exec"
)

const (
	TmuxUp    = "TmuxNavigateUp"
	TmuxDown  = "TmuxNavigateDown"
	TmuxRight = "TmuxNavigateRight"
	TmuxLeft  = "TmuxNavigateLeft"
)

func TmuxNavigate(cmd string) error {
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
