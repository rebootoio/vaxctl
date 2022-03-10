package helpers

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lensesio/tableprinter"
)

func PrintTable(data interface{}) {
	printer := tableprinter.New(os.Stdout)

	printer.NumbersAlignment = tableprinter.AlignLeft
	printer.HeaderLine = false
	printer.Print(data)
}

func PrintTableWithBorders(data [][]table.Row) {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	for _, rows := range data {
		t.AppendRows(rows)
		t.AppendSeparator()
	}
	t.Render()
}
