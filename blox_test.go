package blox_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/sa6mwa/blox"
	"github.com/stretchr/testify/assert"
)

const (
	loremIpsum string = "Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui" + blox.LineBreak +
		"augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu" + blox.LineBreak +
		"nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit" + blox.LineBreak +
		"eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus" + blox.LineBreak +
		"non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque." + blox.LineBreak +
		"Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl" + blox.LineBreak +
		"urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum" + blox.LineBreak +
		"suspendisse." // With or without final new line.
)

func TestNew(t *testing.T) {
	refBlox := &blox.Blox{
		Columns:         blox.InitialCanvasColumns,
		Rows:            blox.InitialCanvasRows,
		LineSpacing:     blox.InitialLineSpacing,
		TrimRightSpaces: blox.InitialTrimRightSpaces,
		Canvas:          make([]rune, blox.InitialCanvasColumns*blox.InitialCanvasRows, blox.InitialCanvasCapacity),
		Cursor:          blox.CursorPosition{blox.InitialCursorPositionX, blox.InitialCursorPositionY, true},
	}
	if len(refBlox.Canvas) > 0 {
		for i := 0; i < len(refBlox.Canvas); i++ {
			refBlox.Canvas[i] = ' '
		}
	}
	b := blox.New()
	assert.Equal(t, refBlox, b, "Should be equal")

	b = blox.New().SetColumnsAndRows(3, 3)
	refCanvas := []rune("         ")
	assert.Equal(t, refCanvas, b.Canvas)
}

func TestWipe(t *testing.T) {
	b := blox.New().SetColumnsAndRows(10, 10).Wipe()
	assert.Empty(t, b.Canvas)
}

func TestResizeCanvas(t *testing.T) {
	b := blox.New().SetColumnsAndRows(100, 100)
	assert.Len(t, b.Canvas, 100*100)
}

func TestSetColumns(t *testing.T) {
	b := blox.New().SetColumnsAndRows(1, 1).SetColumns(2)
	refCanvas := []rune("  ")
	assert.Equal(t, refCanvas, b.Canvas)
}

func TestRows(t *testing.T) {
	b := blox.New().SetColumnsAndRows(2, 1).SetRows(2)
	refCanvas := []rune("    ")
	assert.Equal(t, refCanvas, b.Canvas)
}

func TestSetLineSpacing(t *testing.T) {
	b := blox.New().SetLineSpacing(3)
	assert.Equal(t, 3, b.LineSpacing)
}

func TestSetTrimRightSpaces(t *testing.T) {
	b := blox.New().SetTrimRightSpaces(false)
	assert.False(t, b.TrimRightSpaces)

	refString := "HELLO     " + blox.LineBreak
	got := b.SetColumnsAndRows(10, 1).Move(0, 0).PutText("HELLO").String()
	assert.Equal(t, refString, got)

	b.SetTrimRightSpaces(true)
	assert.True(t, b.TrimRightSpaces)

	refString = "HELLO" + blox.LineBreak
	got = b.Wipe().ResizeCanvas().Move(0, 0).PutText("HELLO").String()
	assert.Equal(t, refString, got)
}

func TestSetTrimFinalEmptyLines(t *testing.T) {
	b := blox.New().SetColumnsAndRows(10, 3).PutText("HELLO")

	refCanvas := []rune("HELLO                         ")
	assert.Equal(t, refCanvas, b.Canvas)

	refString := "HELLO     " + blox.LineBreak + "          " + blox.LineBreak + "          " + blox.LineBreak
	got := b.SetTrimRightSpaces(false).SetTrimFinalEmptyLines(false).String()
	assert.Equal(t, refString, got)

	refString = "HELLO     " + blox.LineBreak
	got = b.SetTrimFinalEmptyLines(true).String()
	assert.Equal(t, refString, got)

	refString = "HELLO" + blox.LineBreak
	got = b.SetTrimRightSpaces(true).String()
	assert.Equal(t, refString, got)
}

func TestMove(t *testing.T) {
	b := blox.New()
	c := blox.New().Move(13, 37)
	assert.Equal(t, b, c)

	b.SetColumnsAndRows(3, 3)
	c.SetColumnsAndRows(3, 3).Move(3, 3)
	b.Cursor.X = 2
	b.Cursor.Y = 2
	b.Cursor.OffCanvas = true
	assert.Equal(t, b, c)

	b = blox.New().SetColumnsAndRows(3, 3)
	b.Cursor.X = 1
	b.Cursor.Y = 2
	c = blox.New().SetColumnsAndRows(3, 3).Move(1, 2)
	assert.Equal(t, b, c)
}

func TestMoveX(t *testing.T) {
	b := blox.New().SetColumnsAndRows(3, 3).Move(2, 0)
	c := blox.New().SetColumnsAndRows(3, 3).Move(1, 0)
	c.MoveX(2)
	assert.Equal(t, b, c)
}

func TestMoveY(t *testing.T) {
	b := blox.New().SetColumnsAndRows(3, 3).Move(1, 2)
	c := blox.New().SetColumnsAndRows(3, 3).Move(1, 0)
	c.MoveY(2)
	assert.Equal(t, b, c)
}
func TestMoveRight(t *testing.T) {
	b := blox.New().SetColumnsAndRows(3, 3).Move(2, 0)
	c := blox.New().SetColumnsAndRows(3, 3).MoveRight(2)
	assert.Equal(t, b, c)
}

func TestMoveLeft(t *testing.T) {
	b := blox.New().SetColumnsAndRows(3, 3).Move(0, 1)
	c := blox.New().SetColumnsAndRows(3, 3).Move(2, 1).MoveLeft(2)
	assert.Equal(t, b, c)
}

func TestMoveDown(t *testing.T) {
	b := blox.New().SetColumnsAndRows(3, 3).Move(1, 2)
	c := blox.New().SetColumnsAndRows(3, 3).Move(1, 0).MoveDown(2)
	assert.Equal(t, b, c)
}

func TestMoveUp(t *testing.T) {
	b := blox.New().SetColumnsAndRows(3, 3).Move(1, 0)
	c := blox.New().SetColumnsAndRows(3, 3).Move(1, 2).MoveUp(2)
	assert.Equal(t, b, c)
	c.MoveUp(3)
	assert.Equal(t, b, c)
}

func TestPutLine(t *testing.T) {
	b := blox.New().SetColumnsAndRows(10, 5).SetTrimFinalEmptyLines(true).SetTrimRightSpaces(true).
		Move(2, 1).PutLine([]rune("HELLO"))
	expect := blox.LineBreak + "  HELLO" + blox.LineBreak
	got := b.String()
	assert.Equal(t, expect, got)
}

func TestPutLines(t *testing.T) {
	b := blox.New().SetColumnsAndRows(10, 5)
	lines := []string{"HELLO", "WORLD"}
	b.PutLines(lines...)

	expect := "HELLO" + blox.LineBreak + "WORLD" + blox.LineBreak + blox.LineBreak + blox.LineBreak + blox.LineBreak
	got := b.String()
	assert.Equal(t, expect, got)

	b.MoveRight(3).PutLines(lines...).SetTrimFinalEmptyLines(true)
	expect = "HELLO" + blox.LineBreak + "WORLD" + blox.LineBreak + "   HELLO" + blox.LineBreak + "   WORLD" + blox.LineBreak
	got = b.String()
	assert.Equal(t, expect, got)
}

func TestPutText(t *testing.T) {
	b := blox.New().SetColumnsAndRows(15, 10).SetTrimFinalEmptyLines(true).SetTrimFinalEmptyLines(true)
	text := "+--------+" + blox.LineBreak
	text += "| HELLO  |" + blox.LineBreak
	text += "+--------+" + blox.LineBreak
	b.PutText(text)
	assert.Equal(t, text, b.String())

	b.Move(2, 2).PutText(text)

	text2 := "+--------+" + blox.LineBreak
	text2 += "| HELLO  |" + blox.LineBreak
	text2 += "+-+--------+" + blox.LineBreak
	text2 += "  | HELLO  |" + blox.LineBreak
	text2 += "  +--------+" + blox.LineBreak
	assert.Equal(t, text2, b.String())
}

func ExampleBlox_PutText() {
	b := blox.New().Trim().SetColumnsAndRows(80, 24)

	text := "ABCDE FGHIJ KLMNO" + blox.LineBreak
	text += "PQRST UVWXY ZABCD" + blox.LineBreak

	heading := "CRYPTO" + blox.LineBreak
	heading += "GROUPS"

	b.PutText(heading).DrawHorizontalLine(0, 6).DrawVerticalLine(0, 1, ':').
		PutChar('+').MoveDown().MoveX(0).PutText(text).
		Move(9, 0).PutText(text).PrintCanvas()

	// Output:
	// CRYPTO : ABCDE FGHIJ KLMNO
	// GROUPS : PQRST UVWXY ZABCD
	// -------+
	// ABCDE FGHIJ KLMNO
	// PQRST UVWXY ZABCD
}

func ExampleBlox_PutText_second() {
	b := blox.New().Trim().SetColumnsAndRows(80, 24)

	text := "ABCDE FGHIJ KLMNO" + blox.LineBreak
	text += "PQRST UVWXY ZABCD" + blox.LineBreak

	heading := "CRYPTO" + blox.LineBreak
	heading += "GROUPS"

	str := b.PutText(heading).DrawHorizontalLine(0, 6).DrawVerticalLine(0, 1, ':').
		PutChar('+').MoveDown().MoveX(0).PutText(text).
		Move(9, 0).PutText(text).Join(blox.LineBreak)

	strEndingInNewLine := b.String()

	fmt.Println(str)
	fmt.Println("--")
	fmt.Print(strEndingInNewLine)

	// Output:
	// CRYPTO : ABCDE FGHIJ KLMNO
	// GROUPS : PQRST UVWXY ZABCD
	// -------+
	// ABCDE FGHIJ KLMNO
	// PQRST UVWXY ZABCD
	// --
	// CRYPTO : ABCDE FGHIJ KLMNO
	// GROUPS : PQRST UVWXY ZABCD
	// -------+
	// ABCDE FGHIJ KLMNO
	// PQRST UVWXY ZABCD
}

func ExampleBlox_PutTextRightAligned() {
	b := blox.New().Trim().SetColumnsAndRows(80, 24)

	// blox.LineBreak is \n on Linux/Darwin and \r\n on Windows.

	text := "Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui" + blox.LineBreak
	text += "augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu" + blox.LineBreak
	text += "nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit" + blox.LineBreak
	text += "eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus" + blox.LineBreak
	text += "non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque." + blox.LineBreak
	text += "Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl" + blox.LineBreak
	text += "urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum" + blox.LineBreak
	text += "suspendisse." + blox.LineBreak

	box := "+----------------------------+" + blox.LineBreak
	box += "|       A BOX WITH TEXT      |" + blox.LineBreak
	box += "+----------------------------+" + blox.LineBreak

	b.PutText(text).MoveY(4).PutTextRightAligned(box).PrintCanvas()

	// Output:
	// Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui
	// augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu
	// nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit
	// eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus
	// non a cursus dui urna, facilisis ipsum dapibus eti+----------------------------+
	// Primis pharetra cursus ultrices vel curabitur duis|       A BOX WITH TEXT      |
	// urna turpis mauris maecenas ac diam, posuere morbi+----------------------------+
	// suspendisse.
}

func ExampleBlox_PrintCanvas() {
	b := blox.New().SetColumnsAndRows(80, 24).SetTrimRightSpaces(true).SetTrimFinalEmptyLines(true)

	text := "Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui" + blox.LineBreak
	text += "augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu" + blox.LineBreak
	text += "nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit" + blox.LineBreak
	text += "eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus" + blox.LineBreak
	text += "non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque." + blox.LineBreak
	text += "Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl" + blox.LineBreak
	text += "urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum" + blox.LineBreak
	text += "suspendisse." + blox.LineBreak

	box := "+----------------------------+" + blox.LineBreak
	box += "|       A BOX WITH TEXT      |" + blox.LineBreak
	box += "+----------------------------+" + blox.LineBreak

	b.PutText(text).Move(13, 3).PutText(box)

	b.PrintCanvas()

	// Output:
	// Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui
	// augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu
	// nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit
	// eleifend null+----------------------------+llam suspendisse quis volutpat luctus
	// non a cursus |       A BOX WITH TEXT      |bus etiam odio lacus feugiat neque.
	// Primis pharet+----------------------------+ur duis taciti semper, tortor nisl
	// urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum
	// suspendisse.
}

func ExampleBlox_FprintCanvas() {
	b := blox.New().SetColumnsAndRows(80, 24).Trim()

	text := "Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui" + blox.LineBreak
	text += "augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu" + blox.LineBreak
	text += "nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit" + blox.LineBreak
	text += "eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus" + blox.LineBreak
	text += "non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque." + blox.LineBreak
	text += "Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl" + blox.LineBreak
	text += "urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum" + blox.LineBreak
	text += "suspendisse." + blox.LineBreak

	box := "+----------------------------+" + blox.LineBreak
	box += "|       A BOX WITH TEXT      |" + blox.LineBreak
	box += "+----------------------------+" + blox.LineBreak

	b.PutText(text).MoveX(23).DrawVerticalLine(0, 7, ':').
		Move(13, 3).PutText(box).FprintCanvas(os.Stdout).Move(0, 0)

	// Output:
	// Lorem ipsum dolor sit a:et consectetur adipiscing elit torquent ante tortor dui
	// augue, dictumst convall:s eget tempor pharetra lectus magnis lacinia lacus eu
	// nostra. Sagittis dolor :attis laoreet justo mollis est varius etiam nisl, sit
	// eleifend null+----------------------------+llam suspendisse quis volutpat luctus
	// non a cursus |       A BOX WITH TEXT      |bus etiam odio lacus feugiat neque.
	// Primis pharet+----------------------------+ur duis taciti semper, tortor nisl
	// urna turpis mauris maec:nas ac diam, posuere morbi mi class tincidunt cum
	// suspendisse.           :
}

func ExampleBlox_DrawSeparator() {
	b := blox.New().Trim().SetColumnsAndRows(80, 15)

	text := "Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui" + blox.LineBreak
	text += "augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu" + blox.LineBreak
	text += "nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit" + blox.LineBreak
	text += "eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus" + blox.LineBreak
	text += "non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque." + blox.LineBreak
	text += "Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl" + blox.LineBreak
	text += "urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum" + blox.LineBreak
	text += "suspendisse." // With or without final new line.

	b.PutText(text).DrawSeparator().PutText("Final line.").PrintCanvas()

	// Output:
	// Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui
	// augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu
	// nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit
	// eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus
	// non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque.
	// Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl
	// urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum
	// suspendisse.
	// --------------------------------------------------------------------------------
	// Final line.
}

func ExampleBlox_DrawSplit() {
	b := blox.New().Trim().SetColumnsAndRows(80, 9)
	text := "Lorem ipsum dolor sit amet consectetur adipiscing elit torquent ante tortor dui" + blox.LineBreak
	text += "augue, dictumst convallis eget tempor pharetra lectus magnis lacinia lacus eu" + blox.LineBreak
	text += "nostra. Sagittis dolor mattis laoreet justo mollis est varius etiam nisl, sit" + blox.LineBreak
	text += "eleifend nullam magna aptent erat vitae. Nullam suspendisse quis volutpat luctus" + blox.LineBreak
	text += "non a cursus dui urna, facilisis ipsum dapibus etiam odio lacus feugiat neque." + blox.LineBreak
	text += "Primis pharetra cursus ultrices vel curabitur duis taciti semper, tortor nisl" + blox.LineBreak
	text += "urna turpis mauris maecenas ac diam, posuere morbi mi class tincidunt cum" + blox.LineBreak
	text += "suspendisse." // With or without final new line.

	str := b.PutText(text).DrawSeparator().MoveX(20).DrawSplit().String()
	fmt.Print(str)

	// Output:
	// Lorem ipsum dolor si| amet consectetur adipiscing elit torquent ante tortor dui
	// augue, dictumst conv|llis eget tempor pharetra lectus magnis lacinia lacus eu
	// nostra. Sagittis dol|r mattis laoreet justo mollis est varius etiam nisl, sit
	// eleifend nullam magn| aptent erat vitae. Nullam suspendisse quis volutpat luctus
	// non a cursus dui urn|, facilisis ipsum dapibus etiam odio lacus feugiat neque.
	// Primis pharetra curs|s ultrices vel curabitur duis taciti semper, tortor nisl
	// urna turpis mauris m|ecenas ac diam, posuere morbi mi class tincidunt cum
	// suspendisse.        |
	// --------------------|-----------------------------------------------------------
}

func TestRowAndColumnCount(t *testing.T) {
	text := loremIpsum
	col, row := blox.RowAndColumnCount(text)
	assert.Equal(t, 8, row)
	assert.Equal(t, 80, col)
	text += blox.LineBreak
	assert.Equal(t, 8, row)
	assert.Equal(t, 80, col)
}

func TestLineCount(t *testing.T) {
	text := loremIpsum
	assert.Equal(t, 1, blox.LineCount("ONELINER"))
	assert.Equal(t, 8, blox.LineCount(text))
	text += blox.LineBreak
	assert.Equal(t, 8, blox.LineCount(text))
}

func TestMaximumLineLength(t *testing.T) {
	text := loremIpsum
	assert.Equal(t, 9, blox.MaximumLineLength("HELLO\nBEAUTIFUL\nWORLD.\n"))
	assert.Equal(t, 80, blox.MaximumLineLength(text))
}

func ExampleCutLinesShort() {
	text := loremIpsum
	fmt.Print(blox.CutLinesShort(text, 13, true))
	// Output:
	// Lorem ipsum d
	// augue, dictum
	// nostra. Sagit
	// eleifend null
	// non a cursus
	// Primis pharet
	// urna turpis m
	// suspendisse.
}

func ExampleCutLineShort() {
	line := "Hello world, I am Blox."
	line = blox.CutLineShort(line, 12, true)
	fmt.Println(line)
	// Output:
	// Hello world…
}

// WrapString is (C) 2014 Mitchell Hashimoto https://github.com/mitchellh
// Imported from https://github.com/mitchellh/go-wordwrap
func TestWrapString(t *testing.T) {
	cases := []struct {
		Input, Output string
		Lim           uint
	}{
		// A simple word passes through.
		{
			"foo",
			"foo",
			4,
		},
		// A single word that is too long passes through.
		// We do not break words.
		{
			"foobarbaz",
			"foobarbaz",
			4,
		},
		// Lines are broken at whitespace.
		{
			"foo bar baz",
			"foo\nbar\nbaz",
			4,
		},
		// Lines are broken at whitespace, even if words
		// are too long. We do not break words.
		{
			"foo bars bazzes",
			"foo\nbars\nbazzes",
			4,
		},
		// A word that would run beyond the width is wrapped.
		{
			"fo sop",
			"fo\nsop",
			4,
		},
		// Do not break on non-breaking space.
		{
			"foo bar\u00A0baz",
			"foo\nbar\u00A0baz",
			10,
		},
		// Whitespace that trails a line and fits the width
		// passes through, as does whitespace prefixing an
		// explicit line break. A tab counts as one character.
		{
			"foo\nb\t r\n baz",
			"foo\nb\t r\n baz",
			4,
		},
		// Trailing whitespace is removed if it doesn't fit the width.
		// Runs of whitespace on which a line is broken are removed.
		{
			"foo    \nb   ar   ",
			"foo\nb\nar",
			4,
		},
		// An explicit line break at the end of the input is preserved.
		{
			"foo bar baz\n",
			"foo\nbar\nbaz\n",
			4,
		},
		// Explicit break are always preserved.
		{
			"\nfoo bar\n\n\nbaz\n",
			"\nfoo\nbar\n\n\nbaz\n",
			4,
		},
		// Complete example:
		{
			" This is a list: \n\n\t* foo\n\t* bar\n\n\n\t* baz  \nBAM    ",
			" This\nis a\nlist: \n\n\t* foo\n\t* bar\n\n\n\t* baz\nBAM",
			6,
		},
		// Multi-byte characters
		{
			strings.Repeat("\u2584 ", 4),
			"\u2584 \u2584" + "\n" +
				strings.Repeat("\u2584 ", 2),
			4,
		},
	}

	for i, tc := range cases {
		actual := blox.WrapString(tc.Input, tc.Lim)
		if actual != tc.Output {
			t.Fatalf("Case %d Input:\n\n`%s`\n\nExpected Output:\n\n`%s`\n\nActual Output:\n\n`%s`", i, tc.Input, tc.Output, actual)
		}
	}
}

func ExampleWrapString() {
	s := blox.WrapString(blox.WithoutLineBreaks(loremIpsum), 50)
	// WrapString will not add a final newline.
	fmt.Println(s)
	// Output:
	// Lorem ipsum dolor sit amet consectetur adipiscing
	// elit torquent ante tortor duiaugue, dictumst
	// convallis eget tempor pharetra lectus magnis
	// lacinia lacus eunostra. Sagittis dolor mattis
	// laoreet justo mollis est varius etiam nisl,
	// siteleifend nullam magna aptent erat vitae. Nullam
	// suspendisse quis volutpat luctusnon a cursus dui
	// urna, facilisis ipsum dapibus etiam odio lacus
	// feugiat neque.Primis pharetra cursus ultrices vel
	// curabitur duis taciti semper, tortor nislurna
	// turpis mauris maecenas ac diam, posuere morbi mi
	// class tincidunt cumsuspendisse.
}

func TestWithoutLineBreaks(t *testing.T) {
	s := blox.WithoutLineBreaks(loremIpsum)
	l := blox.LineCount(s)
	assert.Equal(t, 1, l, "should only render 1 line")
}

func TestReplaceLineBreaks(t *testing.T) {
	l := "String\r\nwith line-breaks\nof various\vkinds.\n"
	s := blox.ReplaceLineBreaks(l, " ")
	c := blox.LineCount(s)
	assert.Equal(t, 1, c)
	assert.Equal(t, 42, utf8.RuneCountInString(s))
}

func ExampleReplaceLineBreaks() {
	l := "String\r\nwith line-breaks\nof various\vkinds.\n"
	s := blox.ReplaceLineBreaks(l, " ")
	fmt.Println(s)
	// Output:
	// String with line-breaks of various kinds.
}
