package main

import (
	"fmt"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type Command string

const (
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
	case Mkdir:
		if len(args) != 1 {
			return fmt.Errorf("mkdir: expected 1 arg, got %d", len(args))
		}
		if err := os.Mkdir(args[0], 0755); err != nil {
			return err
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
	res := pfx
	for {
		e := <-r.events
		if e.Type != ui.KeyboardEvent {
			continue
		}
		switch e.ID {
		case "<Escape>", "<C-c>":
			return "", nil
		case "<Enter>":
			res = strings.TrimSpace(strings.TrimLeft(res, pfx))
			split := strings.Split(res, " ")
			if len(split) > 1 {
				return Command(split[0]), split[1:]
			}
			return Command(split[0]), nil
		case "<Space>":
			res += " "
			r.statusBar.Text = res
			ui.Render(r.statusBar)
		case "<Backspace>":
			res = res[:len(res)-2]
			r.statusBar.Text = res
			ui.Render(r.statusBar)
		default:
			res += e.ID
			r.statusBar.Text = res
			ui.Render(r.statusBar)
		}
	}
}
