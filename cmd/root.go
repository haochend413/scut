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
var globalModel *ui.Model

var rootCmd = &cobra.Command{
	Use:   "scut",
	Short: "scut",
	Long:  "scut",
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
		// Initialize application with AppState
		globalApp = app.NewApp(globalDB)

		// write mode: scut -w "command"
		if write != "" {
			cwd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			sc := models.Shortcut{
				WorkDirectory: cwd,
				Command:       write,
			}

			globalApp.AddShortcut(sc)
			globalApp.OnClose()
			return
		}

		// add latest mode: scut -l
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

			secondLastCmd := parseHistoryLine(secondLast)

			shortcut := models.Shortcut{
				WorkDirectory: cwd,
				Command:       secondLastCmd,
			}
			globalApp.AddShortcut(shortcut)
			globalApp.OnClose()
			return
		}

		// default: launch UI
		model := ui.NewModel(globalApp)
		globalModel = &model

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
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&write, "write", "w", "", "write command")
	rootCmd.Flags().BoolVarP(&addlast, "last", "l", false, "add last command from shell history")
}
