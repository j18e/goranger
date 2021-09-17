package main

import (
	"errors"

	ui "github.com/gizak/termui/v3"
)

func (r *Ranger) HandleEvent(e ui.Event) error {
	switch e.ID {
	case ":":
		return r.HandleCommand(r.EnterCommand(":"))
	case "/":
		// TODO handle search

	case "q", "<C-c>":
		return ErrExit
	case "l", "<Right>":
		if err := r.LevelDown(); err != nil {
			if errors.Is(err, ErrExit) {
				return err
			}
			logger.Error(err)
		}
	case "h", "<Left>":
		r.LevelUp()
	case "j", "<Down>":
		r.Scroll(1)
	case "k", "<Up>":
		r.Scroll(-1)
	case "J":
		r.Scroll(10)
	case "K":
		r.Scroll(-10)
	case "g":
		if r.prevKey == "g" {
			r.ScrollTop()
		}
	case "G", "<End>":
		r.ScrollBottom()
	case "<C-h>", "<C-<Backspace>>":
		if r.showHidden {
			r.showHidden = false
		} else {
			r.showHidden = true
		}
		r.ReloadDirs()
	}
	if err := r.updatePreview(); err != nil {
		logger.Error(err)
	}
	r.DisplayFile()

	if r.prevKey == "g" {
		r.prevKey = ""
	} else {
		r.prevKey = e.ID
	}
	return nil
}
