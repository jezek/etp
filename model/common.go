package model

import "fmt"

type Command struct {
	Description string
	Function    interface{}
}

// ESC/POS commands
var CommonCommands map[string]Command = map[string]Command{
	"init":   {"Initialize printer", Init},
	"cr":     {"Print and carriage return", Cr},
	"lf":     {"Line feed", Lf},
	"pf":     {"(n byte) Print and feed n lines", Pf},
	"b":      {"Print following text with bold", Bold},
	"nob":    {"Print following text without bold", NoBold},
	"ds":     {"Print following text with double size", DoubleSize},
	"nods":   {"Print following text with normal size", NoDoubleSize},
	"u":      {"Print following text with underline 1 dot", Underline1},
	"u2":     {"Print following text with underline 2 dot", Underline2},
	"nou":    {"Print following text without underline", NoUnderline},
	"left":   {"Print following text aligned left", AlignLeft},
	"center": {"Print following text aligned center", AlignCenter},
	"right":  {"Print following text aligned right", AlignRight},
	"fontA":  {"Print following text using font A", FontA},
	"fontB":  {"Print following text using font B", FontB},
	"font": {`(string = [ABubwh]) Specify print mode using string options.
	Options:
	 A - font A
	 B - font B
	 u - underlined
	 b - bold
	 w - double width
	 h - double height`, PrintMode},
	"marginLeft":       {"(int16) Set left margin", MarginLeft},
	"printRegionWidth": {"(int16) Set print region width", PrintRegionWidth},
}

// TODO codepages
// [Name] Select character code table
// [Format] ASCII
// Hex
// Decimal
// [Range] TM-J2000/J2100, TM-T90, TM-T88IV, TM-T70, TM-L90, TM-P60:
// 0 ≤ n ≤ 5, 16 ≤ n ≤ 19, n = 255
// TM-U230:
// 0 ≤ n ≤ 5, n = 16, 254, 255 (Other than the following models)
// 0 ≤ n ≤ 8, n = 16, 254, 255 (Japanese model)
// TM-U220:
// 0 ≤ n ≤ 5, 16 ≤ n ≤ 19, n = 254, 255 (Other than the following models)
// 0 ≤ n ≤ 8, 16 ≤ n ≤ 19, n = 254, 255 (Japanese model)
// [Default] n = 0
// ESC t
// n
// 1B 74 n
// 27 116 n
// [Printers not featuring this command] None
// [Description]
// Paper roll
// Selects a page n from the character code table as follows:
// n Character code table
// 0 Page 0 [PC437 (U.S.A., Standard Europe)]
// 1 Page 1 [Katakana]
// 2 Page 2 [PC850 (Multilingual)]
// 3 Page 3 [PC860 (Portuguese)]
// 4 Page 4 [PC863 (Canadian-French)]
// 5 Page 5 [PC865 (Nordic)]
// 6 Page 6 [Simplified Kanji, Hirakana]
// 7 Page 7 [Simplified Kanji]
// 8 Page 8 [Simplified Kanji]
// 16 Page 16 [WPC1252]
// 17 Page 17 [PC866 (Cyrillic #2)]
// 18 Page 18 [PC852 (Latin 2)]
// 19 Page 19 [PC858 (Euro)]
// 254 Page 254
// 255 Page 255
// ■ The characters of each page are the same for alphanumeric parts (ASCII code: Hexadecimal = 20H to 7FH /
// Decimal = 32 to 127 20H to 7FH), and different for the escape character parts (ASCII code: Hexadecimal =
// 80H to FFH / Decimal = 128 to 255 80H to FFH).
// ■ The selected character code table is valid until ESC @ is executed, the printer is reset, or the power is
// turned off.
// [Model-dependent variations]
// TM-J2000/J2100, TM-T90, TM-T88IV, TM-T70, TM-L90, TM-P60, TM-U230,
// TM-U220.
// [Notes]
// TM-J2000/J2100, TM-T90, TM-L90
// Page 255 is able to be edited by <Function 7> ~ <Function 10> of GS ( E. When the printer is
// shipped, the page is a space page.
// TM-T88IV, TM-T70
// Page 255 is a space page.
// TM-P60
// Page 255 is able to be edited by <Function 7> ~ <Function 10> of GS ( E. When the printer is
// shipped, the page is a space page.
// Settings of this command do not affect special font (24 × 48) printing. Special fonts (24 × 48) print
// page 0[PC437(USA, Standard Europe)] characters irrespective of the settings of this command.
// TM-U230, TM-U220
// Page 254 and 255 are space pages.

// Set CodePage852 (Latin-2) encoding
// ESC 0x74='t' 18
func CP852() []byte {
	return []byte{0x1b, 0x74, 0x12}
}

// Initialize printer
// ESC 0x40='@'
func Init() string {
	return string([]byte{0x1b, 0x40})
}

// Print and carriage return
func Cr() string {
	return string(byte(0xd))
}

// Line feed
func Lf() string {
	return string(byte(0xa))
}

// Print and feed n lines
// ESC 0x64='d' n
func Pf(n byte) string {
	return string([]byte{0x1b, 0x64, n})
}

// Print following text with bold
func Bold() string {
	return string([]byte{0x1b, 0x45, 0x01})
}

// Print following text without bold
func NoBold() string {
	return string([]byte{0x1b, 0x45, 0x00})
}

// Print following text with double size
func DoubleSize() string {
	return string([]byte{0x1b, 0x47, 0x01})
}

// Print following text with normal size
func NoDoubleSize() string {
	return string([]byte{0x1b, 0x47, 0x00})
}

// Specify/cancel underline mode for following text
// ESC 0x2d='-' (0:cancel, 1:one dot width, 2:two-dot width)
func Underline(i byte) (string, error) {
	if i > 2 {
		return "", fmt.Errorf("unknown underline mode")
	}
	return string([]byte{0x1b, 0x2d, i}), nil
}

// Print following text with underline 1 dot
func Underline1() string { u, _ := Underline(1); return u }

// Print following text with underline 2 dot
func Underline2() string { u, _ := Underline(2); return u }

// Print following text without underline
func NoUnderline() string { u, _ := Underline(0); return u }

// Align text if on begining of line
// ESC 0x61='a' (0:left, 1:center, 2:right)
func Align(i byte) (string, error) {
	if i > 2 {
		return "", fmt.Errorf("unknown align mode")
	}
	return string([]byte{0x1b, 0x61, i}), nil
}

// Print following text aligned left
func AlignLeft() string { a, _ := Align(0); return a }

// Print following text aligned center
func AlignCenter() string { a, _ := Align(1); return a }

// Print following text aligned right
func AlignRight() string { a, _ := Align(2); return a }

// Select character font
// ESC 0x4d='M' (0:fontA, 1:fontB):
func Font(i byte) (string, error) {
	if i > 1 {
		return "", fmt.Errorf("unknown font type")
	}
	return string([]byte{0x1b, 0x4d, i}), nil
}

// Print following text using font A
func FontA() string { f, _ := Font(0); return f }

// Print following text using font B
func FontB() string { f, _ := Font(1); return f }

// Specify print mode using string
// A - font A
// B - font B
// u - underlined
// b - bold
// w - double width
// h - double height
func PrintMode(s string) (string, error) {
	used := map[byte]bool{}
	pm := byte(0)
	for _, a := range []byte(s) {
		if used[a] {
			return "", fmt.Errorf("duplicate '%s'", a)
		}
		switch a {
		case 'A':
			if used['B'] {
				return "", fmt.Errorf("allready selected font B")
			}
		case 'B':
			if used['A'] {
				return "", fmt.Errorf("allready selected font A")
			}
			pm |= 1
		case 'b':
			pm |= 1 << 3
		case 'h':
			pm |= 1 << 4
		case 'w':
			pm |= 1 << 5
		case 'u':
			pm |= 1 << 7
		default:
			return "", fmt.Errorf("unknown '%s'", a)
		}
		used[a] = true
	}
	if !used['A'] && !used['B'] {
		return "", fmt.Errorf("no font selected")
	}
	return string([]byte{0x1b, 0x21, pm}), nil
}

// Set left margin
func MarginLeft(i int16) string {
	l := byte(i)
	h := byte(i >> 8)
	return string([]byte{0x1d, 0x4c, l, h})
}

// Set print region width
func PrintRegionWidth(i int16) string {
	l := byte(i)
	h := byte(i >> 8)
	return string([]byte{0x1d, 0x57, l, h})
}
