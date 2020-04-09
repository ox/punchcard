package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	punchcard "github.com/ox/punchcard"
)

func main() {
	encodingStr := flag.String("encoding", "029ftn", "Card encoding")
	outputStr := flag.String("output", "", "File to write to")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		log.Println("No input files specified")
		os.Exit(2)
		return
	}

	enc, err := punchcard.EncodingFromString(*encodingStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	allFiles := bytes.NewBuffer([]byte{})

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		allFiles.Write(content)
	}

	buff := enc.Encode(allFiles)

	if *outputStr != "" {
		if err := ioutil.WriteFile(*outputStr, buff.Bytes(), 0644); err != nil {
			log.Fatalf("Could not write to %s: %s", *outputStr, err)
		}

		return
	}

	fmt.Print(buff.Bytes())
}
