package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	fileSizeLimit = 1024 * 1024 // 1MB
)

func (r *Ranger) updatePreview() error {
	previewFile := r.path
	if len(r.mainDir) > 0 {
		previewFile = filepath.Join(previewFile, r.baseName())
	}

	f, err := os.Stat(previewFile)
	if err != nil {
		return err
	}
	if f.IsDir() {
		contents, err := ioutil.ReadDir(previewFile)
		if err != nil {
			return err
		}
		if len(contents) == 0 {
			r.previewPane.Text = "empty directory"
			return nil
		}
		txt := strings.Join(colorFiles(contents), "\n")
		r.previewPane.Text = txt
		return nil
	}
	if f.Size() > fileSizeLimit {
		r.previewPane.Text = "file too large to display"
		return nil
	}
	if f.Size() == 0 {
		r.previewPane.Text = "empty file"
		return nil
	}
	bs, err := ioutil.ReadFile(previewFile)
	if err != nil {
		return err
	}
	txt := strings.ReplaceAll(string(bs), "\t", "    ")
	r.previewPane.Text = txt
	return nil
}
