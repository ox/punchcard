package main

import (
  "errors"
  "flag"
  "fmt"
	"io/ioutil"
	"log"
  "os"

  punchcard "github.com/ox/punchcard"
)

func main() {
  encodingStr := flag.String("encoding", "029ftn", "Card encoding")
  flag.Parse()

  files := flag.Args()
  if len(files) == 0 {
    fmt.Println(errors.New("No card files specified"))
    flag.Usage()
    return
  }

  encoding, err := punchcard.EncodingFromString(*encodingStr)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  for _, card := range files {
    content, err := ioutil.ReadFile(card)
    if err != nil {
      log.Fatal(err)
    }

    text, err := punchcard.Read(content, encoding)
    if err != nil {
      log.Fatalf("Could not read card: %v", err)
    }

    fmt.Print(text)
  }
}
