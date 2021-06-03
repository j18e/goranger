package main

import (
	"fmt"
	"os"

	ui "github.com/gizak/termui/v3"
	flags "github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger

	opts struct {
		ChooseFiles string `long:"choosefiles" description:"If specified, will cause ranger to write the paths to all selected files to the given file."`
		SelectFile  string `long:"selectfile" description:"If specified, will cause ranger to open with the given file selected."`
	}
)

func main() {
	logger = logrus.New()

	logFile, err := os.Create("/tmp/goranger.log")
	if err != nil {
		logger.Fatal(err)
	}
	defer logFile.Close()
	logger.SetOutput(logFile)
	logger.SetLevel(logrus.DebugLevel)

	_, err = flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}
	if _, ok := err.(*flags.Error); ok {
		os.Exit(1) // if it's a flags.Error the output is already printed
	}
	if err != nil {
		logrus.Fatal(err)
	}

	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	logger.Debugf("starting with choosefiles=%s selectfile=%s", opts.ChooseFiles, opts.SelectFile)
	if err := ui.Init(); err != nil {
		return fmt.Errorf("failed to initialize termui: %w", err)
	}
	defer ui.Close()

	ranger, err := NewRanger(os.Getenv("PWD"), opts.SelectFile, opts.ChooseFiles)
	if err != nil {
		return err
	}

	return ranger.RunLoop()
}
