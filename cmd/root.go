package cmd

import (
	"fmt"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/haochend413/scut/config"
	"github.com/haochend413/scut/internal/app"
	"github.com/haochend413/scut/internal/db"
	"github.com/haochend413/scut/internal/ui"
)

var globalDB *db.DB

var rootCmd = &cobra.Command{
	Use:   "scut",
	Short: "scut — per-directory command shortcut manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadOrCreateConfig()
		dbPath := cfg.DataFilePath + "/shortcuts.db"

		var err error
		globalDB, err = db.NewDB(dbPath)
		if err != nil {
			log.Fatalf("scut: failed to open database: %v\n", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		globalApp := app.NewApp(globalDB)
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
