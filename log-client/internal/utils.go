package internal

import (
	"bufio"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

func WatchFile(filePath string, processLine func(string), stop chan struct{}) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}()

	// start reading from the end of file
	info, _ := file.Stat()
	_, err = file.Seek(info.Size(), 0)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	// setup fsnotify watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			log.Printf("failed to close gRPC connection: %v", err)
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				for {
					line, err := reader.ReadString('\n')

					// no more new lines yet
					if err != nil {
						break
					}
					processLine(line)
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		case <-stop:
			return
		}
	}
}
