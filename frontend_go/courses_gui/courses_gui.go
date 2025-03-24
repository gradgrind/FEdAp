package courses_gui

import (
	"fmt"

	"github.com/pwiecz/go-fltk"
)

const (
	nrows = 5
	ncols = 6
)

type pair struct {
	pos  int
	size int
}

var (
	rows          []pair
	cols          []pair
	CourseWidgets map[string]fltk.Widget
)

const (
	w0 = 1000
	h0 = 700
)

func Ui() {
	CourseWidgets = map[string]fltk.Widget{}
	win := fltk.NewWindow(w0, h0)

	//

	courses_box := fltk.NewGrid(0, 0, win.W(), win.H())
	courses_box.SetLayout(2, 2, 3, 3)
	courses_box.SetColor(fltk.WHITE)

	topPanel := fltk.NewFlex(0, 0, 0, 100)
	topPanel.SetType(fltk.COLUMN)
	fltk.NewBox(fltk.FLAT_BOX, 0, 0, 0, 0, "TODO")
	panelBox := fltk.NewFlex(0, 0, 0, 0)
	panelBox.SetType(fltk.ROW)
	panelBox.SetMargin(5, 5)

	tableType := fltk.NewChoice(0, 0, 0, 0)
	tableType.ClearVisibleFocus()
	panelBox.Fixed(tableType, 100)
	tableType.Add("Class", nil)
	tableType.Add("Teacher", nil)
	tableType.Add("Subject", nil)
	CourseWidgets["TableType"] = tableType

	tableRow := fltk.NewChoice(0, 0, 0, 0)
	tableRow.ClearVisibleFocus()
	panelBox.Fixed(tableRow, 150)
	tableRow.Add("FB: Fritz Jolander Jeremias Braun", nil)
	tableRow.Add("DG: Diego Garcia", nil)
	tableRow.Add("PM: Pamela Masterson", nil)
	CourseWidgets["TableRow"] = tableRow

	// The label of the output widget allocates no space of its own,
	// overwriting anything to the left of the widget, so add padding.
	pad := fltk.NewBox(fltk.NO_BOX, 0, 0, 0, 0)
	tableTotals := fltk.NewOutput(0, 0, 0, 0)
	tableTotals.SetValue("Read only ˝Öößŋħĸ€")
	tableTotals.ClearVisibleFocus() // no cursor, but text cannot be copied
	tableTotals.SetLabel("Total lessons:")
	w, _ := tableTotals.MeasureLabel()
	panelBox.Fixed(pad, w+20)
	CourseWidgets["TableTotals"] = tableTotals

	panelBox.End()
	topPanel.Fixed(panelBox, 40)

	topPanel.End()
	courses_box.SetWidgetWithSpan(topPanel, 0, 0, 1, 2, fltk.GridFill)

	table := fltk.NewTableRow(0, 0, 0, 0)
	table.EnableColumnHeaders()
	table.EnableRowHeaders()
	table.SetColumnCount(ncols)
	table.SetRowCount(nrows)

	table.SetColumnWidthAll(200)
	table.SetRowHeightAll(50)
	table.SetColumnHeaderHeight(40)
	table.DisableRowHeaders()
	table.SetType(fltk.SelectSingle)

	table.SetDrawCellCallback(func(tc fltk.TableContext, row, col, x, y, w, h int) {
		if tc == fltk.ContextCell {
			fltk.SetDrawFont(fltk.HELVETICA, 14)
			fltk.DrawBox(fltk.FLAT_BOX, x, y, w, h, fltk.BLACK)
			bg := fltk.WHITE
			if table.IsRowSelected(row) {
				bg = fltk.YELLOW
			}
			fltk.DrawBox(fltk.FLAT_BOX, x+1, y+1, w-2, h-2, bg)
			fltk.SetDrawColor(fltk.BLACK)
			fltk.Draw(fmt.Sprintf("%d", row+col), x, y, w, h, fltk.ALIGN_CENTER)
		}
		if tc == fltk.ContextRowHeader {
			fltk.SetDrawFont(fltk.HELVETICA_BOLD, 14)
			fltk.DrawBox(fltk.UP_BOX, x, y, w, h, fltk.BACKGROUND_COLOR)
			fltk.SetDrawColor(fltk.BLACK)
			fltk.Draw(fmt.Sprintf("row %d", row+1), x, y, w, h, fltk.ALIGN_CENTER)
		}
		if tc == fltk.ContextColHeader {
			fltk.SetDrawFont(fltk.HELVETICA_BOLD, 14)
			fltk.DrawBox(fltk.UP_BOX, x, y, w, h, fltk.BACKGROUND_COLOR)
			fltk.SetDrawColor(fltk.BLACK)
			fltk.Draw(fmt.Sprintf("col %d", col+1), x, y, w, h, fltk.ALIGN_CENTER)
		}
	})

	table.End()
	//		win.Resizable(table)

	courses_box.SetWidget(table, 1, 0, fltk.GridFill)

	b2 := fltk.NewBox(fltk.BORDER_BOX, 0, 0, 300, 0, "Right Panel")
	courses_box.SetWidget(b2, 1, 1, fltk.GridFill)

	courses_box.SetColumnWeight(1, 0)
	courses_box.SetRowWeight(0, 0)

	courses_box.End()

	//

	win.Resizable(courses_box)

	win.End()
	win.Show()
	fltk.Run()
}
