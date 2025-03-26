#include "courses_gui.h"
#include <FL/Fl_Choice.H>
#include <FL/Fl_Output.H>
#include <FL/fl_draw.H>

CoursesGui::CoursesGui()
    : Fl_Grid(0, 0, 0, 0)
{
    layout(2, 2, 3, 3);
    row_weight(0, 0);
    color(FL_YELLOW);

    // Top Panel
    auto topPanel = new Fl_Flex(0, 0, 0, 0);

    // Top Panel – whatever is still needed ...
    auto todo = new Fl_Box(FL_FLAT_BOX, 0, 0, 0, 0, "TODO");
    topPanel->fixed(todo, 100);

    // Top Panel – the selectors and totals info, at the panel bottom
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

    // Top Panel – end
    topPanel->end();
    widget(topPanel, 0, 0, 1, 2, FL_GRID_FILL);

    // The course list/table
    auto table = new CourseTable();
    Widgets["Table"] = table;
    widget(table, 1, 0, FL_GRID_FILL);

    // The course editor form
    //TODO ...
    auto todo2 = new Fl_Box(FL_FLAT_BOX, 0, 0, 0, 0, "TODO");
    widget(todo2, 1, 1, FL_GRID_FILL);

    end();
}

CourseTable::CourseTable()
    : Fl_Table_Row(0, 0, 0, 0)
{
    // Rows
    rows(5);            // how many rows
    row_header(0);      // enable row headers (along left)
    row_height_all(30); // default height of rows
    row_resize(0);      // disable row resizing
    // Cols
    cols(6);               // how many columns
    col_header(1);         // enable column headers (along top)
    col_header_height(30); // enable column headers (along top)
    col_width_all(150);    // default width of columns
    col_resize(0);         // enable column resizing
    type(Fl_Table_Row::SELECT_SINGLE);
    end(); // end the Fl_Table group

    // Temporary (test) data
    headers = {"Subject", "Groups", "Teachers", "Rooms", "Units", "Properties"};
    data = {{"Fr", "11.A", "PM", "fs1", "1, 1, 1", ""},
            {"Fr", "11.B", "PM", "fs1", "1, 1, 1", ""},
            {"Ma", "11", "AH", "k11", "3, 3", "HU"},
            {"Ma", "11.A", "AH", "fs1", "1, 1", ""},
            {"Fr", "11.B", "AH", "fs1", "1, 1", ""}};
}

void CourseTable::DrawHeader(
    const char *s, int X, int Y, int W, int H)
{
    fl_push_clip(X, Y, W, H);
    fl_draw_box(FL_THIN_UP_BOX, X, Y, W, H, row_header_color());
    fl_color(FL_BLACK);
    fl_draw(s, X, Y, W, H, FL_ALIGN_CENTER);
    fl_pop_clip();
}

void CourseTable::DrawData(
    const char *s, int X, int Y, int W, int H)
{
    fl_push_clip(X, Y, W, H);
    // Draw cell bg
    fl_color(FL_WHITE);
    fl_rectf(X, Y, W, H);
    // Draw cell data
    fl_color(FL_GRAY0);
    fl_draw(s, X, Y, W, H, FL_ALIGN_CENTER);
    // Draw box border
    fl_color(color());
    fl_rect(X, Y, W, H);
    fl_pop_clip();
}

void CourseTable::draw_cell(
    TableContext context, int ROW, int COL, int X, int Y, int W, int H)
{
    static char s[40];
    switch (context) {
    case CONTEXT_STARTPAGE: // before page is drawn..
        //fl_font(FL_HELVETICA, 16); // set the font for our drawing operations
        return;
    case CONTEXT_COL_HEADER:         // Draw column headers
        sprintf(s, "%c", 'A' + COL); // "A", "B", "C", etc.
        DrawHeader(headers[COL].c_str(), X, Y, W, H);
        return;
    //case CONTEXT_ROW_HEADER:      // Draw row headers
    //    sprintf(s, "%03d:", ROW); // "001:", "002:", etc
    //    DrawHeader(s, X, Y, W, H);
    //    return;
    case CONTEXT_CELL: // Draw data in cells
        DrawData(data[ROW][COL].c_str(), X, Y, W, H);
        return;
    default:
        return;
    }
}
