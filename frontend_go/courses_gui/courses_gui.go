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

type EditForm struct {
	Widget  *fltk.Grid
	Entries map[string]fltk.Widget
}

func NewEditForm(lines int) EditForm {
	editPanel := fltk.NewGrid(0, 0, 300, 0)
	editPanel.SetLayout(lines, 2, 3, 3)
	editPanel.SetBox(fltk.BORDER_FRAME)
	//	editPanel.End()
	return EditForm{editPanel, map[string]fltk.Widget{}}
}

func (editform EditForm) AddTextline(name string, label string) {
	entry := fltk.NewButton(0, 0, 0, 0)
	entry.SetAlign(fltk.ALIGN_LEFT | fltk.ALIGN_INSIDE)
	entry.SetBox(fltk.PLASTIC_DOWN_BOX)
	entry.SetColor(fltk.ColorFromRgb(255, 255, 200))
	entry.SetCallback(func() { fmt.Println("Clicked") })
	editform.Widget.SetWidget(entry, 0, 1, fltk.GridFill)
	editform.Widget.SetWidget(
		fltk.NewBox(fltk.NO_BOX, 0, 0, 0, 0, label),
		0, 0, fltk.GridFill)
	editform.Entries[name] = entry
}

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
	tableTotals.SetColor(fltk.ColorFromRgb(255, 255, 200))
	//tableTotals.SetCallback(func() { fmt.Println("Hello there!") })
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
	table.DisableRowHeaders()
	table.SetColumnCount(ncols)
	table.SetRowCount(nrows)

	table.SetColumnWidthAll(200)
	table.SetRowHeightAll(50)
	table.SetColumnHeaderHeight(40)
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

	// The editing form (right panel)
	editPanel := NewEditForm(6)
	editPanel.AddTextline("CourseBlock", "Course Block")

	/*
		editPanel := fltk.NewGrid(0, 0, 300, 0, "Edit Panel")
		editPanel.SetLayout(6, 2, 3, 3)
		editPanel.SetBox(fltk.BORDER_FRAME)
		//editPanel.SetColor(fltk.WHITE)

		courseBlock := fltk.NewButton(0, 0, 0, 0, "A course")
		courseBlock.SetAlign(fltk.ALIGN_LEFT | fltk.ALIGN_INSIDE)
		courseBlock.SetBox(fltk.PLASTIC_DOWN_BOX)
		courseBlock.SetColor(fltk.ColorFromRgb(255, 255, 200))
		courseBlock.SetCallback(func() { fmt.Println("Clicked") })
		editPanel.SetWidget(courseBlock, 0, 1, fltk.GridFill)
		editPanel.SetWidget(
			fltk.NewBox(fltk.NO_BOX, 0, 0, 0, 0, "Course Block"),
			0, 0, fltk.GridFill)
		editPanel.End()
	*/
	editPanel.Widget.End()
	//editPanel.Entries["CourseBlock"].(*fltk.Button).Deactivate()
	editPanel.Entries["CourseBlock"].(*fltk.Button).SetLabel("This entry is really rather long")

	courses_box.SetWidget(editPanel.Widget, 1, 1, fltk.GridFill)

	//b2 := fltk.NewBox(fltk.BORDER_BOX, 0, 0, 300, 0, "Right Panel")
	//courses_box.SetWidget(b2, 1, 1, fltk.GridFill)

	courses_box.SetColumnWeight(1, 0)
	courses_box.SetRowWeight(0, 0)

	courses_box.End()

	//

	win.Resizable(courses_box)

	win.End()
	win.Show()
	fltk.Run()
}
