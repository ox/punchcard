package punchcard

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// inval is an invalid character, shows up as `~` in the output
var inval uint32 = 00404

var decodings = make(map[Encoding][]uint32)
var encodings = make(map[Encoding][]uint32)

type Encoding interface {
	start() int
	end() int
	table() []uint32

	Encode(*bytes.Buffer) *bytes.Buffer
	Decode(*bytes.Buffer) (*bytes.Buffer, error)
}

func EncodingFromString(s string) (Encoding, error) {
	switch s {
	case "026comm":
		return Zero26Comm{}, nil
	case "026ftn":
		return Zero26Ftn{}, nil
	case "029ftn":
		return Zero29Ftn{}, nil
	case "EBCDIC":
		return EBCDIC{}, nil
	default:
		return Zero26Comm{}, fmt.Errorf("Invalid encoding '%s'", s)
	}
}

func init() {
	encs := []Encoding{
		Zero26Comm{},
		Zero26Ftn{},
		Zero29Ftn{},
		EBCDIC{},
	}

	for _, encoding := range encs {
		decodings[encoding] = make([]uint32, 4096)
		for i := 0; i < len(decodings[encoding]); i++ {
			decodings[encoding][i] = '~'
		}

		table := encoding.table()
		encodings[encoding] = table

		for i := encoding.start(); i <= encoding.end(); i++ {
			decodings[encoding][table[i]] = uint32(i)
		}
	}
}

func encode(contents *bytes.Buffer, e Encoding) *bytes.Buffer {
	encodingTable := encodings[e]
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

	return body
}

func decode(contents *bytes.Buffer, e Encoding) (*bytes.Buffer, error) {
	// verify that it has the H80/H82 header format and skip over it
	format := 80
	endCol := 80

	header := contents.Next(3)

	if bytes.Equal(header, []byte("H80")) {

	} else if bytes.Equal(header, []byte("H82")) {
		endCol = 82
		format = 82
	} else {
		return nil, errors.New("Invalid card, missing valid header")
	}

	decodingTable := decodings[e]
	text := bytes.NewBuffer([]byte{})

	for {
		// First 3 bytes are the card metadata
		metadata := make([]byte, 3)
		if _, err := contents.Read(metadata); err == io.EOF {
			return text, nil
		}

		line := bytes.NewBufferString("")

		// Read `format` columns from the card body
		for i := 0; i < endCol; {
			// Get 3 bytes, and cast them to uint32
			cols := make([]byte, 3)
			if _, err := contents.Read(cols); err != nil {
				return nil, fmt.Errorf("Card body read error: %w", err)
			}

			a := uint32(cols[0])
			b := uint32(cols[1])
			c := uint32(cols[2])

			// Left column
			left := (a << 4) | (b >> 4)
			line.Write([]byte(string(decodingTable[left])))
			i++

			// Right column
			right := ((b & 0017) << 8) | c
			line.Write([]byte(string(decodingTable[right])))
			i++
		}

		if format == 82 {
			text.Write(line.Bytes()[1 : line.Len()-2])
		} else {
			text.ReadFrom(line)
		}

		text.Write([]byte("\n"))
	}
}
