package main

import (
	"fmt"
	"github.com/bcampbell/ui"
)

type Dat struct{}

func (d *Dat) NumColumns(m *ui.TableModel) int {
	return 2
}

func (d *Dat) ColumnType(m *ui.TableModel, col int) ui.TableModelColumnType {
	return ui.StringColumn
}

func (d *Dat) NumRows(m *ui.TableModel) int {
	return 1000
}

func (d *Dat) CellValue(m *ui.TableModel, row int, col int) interface{} {
	return fmt.Sprintf("value %d,%d", row, col)
}

func (d *Dat) SetCellValue(*ui.TableModel, int, int, interface{}) {
}

func main() {
	err := ui.Main(func() {
		input := ui.NewEntry()
		button := ui.NewButton("Greet")
		greeting := ui.NewLabel("")
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Enter your name:"), false)
		box.Append(input, false)
		box.Append(button, false)
		box.Append(greeting, false)

		dat := &Dat{}
		model := ui.NewTableModel(dat)
		table := ui.NewTable(model, ui.TableStyleMultiSelect)
		table.AppendTextColumn("one", 0)
		table.AppendTextColumn("two", 1)
		table.OnSelectionChanged(func(t *ui.Table) {
			selected := t.GetSelection()
			fmt.Printf("selected: %v\n", selected)
		})
		box.Append(table, true)

		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			greeting.SetText("Hello, " + input.Text() + "!")
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
