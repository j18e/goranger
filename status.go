package main

import (
	"fmt"
	"io/ioutil"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const modTimeLayout = `2006-01-02 15:04`

func newStatus(width int) *widgets.Paragraph {
	st := widgets.NewParagraph()
	st.Border = false
	return st
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

func (r *Ranger) DisplayError(format string, a ...interface{}) {
	r.statusBar.Text = colorText(fmt.Sprintf(format, a...), colorRed)
	ui.Render(r.statusBar)
}

func (r *Ranger) DisplayFile() {
	if len(r.mainDir) == 0 {
		r.statusBar.Text = ""
		return
	}
	cnt := 1
	info := r.currentFile()
	if info.IsDir() {
		contents, err := ioutil.ReadDir(r.path)
		if err != nil {
			logger.Error(err)
			cnt = -1
		}
		cnt = len(contents)
	}
	r.statusBar.Text = fmt.Sprintf("%s  %d  %s  %s",
		colorText(info.Mode().String(), colorCyan),
		cnt,
		byteCount(info.Size()),
		info.ModTime().Format(modTimeLayout),
	)
}

// shamelessly copied from gitlab.com/pirivan/golearn/cmd/sysprog/ugid
func byteCount(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
