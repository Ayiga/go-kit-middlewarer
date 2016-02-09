package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
)

func sliceContains(slice []string, entry string) bool {
	for _, s := range slice {
		if entry == s {
			return true
		}
	}
	return false
}

func determineLocalName(suggestedName string, currentNames []string) string {
	if !sliceContains(currentNames, suggestedName) {
		return suggestedName
	}

	f := func(prefix string) string {
		p := []byte(prefix)
		for i := 97; i < 97+26; i++ {
			b := append([]byte{}, p...)
			b = append(b, byte(i))
			if !sliceContains(currentNames, string(b)) {
				return string(b)
			}
		}
		return ""
	}

	for sliceContains(currentNames, suggestedName) {
		suggestedName = suggestedName + "a"

		suggestedName = f(suggestedName)
	}

	return suggestedName
}

func formatBuffer(buf bytes.Buffer, filename string) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("Warning: internal Error: invalid go generated in file %s: %s", filename, err)
		log.Printf("Warning: compile the package to analyze the error: %s", err)
		return buf.Bytes()
	}
	return src
}

func openFile(dirname, filename string) *os.File {
	_, err := os.Stat(dirname)
	err = os.Mkdir(dirname, 0744)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Unable to Make Directory: %s: %s", dirname, err)
	}
	fname := dirname + "/" + filename

	file, err := os.Create(fname)
	if err != nil && os.IsExist(err) {
		file, err = os.Open(fname)
	}

	if err != nil {
		log.Fatalf("Unable to open or create file %s: %s", fname, err)
	}
	return file
}
