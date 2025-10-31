package main

import (
	"bibFileCheck/checkFile"
	"os"
)

func main() {
	file, err := os.Open("Bib.bib")
	if err != nil {
		panic("Unable to open file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	checkFile.CheckFile(file)
}
