package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
)

//go:embed public
var staticContent embed.FS
var staticContentFS, _ = fs.Sub(staticContent, "public")

func Static(mux *http.ServeMux) {
	h := http.FileServer(http.FS(staticContentFS))
	// loop through staticContent and print the paths
	printDir(staticContent, ".")

	mux.Handle("/styles/", h)
}

func printDir(fs embed.FS, basePath string) {
	paths, _ := fs.ReadDir(basePath)
	for _, path := range paths {
		if path.IsDir() {
			printDir(fs, filepath.Join(basePath, path.Name()))
		} else {
			slog.Info(filepath.Join(basePath, path.Name()))
		}
	}
}
