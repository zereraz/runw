package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func runMainFile(mainFile string) {
	cmd := exec.Command("go", "run", mainFile)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		// print the error and not close the watcher
		fmt.Printf("Error %v \n %s", err, errBuf.String())
	}
	fmt.Println(out.String())
}

func startWatcher(watchPath, mainFile string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {

					fileExt := event.Name[len(event.Name)-3:]

					if fileExt == ".go" {
						log.Printf("%s changed \n", event.Name)
						log.Printf("running %s \n", mainFile)
						runMainFile(mainFile)
					}

				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()

	err = watcher.Add(watchPath)

	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal(errors.New("Error : arguments empty"))
	}
	startWatcher(args[0], args[1])
	fmt.Println("vim-go")
}
