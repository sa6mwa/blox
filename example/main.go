package main

import (
	_ "embed"

	"github.com/sa6mwa/blox"
)

//go:embed block1.txt
var block1 string

//go:embed block2.txt
var block2 string

//go:embed block3.txt
var loremIpsum string

func main() {
	b := blox.New().SetColumnsAndRows(80, 24).
		SetTrimRightSpaces(true).SetTrimFinalEmptyLines(true) // or just SetTrim(true)

	b.Move(4, 10).PutText(block1).
		Move(0, 0).PutText(loremIpsum).Move(1, 1).PushPos().Move(2, 2).PushPos().
		MoveY(10).SetLineSpacing(1).PutText(loremIpsum).SetLineSpacing(1).
		Move(11, 3).PutText(block2).PopPos().PutText("ÅÄÖÅÄÖ").PopPos().PutText("OKILI").
		Move(58, 2).PutText(block2).PutText(block1).Move(11, 5).DrawVerticalLine(9, 20).
		Move(15, 18).PutText(block1).PushPos().
		DrawSeparator().DrawSeparator('_').PopPos().MoveRight(30).
		DrawSplit()

	b.PrintCanvas() // or: fmt.Print(b.String())
}
