package op

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/Phillip-England/pride/internal/cmd"
	"github.com/Phillip-England/pride/internal/server"
	"github.com/Phillip-England/pride/internal/site"
	"github.com/Phillip-England/pride/internal/syserr"
	"github.com/fsnotify/fsnotify"
)

var currentServer *http.Server
var mu sync.Mutex

type OpServe struct {
	Code int
	Cmd  cmd.Cmd
}

func (op *OpServe) Exec(c cmd.Cmd) *syserr.Err {
	cmdServe, ok := c.(*cmd.CmdServe)
	if !ok {
		return syserr.New(syserr.Here(), "type assertion failure")
	}
	port := cmdServe.Port
	if serr := OperationStartServer(port); serr != nil {
		return serr
	}
	dir, serr := site.LoadPrideDir()
	if serr != nil {
		return serr
	}
	serr = startFileWacher(dir.Path, port)
	if serr != nil {
		return serr
	}
	return nil
}


func OperationStartServer(port int) *syserr.Err {
    dir, serr := site.LoadPrideDir()
    if serr != nil {
        return serr
    }
    svr, serr := server.NewServer(port, dir)
    if serr != nil {
        return serr
    }
    srv := &http.Server{
        Addr:    svr.Addr,
        Handler: svr.Mux,
    }
    mu.Lock()
    currentServer = srv
    mu.Unlock()
    go func() {
        fmt.Printf("🚀 serving on port %d\n", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("server error: %s\n", err)
        }
    }()
    return nil
}


func startFileWacher(prideDirPath string, port int) *syserr.Err {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return syserr.New(syserr.Here(), "%s", err.Error())
	}
	defer watcher.Close()
	err = addRecursive(watcher, prideDirPath)
	if err != nil {
		return syserr.New(syserr.Here(), "%s", err.Error())
	}
	errChan := make(chan *syserr.Err, 1)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					errChan <- syserr.New(syserr.Here(), "fsnotify: event channels closed")
					return
				}
				serr := onChange(port, event)
				if serr != nil {
					errChan <- serr
					return
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					errChan <- syserr.New(syserr.Here(), "fsnotify: error channels closed")
					return
				}
				errChan <- syserr.New(syserr.Here(), "%s", err.Error())
				return
			}
		}
	}()
	if serr := <-errChan; serr != nil {
		return serr
	}
	return nil
}

func onChange(port int, event fsnotify.Event) *syserr.Err {
    if event.Op != fsnotify.Write {
        return nil
    }
	fmt.Println("✍️ restarting server..")
    mu.Lock()
    if currentServer != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := currentServer.Shutdown(ctx); err != nil {
            mu.Unlock()
            return syserr.New(syserr.Here(), "failed to shutdown: %s", err.Error())
        }
    }
    mu.Unlock()
    return OperationStartServer(port)
}

func addRecursive(w *fsnotify.Watcher, root string) error {
    return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            if err := w.Add(path); err != nil {
                return err
            }
        }
        return nil
    })
}