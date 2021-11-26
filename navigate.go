package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	ui "github.com/gizak/termui/v3"
)

func (r *Ranger) LevelUp() {
	if err := r.LoadPath(filepath.Dir(r.path), r.path); err != nil {
		logger.Error(err)
	}
}

func (r *Ranger) LevelDown() error {
	f := r.currentFile()
	if f.IsDir() {
		return r.LoadPath(r.pathToSelection(), "")
	}
	if r.chooseFiles != "" {
		content := filepath.Join(r.path, r.baseName())
		if err := ioutil.WriteFile(r.chooseFiles, []byte(content), 0777); err != nil {
			logger.Error(err)
		}
		return ErrExit
	}
	cmd := exec.Command(os.Getenv("EDITOR"), r.pathToSelection())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer func() {
		ui.Close()
		ui.Init()
	}()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", cmd, err)
	}
	return nil
}

func (r *Ranger) Scroll(i int) {
	r.mainPane.ScrollAmount(i)
}

func (r *Ranger) ScrollTop() {
	r.mainPane.ScrollTop()
}

func (r *Ranger) ScrollBottom() {
	r.mainPane.ScrollBottom()
}

func (r *Ranger) SelectedLast() bool {
	return r.mainPane.SelectedRow+1 == len(r.mainPane.Rows)
}
