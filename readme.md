# Punchcard Reader

This program reads a binary file into the punchcard format specified in [Emulated Punched Card Decks](https://homepage.divms.uiowa.edu/~jones/cards/format.html). After reading the deck and card preambles, it converts 12 bits at a time into ASCII. The end goal is to be able to pipe these into Fortran or Algol programs and function as sort of an OS.

The codebase references [`cardlist.txt`](https://homepage.divms.uiowa.edu/~jones/cards/cardlist.txt), and [`cardcode.i.txt`](https://homepage.divms.uiowa.edu/~jones/cards/cardcode.i.txt); both by Douglas Jones.

Features:

- [x] Read H80 cards
- [x] Read H82 cards
- [x] Read 029 Encoding
- [x] Read 026comm Encoding
- [x] Read 026ftn Encoding
- [x] Read EBCDIC Encoding

## Installation

```
go install github.com/ox/punchcard/cmd/punchcard
```

## Usage

You can read some of the example cards like so:

```
$ punchcard ./brainfuck.card
...
```

---

> a card-image file, 12 bits/column, 80 columns/card.

Notes from the "Emulated Punched Card Decks" page:

> One column of a card holds 12 bits; in the file, we lay them out as follows, with ones representing punched holes:

```
Top                  Bottom
  _ _ _ _ _ _ _ _ _ _ _ _
 |_|_|_|_|_|_|_|_|_|_|_|_|
12 11 0 1 2 3 4 5 6 7 8 9
 |     |                 |
 |Zone |     Numeric     |
```

```
column 1                column 2
|_ _ _ _ _ _ _ _ _ _ _ _|_ _ _ _ _ _ _ _ _ _ _ _|
|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|_|
|               |               |               |
byte 1          byte 2          byte 3
```

> Card files need a distinguished magic number or prefix to prevent accidental interpretation of random files as virtual card decks. Here, we will use the ASCII prefix "H80"
