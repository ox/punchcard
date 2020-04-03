package punchcard

import (
  "bytes"
  "errors"
  "fmt"
  "io"
)

func Read(card []byte, encoding Encoding) (*bytes.Buffer, error) {
  // verify that it has the H80/H82 header format and skip over it
  format := 80
  endCol := 80

  if bytes.Equal(card[0:3], []byte("H80")) {

  } else if bytes.Equal(card[0:3], []byte("H82")) {
    endCol = 82
    format = 82
  } else {
    return nil, errors.New("Invalid card, missing valid header")
  }

  buff := bytes.NewBuffer(card[3:])
  ascii_code := GetEncodingTable(encoding)
  text := bytes.NewBuffer([]byte{})

  for {
    // First 3 bytes are the card metadata
    metadata := make([]byte, 3)
    if _, err := buff.Read(metadata); err == io.EOF {
      return text, nil
    }

    line := bytes.NewBufferString("")

    // Read `format` columns from the card body
    for i := 0; i < endCol; {
      // Get 3 bytes, and cast them to uint32
      cols := make([]byte, 3)
      if _, err := buff.Read(cols); err != nil {
        return nil, fmt.Errorf("Card body read error: %w", err)
      }

      a := uint32(cols[0])
      b := uint32(cols[1])
      c := uint32(cols[2])

      // Left column
      left := (a << 4) | (b >> 4)
      line.Write([]byte(string(ascii_code[left])))
      i += 1

      // Right column
      right := ((b & 0017) << 8) | c
      line.Write([]byte(string(ascii_code[right])))
      i += 1
    }

    if format == 82 {
      text.Write(line.Bytes()[1:line.Len() - 2])
    } else {
      text.ReadFrom(line)
    }

    text.Write([]byte("\n"))
  }

  return text, nil
}
