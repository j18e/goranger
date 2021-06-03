package main

import (
	"errors"
	"io/fs"
	"io/ioutil"
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

func NewRanger(path, selectFile, chooseFiles string) (*Ranger, error) {
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

	if err := r.LoadPath(path, selectFile); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Ranger) baseName() string {
	return r.currentFile().Name()
}

func (r *Ranger) currentFile() fs.FileInfo {
	return r.mainDir[r.mainPane.SelectedRow]
}

func (r *Ranger) pathToSelection() string {
	return filepath.Join(r.path, r.baseName())
}

func (r *Ranger) LoadPath(path, selectFile string) error {
	f, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !f.IsDir() {
		path = filepath.Dir(path)
	}
	r.path = path

	mainDir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	r.mainPane.Rows = colorFiles(mainDir)
	r.mainPane.SelectedRow = fileIdx(selectFile, mainDir)
	r.mainDir = mainDir

	parentDir, err := ioutil.ReadDir(filepath.Dir(path))
	if err != nil {
		return err
	}
	r.parentPane.Rows = colorFiles(parentDir)
	r.parentPane.SelectedRow = fileIdx(path, parentDir)
	r.parentDir = parentDir

	if err := r.updatePreview(); err != nil {
		return err
	}
	r.DisplayFile()

	r.render()
	return nil
}

func fileIdx(name string, dir []fs.FileInfo) int {
	name = filepath.Base(name)
	if name == "" {
		return 0
	}
	for i, f := range dir {
		if f.Name() == name {
			return i
		}
	}
	return 0
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
			return err
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

func (r *Ranger) HandleEvent(e ui.Event) error {
	switch e.ID {
	case ":":
		// TODO handle commands
		// cmd := r.Bar.HandleCommand(r.events)
		// switch cmd {
		// case "TmuxNavigateRight":
		// default:
		// 	break
		// }
		// logger.Debugf("got user command: %s", cmd)

	case "/":
		// TODO handle search

	case "q", "<C-c>":
		return ErrExit
	case "l", "<Right>":
		if err := r.LevelDown(); err != nil {
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
