/**
 * IO
 */

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
)

// Sio represents all IO functionality
type Sio struct {
	Files    map[string]*SrcFile
	TmpFiles []string
	Routes   []Route
	Models   map[string]Model
}

// GetAllFilePaths recursively looks for golang files starting with a root directory at path `rootPath`
func (s *Sio) GetAllFilePaths(rootPath string) {
	fileInfos, err := ioutil.ReadDir(rootPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range fileInfos {
		fileName := fileInfo.Name()
		switch mode := fileInfo.Mode(); {
		case mode.IsRegular():
			// Check for .go extension
			if fileName[len(fileName)-3:] != ".go" {
				continue
			}

			s.TmpFiles = append(s.TmpFiles, rootPath+"/"+fileName)
			break
		case mode.IsDir():
			s.GetAllFilePaths(rootPath + "/" + fileName)
			break
		}
	}
}

// @deprecated
func readDirFiles(dirPath string) (files []string) {

	fileInfos, err := ioutil.ReadDir(dirPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, fileInfo := range fileInfos {

		switch mode := fileInfo.Mode(); {
		// case mode.IsDir():
		// 	fmt.Printf("DIR: %s\n", fileInfo.Name())
		case mode.IsRegular():
			files = append(files, dirPath+"/"+fileInfo.Name())
			// fmt.Printf("File: %s\n", fileInfo.Name())
		}
	}

	return
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
