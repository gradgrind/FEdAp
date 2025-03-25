#include "courses_gui.h"
#include <FL/Fl_Choice.H>
#include <FL/Fl_Output.H>

CoursesGui::CoursesGui()
    : Fl_Grid(0, 0, 0, 0)
{
    layout(2, 2, 3, 3);
    color(FL_YELLOW);

    auto topPanel = new Fl_Flex(0, 0, 0, 100); // 100 -> 0?
    auto todo = new Fl_Box(FL_FLAT_BOX, 0, 0, 0, 0, "TODO");
    auto panelBox = new Fl_Flex(0, 0, 0, 0);
    panelBox->type(Fl_Flex::ROW);
    panelBox->margin(5, 5);
    auto tableType = new Fl_Choice(0, 0, 0, 0);
    tableType->clear_visible_focus();
    panelBox->fixed(tableType, 100);
    tableType->add("Class");
    tableType->add("Teacher");
    tableType->add("Subject");
    Widgets["TableType"] = tableType;

    auto tableRow = new Fl_Choice(0, 0, 0, 0);
    tableRow->clear_visible_focus();
    panelBox->fixed(tableRow, 150);
    tableRow->add("FB: Fritz Jolander Jeremias Braun");
    tableRow->add("DG: Diego Garcia");
    tableRow->add("PM: Pamela Masterson");
    Widgets["TableRow"] = tableRow;

    // The label of the output widget allocates no space of its own,
    // overwriting anything to the left of the widget, so add padding.
    auto pad = new Fl_Box(FL_NO_BOX, 0, 0, 0, 0);
    auto tableTotals = new Fl_Output(0, 0, 0, 0);
    tableTotals->color(0xFFFFC800);
    //tableTotals.SetCallback(func() { fmt.Println("Hello there!") })
    tableTotals->value("Read only ˝Öößŋħĸ€");
    tableTotals->clear_visible_focus(); // no cursor, but text cannot be copied
    tableTotals->label("Total lessons:");
    int wl, hl;
    tableTotals->measure_label(wl, hl);
    panelBox->fixed(pad, wl + 20);
    Widgets["TableTotals"] = tableTotals;
    panelBox->end();
    topPanel->fixed(panelBox, 40);
    topPanel->end();
    widget(topPanel, 0, 0, 1, 2, FL_GRID_FILL);
}
