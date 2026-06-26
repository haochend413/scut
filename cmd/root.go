package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/haochend413/scut/config"
	"github.com/haochend413/scut/internal/app"
	"github.com/haochend413/scut/internal/db"
	"github.com/haochend413/scut/internal/models"
	"github.com/haochend413/scut/internal/ui"
	"github.com/spf13/cobra"
)

func parseHistoryLine(line string) string {
	if idx := strings.Index(line, ";"); idx != -1 {
		return strings.TrimSpace(line[idx+1:])
	}
	return strings.TrimSpace(line)
}

var write string
var addlast bool

var globalCfg *config.Config
var globalDB *db.DB
var globalApp *app.App

var rootCmd = &cobra.Command{
	Use:   "scut",
	Short: "scut — per-directory command shortcut manager",
	Long:  "scut — per-directory command shortcut manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadOrCreateConfig()
		globalCfg = &cfg

		dbPath := cfg.DataFilePath + "/shortcuts.db"

		var err error
		globalDB, err = db.NewDB(dbPath)
		if err != nil {
			log.Fatalf("Failed to connect to database %q: %v\n", dbPath, err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		globalApp = app.NewApp(globalDB)

		// write mode: scut -w "command"
		if write != "" {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			sc := models.Shortcut{WorkDirectory: cwd, Command: write}
			if err := globalApp.AddShortcut(sc); err != nil {
				fmt.Fprintf(os.Stderr, "scut: %v\n", err)
			}
			return
		}

		// add-last mode: scut -l
		if addlast {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			u, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}

			historyPath := filepath.Join(u.HomeDir, ".zsh_history")
			f, err := os.Open(historyPath)
			if err != nil {
				historyPath = filepath.Join(u.HomeDir, ".bash_history")
				f, err = os.Open(historyPath)
				if err != nil {
					log.Fatal(err)
				}
			}
			defer f.Close()

			var last, secondLast string
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					secondLast = last
					last = line
				}
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			if secondLast == "" {
				log.Fatal("history file does not have enough entries")
			}

			sc := models.Shortcut{
				WorkDirectory: cwd,
				Command:       parseHistoryLine(secondLast),
			}
			if err := globalApp.AddShortcut(sc); err != nil {
				fmt.Fprintf(os.Stderr, "scut: %v\n", err)
			}
			return
		}

		// default: launch TUI
		model := ui.NewModel(globalApp)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	defer func() {
		if globalDB != nil {
			globalDB.Close()
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "scut: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&write, "write", "w", "", "save a command to the current directory")
	rootCmd.Flags().BoolVarP(&addlast, "last", "l", false, "save the previous shell command")
}
