package main

import (
	"fmt"
	"io/fs"
	"os"
)

const (
	colorBlue   = "blue"
	colorYellow = "yellow"
	colorRed    = "red"
	colorCyan   = "cyan"
)

func colorFiles(files []fs.FileInfo) []string {
	var res []string
	for _, f := range files {
		s := f.Name()
		if f.IsDir() {
			s = colorText(s, colorBlue)
		} else if isExec(f.Mode()) {
			s = colorText(s, colorYellow)
		}
		res = append(res, s)
	}
	return res
}

func colorText(s, c string) string {
	return fmt.Sprintf("[%s](fg:%s)", s, c)
}

func isExec(mode os.FileMode) bool {
	return mode&0111 != 0
}
