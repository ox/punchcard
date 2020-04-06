package punchcard

import (
	"fmt"
)

/* cardcode.i
 *
 * card code arrays are indexed by 7-bit ASCII code, and
 * give corresponding 12-bit card codes using the indicated
 * collating sequence.
 *
 * inval should be externally defined, either as an illegal
 * card code (on conversion from ASCII to card codes) or as
 * a code with a bit set outside the least significant 12.
 *
 * original author:  Douglas Jones, jones@cs.uiowa.edu
 * revisions:
 *	    March 5, 1996
 *	    Feb  18, 1997 to add 026 and EBCDIC converstion tables
 */

type Encoding int

func (e Encoding) String() string {
	switch e {
	case Zero26comm:
		return "Zero26comm"
	case Zero26ftn:
		return "Zero26ftn"
	case Zero29ftn:
		return "Zero29ftn"
	case EBCDIC:
		return "EBCDIC"
	default:
		return "Unknown"
	}
}

const (
	Zero26comm Encoding = iota
	Zero26ftn
	Zero29ftn
	EBCDIC
)

var decodings = make(map[Encoding][]uint32)
var encodings = make(map[Encoding][]uint32)

func GetEncodingTable(e Encoding) []uint32 {
	return encodings[e]
}

func GetDecodingTable(e Encoding) []uint32 {
	return decodings[e]
}

func EncodingFromString(s string) (Encoding, error) {
	switch s {
	case "026comm":
		return Zero26comm, nil
	case "026ftn":
		return Zero26ftn, nil
	case "029ftn":
		return Zero29ftn, nil
	case "EBCDIC":
		return EBCDIC, nil
	default:
		return Zero26comm, fmt.Errorf("Invalid encoding '%s'", s)
	}
}

func init() {
	for _, encoding := range []Encoding{Zero26comm, Zero26ftn, Zero29ftn, EBCDIC} {
		decodings[encoding] = make([]uint32, 4096)
		for i := 0; i < len(decodings[encoding]); i += 1 {
			decodings[encoding][i] = '~'
		}

		start := ' '
		end := '_'
		table := zero_26_comm_code

		switch encoding {
		case Zero26comm:
			break
		case Zero26ftn:
			table = zero_26_ftn_code
			break
		case Zero29ftn:
			table = zero_29_code
			break
		case EBCDIC:
			start = 0
			end = 0177
			table = ebcdic_code
		}

		encodings[encoding] = table

		for i := start; i <= end; i += 1 {
			decodings[encoding][table[i]] = uint32(i)
		}
	}
}

// inval is an invalid character, these will show up as `~` in the output
var inval uint32 = 00404

var zero_26_comm_code = []uint32{
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, inval, inval, inval, inval, /* chars   */
	inval, inval, inval, inval, inval, inval, inval, inval, /* control */
	inval, inval, inval, inval, inval, inval, inval, inval, /* chars   */
	00000, inval, inval, 00102, 02102, 01042, 04000, inval, /*  !"#$%&' */
	inval, inval, 02042, inval, 01102, 02000, 04102, 01400, /* ()*+,-./ */
	01000, 00400, 00200, 00100, 00040, 00020, 00010, 00004, /* 01234567 */
	00002, 00001, inval, inval, 04042, inval, inval, inval, /* 89:;<=>? */
	00042, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* @ABCDEFG */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* HIJKLMNO */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* PQRSTUVW */
	01004, 01002, 01001, inval, inval, inval, inval, inval, /* XYZ[\]^_ */
	inval, 04400, 04200, 04100, 04040, 04020, 04010, 04004, /* `abcdefg */
	04002, 04001, 02400, 02200, 02100, 02040, 02020, 02010, /* hijklmno */
	02004, 02002, 02001, 01200, 01100, 01040, 01020, 01010, /* pqrstuvw */
	01004, 01002, 01001, inval, inval, inval, inval, inval, /* xyz{|}~  */
}

/* Bare bones 026 kepunch encodings */
var zero_26_ftn_code = []uint32{
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

/* DEC's version of the IBM 029 kepunch encoding, (thus avoiding IBM's
   use of non-ASCII punctuation), based on that given in the appendix
   to Digital's "Small Computer Handbook, 1973", and augmented to
   translate lower case to upper case.  As a result of this modification,
   inversion of this table should be done with care! */
var zero_29_code = []uint32{
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

/* FULL EBCDIC, from Appendix C of System 360 Programming by Alex Thomas,
   1977, Reinhart Press, San Francisco.  Codes not in that table have been
   left compatable with DEC's 029 table.  Some control codes have been
   left out */
var ebcdic_code = []uint32{
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
