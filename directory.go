package main

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gizak/termui/v3/widgets"
)

func (r *Ranger) ReloadDirs() error {
	mainDir, err := r.RenderDir(r.mainPane, r.path, r.baseName())
	if err != nil {
		return err
	}
	r.mainDir = mainDir

	parentDir, err := r.RenderDir(r.parentPane, filepath.Dir(r.path), filepath.Base(r.path))
	if err != nil {
		return err
	}
	r.parentDir = parentDir
	return nil
}

func (r *Ranger) RenderDir(l *widgets.List, path, selectFile string) ([]fs.FileInfo, error) {
	logger.Debugf("render %s, select '%s'", path, selectFile)
	d, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	dir := r.filterDir(d)
	l.Rows = colorFiles(dir)
	l.SelectedRow = fileIdx(selectFile, dir)
	return dir, nil
}

func (r *Ranger) filterDir(dir []fs.FileInfo) []fs.FileInfo {
	if r.showHidden {
		return dir
	}
	var res []fs.FileInfo
	for _, f := range dir {
		if !strings.HasPrefix(f.Name(), ".") {
			res = append(res, f)
		}
	}
	return res
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
