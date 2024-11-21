package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

			match, _ := regexp.MatchString("state[a-zA-Z0-9\\-_]*\\.json$", filePath)
			if !match {
				return nil
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

func performInitialLoad() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		profileManager, err := NewProfileManager()
		if err != nil {
			return initLoadCompleteMsg{
				err: err,
			}
		}
		return initLoadCompleteMsg{
			profileManager: profileManager,
		}
	})
}

func loadProfile(profileManager *ProfileManager, profile string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := profileManager.SwitchProfile(profile)
		return switchedProfileMsg{
			err: err,
		}
	})
}

func createProfile(profileManager *ProfileManager, profileName, baseProfile string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := profileManager.CreateProfile(profileName, baseProfile)
		return createdProfileMsg{
			err: err,
		}
	})
}
