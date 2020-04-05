package main

import (
  "bytes"
  "fmt"
  "flag"
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

  encoding, err := punchcard.EncodingFromString(*encodingStr)
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

  card, err := punchcard.New(encoding, allFiles)
  if err != nil {
    log.Fatalf("Could not make new card: %w", err)
  }

  if *outputStr != "" {
    if err := ioutil.WriteFile(*outputStr, card.Bytes(), 0644); err != nil {
      log.Fatalf("Could not write to %s: %w", *outputStr, err)
    }

    return
  }

  fmt.Println(card.Bytes())
}
