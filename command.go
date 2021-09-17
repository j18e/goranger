package main

import (
	"fmt"
)

func (r *Ranger) HandleCommand(cmd string) error {
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
