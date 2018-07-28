package ui

// #include <stdlib.h>
// #include "ui.h"
//
// //TODO: we use this typedef to cast our (go) callback with a non-const uiTableValue
// // into what the C side expects.
// // Would be much better to just typedef a const version of uiTableValue and define
// // our go function using it, but that seems to confuse cgo about other uses
// // of uiTableValue.
// // Could be a cgo bug? Try again some time.
// typedef void (*setCellValueFn)(uiTableModelHandler *, uiTableModel *, int, int, const uiTableValue *);
//
// extern int doTableModelHandlerNumColumns(uiTableModelHandler *, uiTableModel *);
// extern uiTableValueType doTableModelHandlerColumnType(uiTableModelHandler *, uiTableModel *, int);
// extern int doTableModelHandlerNumRows(uiTableModelHandler *, uiTableModel *);
// extern uiTableValue* doTableModelHandlerCellValue(uiTableModelHandler *, uiTableModel *, int, int);
// extern void doTableModelHandlerSetCellValue(uiTableModelHandler *, uiTableModel *, int, int, uiTableValue *);
//
// static inline uiTableModelHandler *allocTableModelHandler(void)
// {
// 	uiTableModelHandler *h;
//
// 	h = (uiTableModelHandler *) malloc(sizeof (uiTableModelHandler));
// 	if (h == NULL)		// TODO
// 		return NULL;
//  h->NumColumns = doTableModelHandlerNumColumns;
//  h->ColumnType = doTableModelHandlerColumnType;
//  h->NumRows = doTableModelHandlerNumRows;
//  h->CellValue = doTableModelHandlerCellValue;
//  // note cast, to handle go version with non-const uiTableValue param
//  h->SetCellValue = (setCellValueFn)doTableModelHandlerSetCellValue;
// 	return h;
// }
// static inline void freeTableModelHandler(uiTableModelHandler *h)
// {
// 	free(h);
// }
import "C"

// -------------

// Note: these must be numerically identical to their libui equivalents.
type TableValueType uint

const (
	StringColumn TableValueType = iota
	ImageColumn
	IntColumn
	ColorColumn
)

// -------------
// TableModelHandler

// no need to lock this; only the GUI thread can access it
var tablemodelhandlers = make(map[*C.uiTableModelHandler]TableModelHandler)

type TableModelHandler interface {
	NumColumns(*TableModel) int
	ColumnType(*TableModel, int) TableValueType
	NumRows(*TableModel) int
	CellValue(*TableModel, int, int) interface{}
	SetCellValue(*TableModel, int, int, interface{})
}

func registerTableModeHandler(h TableModelHandler) *C.uiTableModelHandler {
	uh := C.allocTableModelHandler()

	tablemodelhandlers[uh] = h
	return uh
}

func unregisterTableModelHandler(uh *C.uiTableModelHandler) {
	delete(tablemodelhandlers, uh)
	C.freeTableModelHandler(uh)
}

//export doTableModelHandlerNumColumns
func doTableModelHandlerNumColumns(uh *C.uiTableModelHandler, um *C.uiTableModel) C.int {
	h := tablemodelhandlers[uh]
	m := tablemodels[um]
	return C.int(h.NumColumns(m))
}

//export doTableModelHandlerColumnType
func doTableModelHandlerColumnType(uh *C.uiTableModelHandler, um *C.uiTableModel, col C.int) C.uiTableValueType {
	h := tablemodelhandlers[uh]
	m := tablemodels[um]
	typ := h.ColumnType(m, int(col))
	return C.uiTableValueType(typ)
}

//export doTableModelHandlerNumRows
func doTableModelHandlerNumRows(uh *C.uiTableModelHandler, um *C.uiTableModel) C.int {
	h := tablemodelhandlers[uh]
	m := tablemodels[um]
	return C.int(h.NumRows(m))
}

//export doTableModelHandlerCellValue
func doTableModelHandlerCellValue(uh *C.uiTableModelHandler, um *C.uiTableModel, row C.int, col C.int) *C.uiTableValue {
	h := tablemodelhandlers[uh]
	m := tablemodels[um]

	raw := h.CellValue(m, int(row), int(col))

	switch v := raw.(type) {
	case string:
		tmpstr := C.CString(v)
		uv := C.uiNewTableValueString(tmpstr)
		freestr(tmpstr)
		return uv
	case int:
		return C.uiNewTableValueInt(C.int(v))
	// TODO - support other numeric types?
	//	case color:
	//	case image:
	default:
		panic("unsupported type")
	}
}

//export doTableModelHandlerSetCellValue
func doTableModelHandlerSetCellValue(uh *C.uiTableModelHandler, um *C.uiTableModel, row C.int, col C.int, value *C.uiTableValue) {
	// TODO!!!
}
