package main

import (
	//	"fmt"
	"github.com/bcampbell/ui"
	"github.com/icrowley/fake"
	"math/rand"
	"strconv"
)

type Person struct {
	FirstName string
	LastName  string
	ShoeSize  int
}

type PersonDB struct {
	People []Person
}

// implement the TableModelHandler interface

func (db *PersonDB) NumColumns(m *ui.TableModel) int {
	return 3
}

func (db *PersonDB) ColumnType(m *ui.TableModel, col int) ui.TableValueType {
	return ui.StringColumn
}

func (db *PersonDB) NumRows(m *ui.TableModel) int {
	//fmt.Printf("numrows: %d\n", len(db.People))
	return len(db.People)
}

func (db *PersonDB) CellValue(m *ui.TableModel, row int, col int) interface{} {
	//fmt.Printf("CellValue %d %d\n", row, col)
	prod := &db.People[row]
	switch col {
	case 0:
		return prod.FirstName
	case 1:
		return prod.LastName
	case 2:
		return strconv.Itoa(prod.ShoeSize)
	}
	return nil
}

func (db *PersonDB) SetCellValue(*ui.TableModel, int, int, interface{}) {
	// TODO
}

func RandomPerson() Person {
	p := Person{}
	p.FirstName = fake.FirstName()
	p.LastName = fake.LastName()
	p.ShoeSize = 8 + rand.Intn(10)
	return p
}

type App struct {
	db PersonDB

	table          *ui.Table
	model          *ui.TableModel
	firstNameInput *ui.Entry
	lastNameInput  *ui.Entry
	shoeSizeInput  *ui.Spinbox
	createButton   *ui.Button
	//deleteButton    *ui.Button
	//selSummaryLabel *ui.Label
}

func (app *App) Init() {
	for i := 0; i < 10; i++ {
		app.db.People = append(app.db.People, RandomPerson())
	}
}

func (app *App) rethink() {
	invalid := app.firstNameInput.Text() == "" || app.lastNameInput.Text() == ""

	if invalid {
		app.createButton.Disable()
	} else {
		app.createButton.Enable()
	}
}

func (app *App) buildGUI() ui.Control {

	app.firstNameInput = ui.NewEntry()
	app.lastNameInput = ui.NewEntry()
	app.shoeSizeInput = ui.NewSpinbox(8, 18)

	vbox := ui.NewVerticalBox()

	// table display

	app.model = ui.NewTableModel(&app.db)
	table := ui.NewTable(&ui.TableParams{app.model, -1})
	table.AppendTextColumn("FirstName", 0, -1, nil)
	table.AppendTextColumn("LastName", 1, -1, nil)
	table.AppendTextColumn("ShoeSize", 2, -1, nil)
	/*
		table.OnSelectionChanged(func(t *ui.Table) {
			app.HandleSelectionChanged()
		})
	*/
	vbox.Append(table, true)
	app.table = table

	// data entry
	grp := ui.NewGroup("Enter person details")
	{
		// Hmm... think we really want a form here, but ui doesn't seem to support it yet
		box := ui.NewVerticalBox()
		box.SetPadded(true)
		box.Append(ui.NewLabel("First name"), false)
		box.Append(app.firstNameInput, false)
		box.Append(ui.NewLabel("Last name"), false)
		box.Append(app.lastNameInput, false)
		box.Append(ui.NewLabel("Shoe size"), false)
		box.Append(app.shoeSizeInput, false)

		app.createButton = ui.NewButton("Create")
		box.Append(app.createButton, false)
		grp.SetChild(box)
	}

	vbox.Append(grp, false)

	//	app.selSummaryLabel = ui.NewLabel("")
	//	vbox.Append(app.selSummaryLabel, false)
	//	app.deleteButton = ui.NewButton("Delete")
	//	vbox.Append(app.deleteButton, false)

	app.firstNameInput.OnChanged(func(e *ui.Entry) { app.rethink() })
	app.lastNameInput.OnChanged(func(e *ui.Entry) { app.rethink() })

	app.createButton.OnClicked(func(b *ui.Button) {
		p := Person{
			FirstName: app.firstNameInput.Text(),
			LastName:  app.lastNameInput.Text(),
			ShoeSize:  app.shoeSizeInput.Value(),
		}
		app.db.People = append(app.db.People, p)
		app.model.RowInserted(len(app.db.People) - 1)
		app.firstNameInput.SetText("")
		app.lastNameInput.SetText("")
		app.rethink()
	})
	//	app.deleteButton.OnClicked(func(b *ui.Button) { app.DeleteSelected() })
	app.rethink()
	return vbox
}

/*
func (app *App) HandleSelectionChanged() {

	sel := app.table.GetSelection()

	summary := fmt.Sprintf("%d selected", len(sel))
	app.selSummaryLabel.SetText(summary)

	fmt.Printf("selected: %v\n", sel)
	if len(sel) > 0 {
		app.deleteButton.Enable()
	} else {
		app.deleteButton.Disable()
	}

}

func (app *App) DeleteSelected() {
	sel := app.table.GetSelection()
	// remove highest-first so we don't screw up our indices
	sort.Sort(sort.Reverse(sort.IntSlice(sel)))
	for _, idx := range sel {
		app.db.People = append(app.db.People[:idx], app.db.People[idx+1:]...)
		app.model.RowDeleted(idx)
	}
	app.HandleSelectionChanged()
}
*/
func main() {

	err := ui.Main(func() {
		app := &App{}

		app.Init()

		gui := app.buildGUI()

		window := ui.NewWindow("Hello", 200, 100, false)
		window.SetMargined(true)
		window.SetChild(gui)
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
