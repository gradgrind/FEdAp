#include "courses_gui.h"
#include <FL/Fl_Choice.H>
#include <FL/Fl_Output.H>
#include <FL/fl_draw.H>

#include <iostream>

CoursesGui::CoursesGui()
    : Fl_Flex(0, 0, 0, 0)
{
    gap(3);
    color(FL_YELLOW);

    // *** Top Panel ***

    // Top Panel – whatever is still needed ...
    auto todo = new Fl_Box(FL_FLAT_BOX, 0, 0, 0, 0, "TODO");
    fixed(todo, 50);

    // Top Panel – the selectors and totals info, at the panel bottom
    auto panelBox = new Fl_Flex(0, 0, 0, 0);
    panelBox->type(Fl_Flex::ROW);
    panelBox->margin(5, 5);
    fixed(panelBox, 40);

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

    // Top Panel – end

    // *** The course list/table and course editor ***
    auto mainview = new Fl_Flex(0, 0, 0, 0);
    mainview->type(Fl_Flex::ROW);
    mainview->gap(3);
    //mainview->margin(5, 5);

    auto table = new CourseTable();
    Widgets["Table"] = table;

    // The course editor form
    //TODO ...
    auto todo2 = new Fl_Box(FL_FLAT_BOX, 0, 0, 0, 0, "TODO");

    // End of mainview
    mainview->fixed(todo2, 300);
    mainview->end();

    // End of CoursesGui
    end();
}

void CourseTable::_row_cb(
    Fl_Widget* w, void* a)
{
    ((CourseTable*) w)->row_cb(a);
}

void CourseTable::row_cb(
    void* a)
{
    int r = callback_row();
    std::cout << "ROW: " << r << std::endl;

    select_row(r, 1);
}

//TODO: Adjust cols to fit space, but with a set minimum size
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
            {"Ma", "11.A", "AH", "k11", "1, 1", ""},
            {"Fr", "11.B", "AH", "k11", "1, 1", ""}};

    callback(_row_cb);
    when(FL_WHEN_RELEASE_ALWAYS);
}

void CourseTable::draw_cell(
    TableContext context, int ROW, int COL, int X, int Y, int W, int H)
{
    static char s[40];
    switch (context) {
    case CONTEXT_STARTPAGE: // before page is drawn..
        //fl_font(FL_HELVETICA, 16); // set the font for our drawing operations

        std::cout << "tow: " << tow << std::endl;
        std::cout << "tiw: " << tiw << std::endl;
        return;
    case CONTEXT_COL_HEADER: // Draw column headers
        fl_push_clip(X, Y, W, H);
        fl_draw_box(FL_THIN_UP_BOX, X, Y, W, H, row_header_color());
        fl_color(FL_BLACK);
        fl_draw(headers[COL].c_str(), X, Y, W, H, FL_ALIGN_CENTER);
        fl_pop_clip();
        return;
    //case CONTEXT_ROW_HEADER:      // Draw row headers
    //    return;
    case CONTEXT_CELL: // Draw data in cells
        fl_push_clip(X, Y, W, H);
        // Draw cell bg
        if (row_selected(ROW)) {
            fl_color(FL_YELLOW);
        } else {
            fl_color(FL_WHITE);
        }
        fl_rectf(X, Y, W, H);
        // Draw cell data
        fl_color(FL_GRAY0);
        fl_draw(data[ROW][COL].c_str(), X, Y, W, H, FL_ALIGN_CENTER);
        // Draw box border
        fl_color(color());
        fl_rect(X, Y, W, H);
        fl_pop_clip();
        return;
    default:
        return;
    }
}

void CourseTable::size_columns()
{
    int ncols = cols();
    int nrows = rows();
    for (int c = 0; c < ncols; ++c) {
        //int w = headers[c];
        for (int r = 0; r < nrows; ++r) {
        }
    }
}
