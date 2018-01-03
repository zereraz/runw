package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

// run provided main folder and print the output
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

	// print output of mainFile
	fmt.Println(out.String())
}

// watching each folder in dirList
// for go file changes and run mainFile
func watchFiles(mainFile string, dirList []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	// make it wait for the go routine
	done := make(chan bool)

	// handler for events from the watcher
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// write operation
				if event.Op&fsnotify.Write == fsnotify.Write {
					fileName := event.Name

					// get file extension
					fileExt := event.Name[len(fileName)-3:]

					if fileExt == ".go" {
						log.Printf("%s changed \n", fileName)
						log.Printf("running %s \n", mainFile)
						runMainFile(mainFile)
					}

				}
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()

	// watch each dir
	for _, dir := range dirList {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// waiting
	<-done
}

func startRecursiveWatcher(watchDir, mainFile string) {
	// to make watcher watch go files in current directory
	initDir := []string{"."}
	dirs := getAllDirs(watchDir, initDir)
	watchFiles(mainFile, dirs)
}

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		log.Fatal(errors.New("Error : arguments empty"))
	}

	// -r recursive flag for watching all the go files
	// in any depth of directory
	recursiveFlag := flag.Bool("r", false, "recursive watching")
	flag.Parse()

	if *recursiveFlag {
		startRecursiveWatcher(args[1], args[2])
	} else {
		watchFiles(args[1], []string{args[0]})
	}
}
