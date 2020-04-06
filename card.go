package punchcard

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type card struct {
	encoding Encoding
	content  *bytes.Buffer
}

func (c *card) Bytes() []byte {
	return c.content.Bytes()
}

func New(encoding Encoding, contents *bytes.Buffer) (*card, error) {
	body := bytes.NewBufferString("H80")

	lines := make([][]byte, 0)

	for contents.Len() > 0 {
		line := make([]byte, 80)

		for i := 0; i < 80; {
			byte, err := contents.ReadByte()
			if err == io.EOF {
				break
			}

			// blank out the rest of the card if we find a newline
			if byte == '\n' {
				for ; i < 80; i++ {
					line[i] = ' '
				}
			} else if byte == '\t' {
				for k := i; k < 80 && k < i+4; k++ {
					line[i] = ' '
				}
				i += 4
			} else {
				line[i] = byte
				i++
			}
		}

		lines = append(lines, line)
	}

	encodingTable := GetEncodingTable(encoding)
	if len(encodingTable) == 0 {
		return nil, fmt.Errorf("Could not get encoding table '%s'", encoding)
	}

	for _, line := range lines {
		body.Write([]byte{0x80, 0x80, 0x80})

		for i := 0; i < len(line); {
			even := byte(line[i])
			odd := byte(line[i+1])
			i += 2

			evenCol := encodingTable[even]
			oddCol := encodingTable[odd]

			a := evenCol >> 4
			b := ((evenCol & 017) << 4) | (oddCol >> 8)
			c := oddCol & 00377

			body.Write([]byte{byte(a), byte(b), byte(c)})
		}
	}

	return &card{encoding, body}, nil
}

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
	ascii_code := GetDecodingTable(encoding)
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
			text.Write(line.Bytes()[1 : line.Len()-2])
		} else {
			text.ReadFrom(line)
		}

		text.Write([]byte("\n"))
	}

	return text, nil
}
