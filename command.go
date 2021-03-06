package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	ui "github.com/gizak/termui/v3"
)

type Command string

const (
	Touch  Command = "touch"
	Mkdir  Command = "mkdir"
	Delete Command = "delete"

	TmuxUp    Command = "TmuxNavigateUp"
	TmuxDown  Command = "TmuxNavigateDown"
	TmuxRight Command = "TmuxNavigateRight"
	TmuxLeft  Command = "TmuxNavigateLeft"
)

var reMetaKeystroke = regexp.MustCompile(`^<.+>$`)

func (r *Ranger) HandleCommand(text string) error {
	cmd, args := r.EnterCommand(":", text)
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
	case Delete:
		deletePath := r.pathToSelection()
		if err := os.Remove(deletePath); err != nil {
			return err
		}
		if r.SelectedLast() {
			r.Scroll(-1)
		} else {
			r.Scroll(1)
		}
		return r.ReloadDirs(r.baseName())
	default:
		return fmt.Errorf("command not recognized: %s", cmd)
	}
	return nil
}

func (r *Ranger) EnterCommand(pfx string, res string) (Command, []string) {
	r.statusBar.Text = pfx + res
	ui.Render(r.statusBar)
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
			res = res[:len(res)-1]
			r.statusBar.Text = pfx + res
			ui.Render(r.statusBar)
		default:
			// don't listen to other <...> events (eg <C-o>)
			if len(e.ID) > 1 {
				continue
			}
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
