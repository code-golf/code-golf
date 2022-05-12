package hole

import (
	"math/rand"
	"strings"
)

func hexdump() ([]string, string) {
	args := []string{
		"Code Golf",
		" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~",
		"multi-\nline\nstring",
		"Code Golf is a game designed to let you show off your code-fu by solving problems in the least number of characters.",
		`In mathematics and computing, the hexadecimal (also base 16 or simply hex)
numeral system is a positional numeral system that represents numbers using a
radix (base) of 16. Unlike the decimal system representing numbers using 10
symbols, hexadecimal uses 16 distinct symbols, most often the symbols "0"-"9"
to represent values 0 to 9, and "A"-"F" (or alternatively "a"-"f") to represent
values from 10 to 15.

Software developers and system designers widely use hexadecimal numbers because
they provide a human-friendly representation of binary-coded values. Each
hexadecimal digit represents four bits (binary digits), also known as a nibble
(or nybble). For example, an 8-bit byte can have values ranging from 00000000 to
11111111 in binary form, which can be conveniently represented as 00 to FF in
hexadecimal.`,
		"0123456789abcdef",
		"Lorem ipsum dolor sit amet,\n\n consectetur adipiscing elit\n",
	}

	outs := []string{
		"00000000: 436f 6465 2047 6f6c 66                   Code Golf",
		`00000000: 2021 2223 2425 2627 2829 2a2b 2c2d 2e2f   !"#$%&'()*+,-./
00000010: 3031 3233 3435 3637 3839 3a3b 3c3d 3e3f  0123456789:;<=>?
00000020: 4041 4243 4445 4647 4849 4a4b 4c4d 4e4f  @ABCDEFGHIJKLMNO
00000030: 5051 5253 5455 5657 5859 5a5b 5c5d 5e5f  PQRSTUVWXYZ[\]^_
00000040: 6061 6263 6465 6667 6869 6a6b 6c6d 6e6f  ` + "`" + `abcdefghijklmno
00000050: 7071 7273 7475 7677 7879 7a7b 7c7d 7e    pqrstuvwxyz{|}~`,
		`00000000: 6d75 6c74 692d 0a6c 696e 650a 7374 7269  multi-.line.stri
00000010: 6e67                                     ng`,
		`00000000: 436f 6465 2047 6f6c 6620 6973 2061 2067  Code Golf is a g
00000010: 616d 6520 6465 7369 676e 6564 2074 6f20  ame designed to
00000020: 6c65 7420 796f 7520 7368 6f77 206f 6666  let you show off
00000030: 2079 6f75 7220 636f 6465 2d66 7520 6279   your code-fu by
00000040: 2073 6f6c 7669 6e67 2070 726f 626c 656d   solving problem
00000050: 7320 696e 2074 6865 206c 6561 7374 206e  s in the least n
00000060: 756d 6265 7220 6f66 2063 6861 7261 6374  umber of charact
00000070: 6572 732e                                ers.`,
		`00000000: 496e 206d 6174 6865 6d61 7469 6373 2061  In mathematics a
00000010: 6e64 2063 6f6d 7075 7469 6e67 2c20 7468  nd computing, th
00000020: 6520 6865 7861 6465 6369 6d61 6c20 2861  e hexadecimal (a
00000030: 6c73 6f20 6261 7365 2031 3620 6f72 2073  lso base 16 or s
00000040: 696d 706c 7920 6865 7829 0a6e 756d 6572  imply hex).numer
00000050: 616c 2073 7973 7465 6d20 6973 2061 2070  al system is a p
00000060: 6f73 6974 696f 6e61 6c20 6e75 6d65 7261  ositional numera
00000070: 6c20 7379 7374 656d 2074 6861 7420 7265  l system that re
00000080: 7072 6573 656e 7473 206e 756d 6265 7273  presents numbers
00000090: 2075 7369 6e67 2061 0a72 6164 6978 2028   using a.radix (
000000a0: 6261 7365 2920 6f66 2031 362e 2055 6e6c  base) of 16. Unl
000000b0: 696b 6520 7468 6520 6465 6369 6d61 6c20  ike the decimal
000000c0: 7379 7374 656d 2072 6570 7265 7365 6e74  system represent
000000d0: 696e 6720 6e75 6d62 6572 7320 7573 696e  ing numbers usin
000000e0: 6720 3130 0a73 796d 626f 6c73 2c20 6865  g 10.symbols, he
000000f0: 7861 6465 6369 6d61 6c20 7573 6573 2031  xadecimal uses 1
00000100: 3620 6469 7374 696e 6374 2073 796d 626f  6 distinct symbo
00000110: 6c73 2c20 6d6f 7374 206f 6674 656e 2074  ls, most often t
00000120: 6865 2073 796d 626f 6c73 2022 3022 2d22  he symbols "0"-"
00000130: 3922 0a74 6f20 7265 7072 6573 656e 7420  9".to represent
00000140: 7661 6c75 6573 2030 2074 6f20 392c 2061  values 0 to 9, a
00000150: 6e64 2022 4122 2d22 4622 2028 6f72 2061  nd "A"-"F" (or a
00000160: 6c74 6572 6e61 7469 7665 6c79 2022 6122  lternatively "a"
00000170: 2d22 6622 2920 746f 2072 6570 7265 7365  -"f") to represe
00000180: 6e74 0a76 616c 7565 7320 6672 6f6d 2031  nt.values from 1
00000190: 3020 746f 2031 352e 0a0a 536f 6674 7761  0 to 15...Softwa
000001a0: 7265 2064 6576 656c 6f70 6572 7320 616e  re developers an
000001b0: 6420 7379 7374 656d 2064 6573 6967 6e65  d system designe
000001c0: 7273 2077 6964 656c 7920 7573 6520 6865  rs widely use he
000001d0: 7861 6465 6369 6d61 6c20 6e75 6d62 6572  xadecimal number
000001e0: 7320 6265 6361 7573 650a 7468 6579 2070  s because.they p
000001f0: 726f 7669 6465 2061 2068 756d 616e 2d66  rovide a human-f
00000200: 7269 656e 646c 7920 7265 7072 6573 656e  riendly represen
00000210: 7461 7469 6f6e 206f 6620 6269 6e61 7279  tation of binary
00000220: 2d63 6f64 6564 2076 616c 7565 732e 2045  -coded values. E
00000230: 6163 680a 6865 7861 6465 6369 6d61 6c20  ach.hexadecimal
00000240: 6469 6769 7420 7265 7072 6573 656e 7473  digit represents
00000250: 2066 6f75 7220 6269 7473 2028 6269 6e61   four bits (bina
00000260: 7279 2064 6967 6974 7329 2c20 616c 736f  ry digits), also
00000270: 206b 6e6f 776e 2061 7320 6120 6e69 6262   known as a nibb
00000280: 6c65 0a28 6f72 206e 7962 626c 6529 2e20  le.(or nybble).
00000290: 466f 7220 6578 616d 706c 652c 2061 6e20  For example, an
000002a0: 382d 6269 7420 6279 7465 2063 616e 2068  8-bit byte can h
000002b0: 6176 6520 7661 6c75 6573 2072 616e 6769  ave values rangi
000002c0: 6e67 2066 726f 6d20 3030 3030 3030 3030  ng from 00000000
000002d0: 2074 6f0a 3131 3131 3131 3131 2069 6e20   to.11111111 in
000002e0: 6269 6e61 7279 2066 6f72 6d2c 2077 6869  binary form, whi
000002f0: 6368 2063 616e 2062 6520 636f 6e76 656e  ch can be conven
00000300: 6965 6e74 6c79 2072 6570 7265 7365 6e74  iently represent
00000310: 6564 2061 7320 3030 2074 6f20 4646 2069  ed as 00 to FF i
00000320: 6e0a 6865 7861 6465 6369 6d61 6c2e       n.hexadecimal.`,
		"00000000: 3031 3233 3435 3637 3839 6162 6364 6566  0123456789abcdef",
		`00000000: 4c6f 7265 6d20 6970 7375 6d20 646f 6c6f  Lorem ipsum dolo
00000010: 7220 7369 7420 616d 6574 2c0a 0a20 636f  r sit amet,.. co
00000020: 6e73 6563 7465 7475 7220 6164 6970 6973  nsectetur adipis
00000030: 6369 6e67 2065 6c69 740a                 cing elit.`,
	}

	rand.Shuffle(len(args), func(i, j int) {
		args[i], args[j] = args[j], args[i]
		outs[i], outs[j] = outs[j], outs[i]
	})

	return args, strings.Join(outs, "\n\n")
}
