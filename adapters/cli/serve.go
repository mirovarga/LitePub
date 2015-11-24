package main

import (
	"net/http"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

const defaultPort = "2703"

func serve(arguments map[string]interface{}) {
	dir := arguments["<dir>"].(string)

	if arguments["--rebuild"].(int) == 1 {
		build(map[string]interface{}{"<dir>": dir})
	}

	port, ok := arguments["--port"].([]string)
	if !ok {
		port[0] = defaultPort
	}

	watch := arguments["--watch"].(int)

	if watch == 1 {
		go watchDirs(dir)
	}

	log.Printf("Running on http://localhost:%s\n", port[0])
	if watch == 1 {
		log.Printf("Rebuilding when posts or templates change\n")
	}
	log.Printf("Ctrl+C to quit\n")

	http.ListenAndServe(":"+port[0], http.FileServer(http.Dir(filepath.Join(dir, outputDir))))
}

func watchDirs(dir string) {
	watcher, _ := fsnotify.NewWatcher()
	defer watcher.Close()

	watcher.Add(filepath.Join(dir, postsDir))
	watcher.Add(filepath.Join(dir, templatesDir))

	for {
		select {
		case <-watcher.Events:
			build(map[string]interface{}{"<dir>": dir})
		}
	}
}
