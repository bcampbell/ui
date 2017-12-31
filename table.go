package ui

import (
	"unsafe"
)

// #include "ui.h"
// #include <stdlib.h>
// extern void doTableOnSelectionChanged(uiTable *, void *);
// static inline void realuiTableOnSelectionChanged(uiTable *t)
// {
// 	uiTableOnSelectionChanged(t, doTableOnSelectionChanged, NULL);
// }
import "C"

// -------------
// TableModel stuff

// no need to lock this; only the GUI thread can access it
var tablemodels = make(map[*C.uiTableModel]*TableModel)

type TableModel struct {
	m *C.uiTableModel
	// TODO: in libui, could add an accessor to get handler from model...
	// (all uiTableModel implementations hold a pointer to the handler anyway)
	h *C.uiTableModelHandler
}

func NewTableModel(handler TableModelHandler) *TableModel {
	m := &TableModel{}
	m.h = registerTableModeHandler(handler)
	m.m = C.uiNewTableModel(m.h)
	tablemodels[m.m] = m
	return m
}

func (m *TableModel) Destroy() {
	delete(tablemodels, m.m)
	C.uiFreeTableModel(m.m)
	unregisterTableModelHandler(m.h)
}

func (m *TableModel) RowInserted(newIndex int) {
	C.uiTableModelRowInserted(m.m, C.int(newIndex))
}

func (m *TableModel) RowChanged(index int) {
	C.uiTableModelRowChanged(m.m, C.int(index))
}

func (m *TableModel) RowDeleted(oldIndex int) {
	C.uiTableModelRowDeleted(m.m, C.int(oldIndex))
}

// -------------
// Table stuff

// no need to lock this; only the GUI thread can access it
var tables = make(map[*C.uiTable]*Table)

// Table is... TODO
type Table struct {
	c *C.uiControl
	t *C.uiTable

	onSelectionChanged func(*Table)
}

// NewTable creates a new Table control
// TODO: style flags
func NewTable(m *TableModel) *Table {
	t := new(Table)
	var styleFlags int = 0

	//ctext := C.CString(text)
	t.t = C.uiNewTable(m.m, C.int(styleFlags))
	t.c = (*C.uiControl)(unsafe.Pointer(t.t))
	//freestr(ctext)

	C.realuiTableOnSelectionChanged(t.t)
	tables[t.t] = t

	return t
}

// Destroy destroys the Button.
func (t *Table) Destroy() {
	delete(tables, t.t)
	C.uiControlDestroy(t.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Table. This is only used by package ui itself and should
// not be called by programs.
func (t *Table) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(t.c))
}

// Handle returns the OS-level handle associated with this Table.
func (t *Table) Handle() uintptr {
	return uintptr(C.uiControlHandle(t.c))
}

// Show shows the Table control.
func (t *Table) Show() {
	C.uiControlShow(t.c)
}

// Hide hides the Table control.
func (t *Table) Hide() {
	C.uiControlHide(t.c)
}

// Enable enables the Table control.
func (t *Table) Enable() {
	C.uiControlEnable(t.c)
}

// Disable disables the Table control.
func (t *Table) Disable() {
	C.uiControlDisable(t.c)
}

// OnSelectionChanged registers f to be run when the set of selected
// items in the table changes.
// Only one function can be registered at a time.
func (t *Table) OnSelectionChanged(f func(*Table)) {
	t.onSelectionChanged = f
}

//export doTableOnSelectionChanged
func doTableOnSelectionChanged(tt *C.uiTable, data unsafe.Pointer) {
	t := tables[tt]
	if t.onSelectionChanged != nil {
		t.onSelectionChanged(t)
	}
}

func (t *Table) AppendTextColumn(name string, modelColumn int) {
	tmpName := C.CString(name)
	C.uiTableAppendTextColumn(t.t, tmpName, C.int(modelColumn))
	freestr(tmpName)
}
