package main

type Font map[rune][]byte

// Must only be created with constructor
type Text struct {
	fontNormal, fontSlim Font
}

func NewText() *Text {
	t := Text{}

	t.fontNormal = Font{
		'0': {0x3E, 0x41, 0x41, 0x3E},
		'1': {0x00, 0x02, 0x7F, 0x00},
		'2': {0x62, 0x51, 0x49, 0x46},
		'3': {0x22, 0x49, 0x49, 0x36},
		'4': {0x18, 0x14, 0x7F, 0x10},
		'5': {0x27, 0x49, 0x49, 0x31},
		'6': {0x3E, 0x49, 0x49, 0x30},
		'7': {0x01, 0x71, 0x0D, 0x03},
		'8': {0x36, 0x49, 0x49, 0x36},
		'9': {0x06, 0x49, 0x49, 0x3E},
	}
	t.fontSlim = Font{
		'0': {0x3E, 0x41, 0x3E},
		'1': {0x00, 0x02, 0x7F},
		'2': {0x62, 0x51, 0x4E},
		'3': {0x41, 0x49, 0x36},
		'4': {0x0F, 0x08, 0x7F},
		'5': {0x4F, 0x49, 0x31},
		'6': {0x3E, 0x49, 0x31},
		'7': {0x61, 0x19, 0x07},
		'8': {0x36, 0x49, 0x36},
		'9': {0x46, 0x49, 0x3E},
	}
	return &t
}

// Centers the text on the screen
func (t *Text) MakeText(str string) *LedArray {
	// Select font
	// TODO: Only count valid characters
	runes := []rune(str)
	var font Font
	if len(runes) > 4 {
		font = t.fontSlim
	} else {
		font = t.fontNormal
	}

	// Get characters
	ledColumns := []byte{}
	for _, char := range runes {
		charBytes := font[char]
		if len(charBytes) > 0 {
			// Add space
			if len(ledColumns) > 0 {
				ledColumns = append(ledColumns, 0)
			}

			// Add characters
			ledColumns = append(ledColumns, charBytes...)
		}
	}

	// Center text
	x := (BoardWidth - len(ledColumns)) / 2
	if x < 0 {
		// More columns than screen can show, shrink columns
		ledColumns = ledColumns[-x : -x+BoardWidth]
		x = 0
	}

	// Add columns
	ledArray := LedArray{}
	for _, column := range ledColumns {
		for y := byte(0); y < BoardHeight; y++ {
			// Turn on LED
			if (column>>y)&1 == 1 {
				ledArray[x][y] = true
			}
		}
		x++
	}
	return &ledArray
}
