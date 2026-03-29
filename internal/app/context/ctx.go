package context

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Context includes working directory and current cmd history.
type Context struct {
	WorkingDirectory string
	History          []string
}

type ContextMgr struct {
	CurrentCtx *Context
}

func parseHistoryLine(line string) string {
	// : 1774125697:0;git commit -m "msg"
	if idx := strings.Index(line, ";"); idx != -1 {
		return strings.TrimSpace(line[idx+1:])
	}
	return strings.TrimSpace(line)
}

func (cm *ContextMgr) FetchContext() {
	if cm.CurrentCtx == nil {
		cm.CurrentCtx = &Context{}
	}

	// get working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cm.CurrentCtx.WorkingDirectory = cwd

	cm.CurrentCtx.History = nil

	// get history file
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

	var all []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		cmd := parseHistoryLine(line)
		if cmd != "" {
			all = append(all, cmd)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// keep latest 50
	if len(all) > 50 {
		all = all[len(all)-50:]
	}

	cm.CurrentCtx.History = all
}

// display WD
func (cm *ContextMgr) DisplayWD() string {
	return cm.CurrentCtx.WorkingDirectory
}

// display history
func (cm *ContextMgr) DisplayCmdHistory(l int) []string {
	h := cm.CurrentCtx.History

	if l <= 0 {
		return []string{}
	}

	if l >= len(h) {
		return h
	}

	return h[len(h)-l:]
}
