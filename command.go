package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type Command string

const (
	Touch Command = "touch"
	Mkdir Command = "mkdir"

	TmuxUp    Command = "TmuxNavigateUp"
	TmuxDown  Command = "TmuxNavigateDown"
	TmuxRight Command = "TmuxNavigateRight"
	TmuxLeft  Command = "TmuxNavigateLeft"
)

func (r *Ranger) HandleCommand() error {
	cmd, args := r.EnterCommand(":")
	switch cmd {
	case "":
		r.UpdateStatus()
	case TmuxUp, TmuxDown, TmuxLeft, TmuxRight:
		return TmuxNavigate(cmd)
	case Mkdir, Touch:
		if len(args) != 1 {
			return fmt.Errorf("%s: expected 1 arg, got %d", cmd, len(args))
		}
		path := filepath.Join(r.path, args[0])
		if cmd == Mkdir {
			if err := os.Mkdir(path, 0755); err != nil {
				return err
			}
		} else {
			if err := ioutil.WriteFile(path, nil, 0644); err != nil {
				return err
			}
		}
		return r.ReloadDirs("")
	default:
		return fmt.Errorf("command not recognized: %s", cmd)
	}
	return nil
}

func (r *Ranger) EnterCommand(pfx string) (Command, []string) {
	r.statusBar.Text = pfx
	ui.Render(r.statusBar)
	var res string
	for {
		e := <-r.events
		if e.Type != ui.KeyboardEvent {
			continue
		}
		switch e.ID {
		case "<Escape>", "<C-c>":
			return "", nil
		case "<Enter>":
			res = strings.TrimSpace(res)
			split := strings.Split(res, " ")
			if len(split) > 1 {
				return Command(split[0]), split[1:]
			}
			return Command(split[0]), nil
		case "<Space>":
			res += " "
			r.statusBar.Text = pfx + res
			ui.Render(r.statusBar)
		case "<Backspace>":
			res = res[:len(res)-2]
			r.statusBar.Text = pfx + res
			ui.Render(r.statusBar)
		default:
			res += e.ID
			r.statusBar.Text = pfx + res
			ui.Render(r.statusBar)
		}
	}
}

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
