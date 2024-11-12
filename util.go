package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mitchellh/go-ps"
)

func checkProcess() tea.Cmd {
	return tea.Every(time.Millisecond*100, func(t time.Time) tea.Msg {
		processes, err := ps.Processes()
		if err != nil {
			return checkProcessMsg{running: false}
		}

		for _, process := range processes {
			if process.Executable() == "eve-online.exe" {
				return checkProcessMsg{running: true}
			}
		}
		return checkProcessMsg{running: false}
	})
}

// Simulate backup process
func performBackup() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		// It confuses me if this happens to quickly. I think it hasn't worked
		time.Sleep(1 * time.Second)

		home, _ := os.UserHomeDir()
		destinationPath := path.Join(home, fmt.Sprintf("eve-settings-%s.zip", time.Now().Format("2006-01-02_15-04")))
		destinationFile, err := os.Create(destinationPath)
		if err != nil {
			return backupCompleteMsg{
				err: err,
			}
		}
		defer destinationFile.Close()

		pathToZip := getEvePath()

		myZip := zip.NewWriter(destinationFile)
		err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if err != nil {
				return err
			}
			relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
			zipFile, err := myZip.Create(relPath)
			if err != nil {
				return err
			}
			fsFile, err := os.Open(filePath)
			if err != nil {
				return err
			}
			_, err = io.Copy(zipFile, fsFile)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return backupCompleteMsg{
				err: err,
			}
		}
		err = myZip.Close()
		if err != nil {
			return backupCompleteMsg{
				err: err,
			}
		}
		return backupCompleteMsg{
			filename: destinationPath,
		}
	})
}
