package ui

import (
	"unsafe"
)

// #include <stdlib.h>
// #include "ui.h"
// static uiTableTextColumnOptionalParams *newTextColumnOptionalParams(void)
// {
// 	uiTableTextColumnOptionalParams *p;
//
// 	p = (uiTableTextColumnOptionalParams *) malloc(sizeof (uiTableTextColumnOptionalParams));
// 	return p;
// }
// static void freeTextColumnOptionalParams(uiTableTextColumnOptionalParams *p)
// {
// 	free(p);
// }
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

//
type TableModel struct {
	handler TableModelHandler
	// refCount holds number of Table controls currently using this model
	refCount int

	m *C.uiTableModel
	// TODO: in libui, could add an accessor to get handler from model...
	// (all uiTableModel implementations hold a pointer to the handler anyway)
	h *C.uiTableModelHandler
}

func NewTableModel(handler TableModelHandler) *TableModel {
	// leave underlying C objects nil until first incRef
	m := &TableModel{
		handler:  handler,
		refCount: 0,
		m:        nil,
		h:        nil,
	}

	return m
}

func (m *TableModel) incRef() {
	if m.refCount == 0 {
		// first table attaching - time to create the actual libui C objects
		m.h = registerTableModeHandler(m.handler)
		m.m = C.uiNewTableModel(m.h)
		tablemodels[m.m] = m
	}
	m.refCount++
}

func (m *TableModel) decRef() {
	m.refCount--
	if m.refCount == 0 {
		// last table using this model - free up the libui C objects
		delete(tablemodels, m.m)
		C.uiFreeTableModel(m.m)
		unregisterTableModelHandler(m.h)
		m.m = nil
		m.h = nil
	}
}

func (m *TableModel) RowInserted(newIndex int) {
	if m.m != nil {
		C.uiTableModelRowInserted(m.m, C.int(newIndex))
	}
}

func (m *TableModel) RowChanged(index int) {
	if m.m != nil {
		C.uiTableModelRowChanged(m.m, C.int(index))
	}
}

func (m *TableModel) RowDeleted(oldIndex int) {
	if m.m != nil {
		C.uiTableModelRowDeleted(m.m, C.int(oldIndex))
	}
}

// ------------

type TextColumnOptionalParams struct {
	ColorModelColumn int
}

func (p *TextColumnOptionalParams) toC() *C.uiTableTextColumnOptionalParams {
	if p == nil {
		return nil // a bit icky accepting nil receiver, but valid!
	}
	cp := C.newTextColumnOptionalParams()
	cp.ColorModelColumn = C.int(p.ColorModelColumn)
	return cp
}

// -------------
// Table stuff

// no need to lock this; only the GUI thread can access it
var tables = make(map[*C.uiTable]*Table)

// Table is... TODO
type Table struct {
	model *TableModel
	c     *C.uiControl
	t     *C.uiTable

	onSelectionChanged func(*Table)
}

type TableParams struct {
	Model                         *TableModel
	RowBackgroundColorModelColumn int
	MultiSelect                   bool
}

// NewTable creates a new Table control
func NewTable(params *TableParams) *Table {

	params.Model.incRef()

	cp := C.uiTableParams{}
	cp.Model = params.Model.m
	cp.RowBackgroundColorModelColumn = C.int(params.RowBackgroundColorModelColumn)
	if params.MultiSelect {
		cp.MultiSelect = 1
	} else {
		cp.MultiSelect = 0
	}

	t := new(Table)
	t.model = params.Model

	t.t = C.uiNewTable((*C.uiTableParams)(unsafe.Pointer(&cp)))
	t.c = (*C.uiControl)(unsafe.Pointer(t.t))

	C.realuiTableOnSelectionChanged(t.t)
	tables[t.t] = t

	return t
}

// Destroy destroys the Table.
func (t *Table) Destroy() {
	delete(tables, t.t)
	C.uiControlDestroy(t.c)

	t.model.decRef()
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

func (t *Table) AppendTextColumn(name string, modelColumn int, editableModelColumn int, params *TextColumnOptionalParams) {
	tmpName := C.CString(name)
	cp := params.toC()
	C.uiTableAppendTextColumn(t.t, tmpName, C.int(modelColumn), C.int(editableModelColumn), cp)
	if cp != nil {
		C.freeTextColumnOptionalParams(cp)
	}
	freestr(tmpName)
}
