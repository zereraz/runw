package main

import (
	"io/ioutil"
	"log"
	"path"
)

// get all recursive directories in one array
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
			// ignoring .git and other . folders
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
