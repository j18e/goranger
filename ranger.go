package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	ErrExit = errors.New("exit")

	textStyle        = ui.NewStyle(ui.ColorWhite)
	selectedRowStyle = ui.NewStyle(ui.ColorClear, ui.ColorBlue, ui.ModifierBold)
)

type Ranger struct {
	height, width int
	path          string

	parentPane  *widgets.List
	mainPane    *widgets.List
	previewPane *widgets.Paragraph
	statusBar   *widgets.Paragraph

	parentDir []fs.FileInfo
	mainDir   []fs.FileInfo

	events      <-chan ui.Event
	prevKey     string
	chooseFiles string // file to write chosen file names to, if specified

	showHidden bool
}

func newList() *widgets.List {
	l := widgets.NewList()
	l.Border = false
	l.TextStyle = textStyle
	l.SelectedRowStyle = selectedRowStyle
	return l
}

func newParagraph() *widgets.Paragraph {
	p := widgets.NewParagraph()
	p.WrapText = false
	p.Border = false
	return p
}

func NewRanger(path, chooseFiles string) (*Ranger, error) {
	termX, termY := ui.TerminalDimensions()

	mainPane := newList()
	mainPane.SetRect(termX/5, 0, termX/5*2, termY-2)
	parentPane := newList()
	parentPane.SetRect(0, 0, termX/5, termY-2)
	previewPane := newParagraph()
	previewPane.SetRect(termX/5*2, 0, termX, termY-2)
	statusBar := newParagraph()
	statusBar.SetRect(0, termY-1, termX, termY)

	r := &Ranger{
		height: termY,
		width:  termX,

		parentPane:  parentPane,
		mainPane:    mainPane,
		previewPane: previewPane,
		statusBar:   statusBar,

		chooseFiles: chooseFiles,
		events:      ui.PollEvents(),
	}

	if err := r.LoadPath(path, ""); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Ranger) baseName() string {
	f := r.currentFile()
	if f == nil {
		return ""
	}
	return f.Name()
}

func (r *Ranger) currentFile() fs.FileInfo {
	if len(r.mainDir) == 0 {
		return nil
	}
	return r.mainDir[r.mainPane.SelectedRow]
}

func (r *Ranger) pathToSelection() string {
	return filepath.Join(r.path, r.baseName())
}

func (r *Ranger) LoadPath(path, selectFile string) error {
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	logger.Debugf("load %s, select '%s'", path, selectFile)
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		if selectFile == "" {
			selectFile = filepath.Base(path)
		}
		path = filepath.Dir(path)
	}
	r.path = path

	if err := r.ReloadDirs(selectFile); err != nil {
		return err
	}

	if err := r.updatePreview(); err != nil {
		return err
	}
	r.UpdateStatus()

	r.render()
	return nil
}

func (r *Ranger) RunLoop() error {
	for {
		e := <-r.events
		if e.Type != ui.KeyboardEvent {
			continue
		}
		logger.Debugf("received event: %s", e.ID)
		if err := r.HandleEvent(e); err != nil {
			if err == ErrExit {
				return nil
			}
			r.DisplayError("%s", err)
		}
		r.render()
	}
}

func (r *Ranger) render() {
	ui.Render(r.mainPane)
	ui.Render(r.parentPane)
	ui.Render(r.previewPane)
	ui.Render(r.statusBar)
}
