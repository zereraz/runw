package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

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
	fmt.Println(watchPath, mainFile)
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

func getAllDirs(currentDir string, dirList []string) []string {
	curDirList, err := ioutil.ReadDir(currentDir)

	if err != nil {
		log.Fatal(err)
	}
	//base case
	if len(curDirList) == 0 {
		return dirList
	}

	newDirList := []string{}
	for _, file := range curDirList {
		if file.Name()[0] == '.' {
			continue
		}
		if file.IsDir() {
			folderPath := path.Join(currentDir, file.Name())
			// TODO: REFACTOR
			newDirList = append(newDirList, folderPath)
			newDirList = append(append(dirList, getAllDirs(folderPath, dirList)...), newDirList...)
		}
	}

	return newDirList
}

func startRecursiveWatcher(watchDir, mainFile string) {
	dirs := getAllDirs(watchDir, []string{})
	for _, dir := range dirs {
		done := make(chan bool)
		go func() {
			startWatcher(dir, mainFile)
			done <- true
		}()
		<-done
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal(errors.New("Error : arguments empty"))
	}
	recursiveFlag := flag.Bool("r", false, "recursive watching")
	flag.Parse()
	if *recursiveFlag {
		startRecursiveWatcher(args[1], args[2])
	} else {
		startWatcher(args[0], args[1])
	}
}
