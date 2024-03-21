package xpress

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Watcher struct {
	reloadActive   bool
	SourceDir      string
	Interval       time.Duration
	lastCheckTime  time.Time
	server         *http.Server
	serverStopChan chan struct{}
}

func NewWatcher(server *http.Server, interval time.Duration) *Watcher {
	sourceDir := GetRootPath()
	return &Watcher{
		SourceDir:      sourceDir,
		Interval:       interval,
		lastCheckTime:  time.Now(),
		server:         server,
		reloadActive:   false,
		serverStopChan: make(chan struct{}),
	}
}

func (w *Watcher) Start() {
	log.Println("Waiting for file changes...")

	var mainGoPath string
	filepath.Walk(w.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "main.go" {
			mainGoPath = path
			return filepath.SkipDir // Skip further traversal
		}
		return nil
	})

	if mainGoPath == "" {
		log.Fatal("main.go file not found in source directory")
	}

	for {
		w.checkChanges(mainGoPath)
		time.Sleep(w.Interval)
	}
}

func (w *Watcher) checkChanges(mainGoPath string) {
	now := time.Now()
	if now.Sub(w.lastCheckTime) >= w.Interval {
		err := filepath.Walk(w.SourceDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && filepath.Ext(path) == ".go" {
				modTime := info.ModTime()
				if modTime.After(w.lastCheckTime) {
					log.Printf("🔄 Update Alert: File '%s' has been modified. 🔄\n", info.Name())
					w.lastCheckTime = now
					// Now you can use `mainGoFullPath` in your reload function
					if err := w.reload(mainGoPath); err != nil {
						log.Println("Error reloading application:", err)
					}

				}
			}
			return nil
		})
		if err != nil {
			log.Println("Error walking directory:", err)
		}
	}
}

func (w *Watcher) reload(mainGoFullPath string) error {
	w.stopServer()
	<-w.serverStopChan
	log.Printf("🔄 Reload Alert: Reloading application... 🔄\n")
	cmd := exec.Command("go", "run", mainGoFullPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Error restarting application:", err)
		return err
	}
	return nil
}

func (w *Watcher) stopServer() {
	if w.server != nil {
		if err := w.server.Shutdown(context.Background()); err != nil {
			log.Println("Error stopping server:", err)
		} else {
			close(w.serverStopChan)
		}
	}
}
