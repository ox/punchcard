package punchcard

import (
	"bytes"
)

type Zero26Ftn struct{}

func (e Zero26Ftn) Encode(contents *bytes.Buffer) *bytes.Buffer {
	return encode(contents, e)
}

func (e Zero26Ftn) Decode(contents *bytes.Buffer) (*bytes.Buffer, error) {
	return decode(contents, e)
}

func (e Zero26Ftn) start() int {
	return ' '
}

func (e Zero26Ftn) end() int {
	return '_'
}

func (e Zero26Ftn) table() []uint32 {
	return zero26FtnCode
}

/* Bare bones 026 kepunch encodings */
var zero26FtnCode = []uint32{
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, inval, inval, inval, inval, /* chars   */
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, inval, inval, inval, inval, /* chars   */
	00000, inval, inval, inval, 02102, inval, inval, 00042, /*  !"#$%&' */
	01042, 04042, 02042, 04000, 01102, 02000, 04102, 01400, /* ()*+,-./ */
	01000, 00400, 00200, 00100, 00040, 00020, 00010, 00004, /* 01234567 */
	00002, 00001, inval, inval, inval, 00102, inval, inval, /* 89:;<=>? */
	inval, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* @ABCDEFG */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* HIJKLMNO */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* PQRSTUVW */
	01004, 01002, 01001, inval, inval, inval, inval, inval, /* XYZ[\]^_ */
	inval, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* `abcdefg */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* hijklmno */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* pqrstuvw */
	01004, 01002, 01001, inval, inval, inval, inval, inval, /* xyz{|}~  */
}
