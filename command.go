package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
)

func (r *Ranger) HandleCommand() error {
	cmd := r.EnterCommand(":")
	switch cmd {
	case TmuxUp, TmuxDown, TmuxLeft, TmuxRight:
		if err := TmuxNavigate(cmd); err != nil {
			return fmt.Errorf("tmux navigate: %s", err)
		}
	default:
		r.DisplayError("command not recognized: %s", cmd)
	}
	logger.Debugf("got user command: %s", cmd)
	return nil
}

func (r *Ranger) EnterCommand(pfx string) string {
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
			return ""
		case "<Enter>":
			return res
		default:
			res += e.ID
			r.statusBar.Text += e.ID
			ui.Render(r.statusBar)
		}
	}
}
