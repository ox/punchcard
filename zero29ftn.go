package punchcard

import (
	"bytes"
)

type Zero29Ftn struct{}

func (e Zero29Ftn) Encode(contents *bytes.Buffer) *bytes.Buffer {
	return encode(contents, e)
}

func (e Zero29Ftn) Decode(contents *bytes.Buffer) (*bytes.Buffer, error) {
	return decode(contents, e)
}

func (e Zero29Ftn) start() int {
	return ' '
}

func (e Zero29Ftn) end() int {
	return '_'
}

func (e Zero29Ftn) table() []uint32 {
	return zero29FtnCode
}

/* DEC's version of the IBM 029 kepunch encoding, (thus avoiding IBM's
   use of non-ASCII punctuation), based on that given in the appendix
   to Digital's "Small Computer Handbook, 1973", and augmented to
   translate lower case to upper case.  As a result of this modification,
   inversion of this table should be done with care! */
var zero29FtnCode = []uint32{
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, inval, inval, inval, inval, /* chars   */
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, inval, inval, inval, inval, /* chars   */
	00000, 02202, 00006, 00102, 02102, 01042, 04000, 00022, /*  !"#$%&' */
	04022, 02022, 02042, 04012, 01102, 02000, 04102, 01400, /* ()*+,-./ */
	01000, 00400, 00200, 00100, 00040, 00020, 00010, 00004, /* 01234567 */
	00002, 00001, 00202, 02012, 04042, 00012, 01012, 01006, /* 89:;<=>? */
	00042, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* @ABCDEFG */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* HIJKLMNO */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* PQRSTUVW */
	01004, 01002, 01001, 04202, 02006, 01202, 04006, 01022, /* XYZ[\]^_ */
	inval, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* `abcdefg */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* hijklmno */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* pqrstuvw */
	01004, 01002, 01001, inval, inval, inval, inval, inval, /* xyz{|}~  */
}
