# Punchcard Reader

This program reads a binary file into the punchcard format specified in [Emulated Punched Card Decks](https://homepage.divms.uiowa.edu/~jones/cards/format.html). After reading the deck and card preambles, it converts 12 bits at a time into ASCII. The end goal is to be able to pipe these into Fortran or Algol programs and function as sort of an OS.

The codebase references [`cardlist.txt`](https://homepage.divms.uiowa.edu/~jones/cards/cardlist.txt), and [`cardcode.i.txt`](https://homepage.divms.uiowa.edu/~jones/cards/cardcode.i.txt); both by Douglas Jones.

## Installation

```
go install github.com/ox/punchcard/cmd/punchcard-reader
```

## Usage

You can read some of the example cards like so:

```
$ punchcard-reader ./brainfuck.card
...
```

## The Card Format

Punchcards encoded 12 bits/column and usually had 80 columns. A few special cards had 82 columns but also required a compatible machine. Card columns were laid out with the top 3 bits for the zone and the next 9 bits for the numeric portion.

```
Top                  Bottom
  _ _ _ _ _ _ _ _ _ _ _ _
 |_|_|_|_|_|_|_|_|_|_|_|_|
12 11 0 1 2 3 4 5 6 7 8 9
 |     |                 |
 |Zone |     Numeric     |
```

In the card file format the columns are laid out sequentially, which requires reading 3 bytes at a time to reconstruct the 12 bits for each column.

```
column 1                column 2
|_ _ _ _ _ _ _ _ _ _ _ _|_ _ _ _ _ _ _ _ _ _ _ _|
|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|
|               |               |               |
byte 1          byte 2          byte 3
```

Card files are distinguished by the "magic" `H80` prefix.

> Card files need a distinguished magic number or prefix to prevent accidental interpretation of random files as virtual card decks. Here, we will use the ASCII prefix "H80"

Then each card in the file starts with 3 bytes of metadata about the look and format of the card.

```
byte 1          byte 2         byte 3
|_ _ _ _ _ _ _ _|_ _ _ _ _ _ _ _|_ _ _ _ _ _ _ _|
|1|_|_|_|_|_|_|_|1|_|_|_|_|_|_|_|1|_|_|_|_|_|_|_|
| | Color | |cut| | |     |form | |    logo     |
     |       |   |
   corner    | punch
           interp
```


The `punchcard-reader` tool supports:

- [x] H80/H82 encoded cards
- [x] Read 029 Encoding
- [x] Read 026comm Encoding
- [x] Read 026ftn Encoding
- [x] Read EBCDIC Encoding
- [ ] Card metadata

The full and original reference for this file format is available at [Emulated Punched Card Decks](https://homepage.divms.uiowa.edu/~jones/cards/format.html).
