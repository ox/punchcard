package punchcard

import (
	"bytes"
)

type EBCDIC struct{}

func (e EBCDIC) Encode(contents *bytes.Buffer) *bytes.Buffer {
	return encode(contents, e)
}

func (e EBCDIC) Decode(contents *bytes.Buffer) (*bytes.Buffer, error) {
	return decode(contents, e)
}

func (e EBCDIC) start() int {
	return 0
}

func (e EBCDIC) end() int {
	return 0177
}

func (e EBCDIC) table() []uint32 {
	return ebcdicCode
}

/* FULL EBCDIC, from Appendix C of System 360 Programming by Alex Thomas,
   1977, Reinhart Press, San Francisco.  Codes not in that table have been
   left compatable with DEC's 029 table.  Some control codes have been
   left out */
var ebcdicCode = []uint32{
	05403, inval, inval, inval, inval, inval, inval, inval, /* control */
	02011, 04021, 01021, inval, 04041, 02021, inval, inval, /* chars   */
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, 01201, inval, inval, inval, /* chars   */
	00000, 02202, 00006, 00102, 02102, 01042, 04000, 00022, /*  !"#$%&' */
	04022, 02022, 02042, 04012, 01102, 02000, 04102, 01400, /* ()*+,-./ */
	01000, 00400, 00200, 00100, 00040, 00020, 00010, 00004, /* 01234567 */
	00002, 00001, 00202, 02012, 04042, 00012, 01012, 01006, /* 89:;<=>? */
	00042, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* @ABCDEFG */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* HIJKLMNO */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* PQRSTUVW */
	01004, 01002, 01001, 04202, 02006, 01202, 04006, 01022, /* XYZ[\]^_ */
	inval, 05400, 05200, 05100, 05040, 05020, 05010, 05004, /* `abcdefg */
	05002, 05001, 06400, 06200, 06100, 06040, 06020, 06010, /* hijklmno */
	06004, 06002, 06001, 03200, 03100, 03040, 03020, 03010, /* pqrstuvw */
	03004, 03002, 03001, inval, inval, inval, inval, inval, /* xyz{|}~  */
}
