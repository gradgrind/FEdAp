#include "courses_gui.h"
#include "editform.h"
#include <FL/Fl.H>
#include <FL/Fl_Choice.H>
#include <FL/Fl_Output.H>
#include <FL/fl_draw.H>

#include <algorithm>
#include <iostream>

CoursesGui::CoursesGui()
    : Fl_Flex(Fl_Flex::COLUMN)
{
    gap(3);
    color(FL_YELLOW);

    // *** Top Panel ***

    // Top Panel – whatever is still needed ...
    auto todo = new Fl_Box(FL_FLAT_BOX, 0, 0, 0, 0, "TODO");
    fixed(todo, 50);

    // Top Panel – the selectors and totals info, at the panel bottom
    auto panelBox = new Fl_Flex(Fl_Flex::ROW);
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
    auto mainview = new Fl_Flex(Fl_Flex::ROW);
    mainview->gap(3);
    //mainview->margin(5, 5);

    auto table = new CourseTable();
    Widgets["Table"] = table;

    // *** The course editor form ***
    //TODO ...
    auto editpanel = new EditForm();
    current(NULL);
    editpanel->add_value("CourseBlock", "Course Block");   //->value("");
    editpanel->add_value("BlockSubject", "Block Subject"); //->deactivate();
    //editpanel->add_value("Block Subject")->hide();

    editpanel->add_separator();
    editpanel->add_value("Subject", "Subject");
    editpanel->add_value("Teachers", "Teachers");
    editpanel->add_value("Rooms", "Rooms");
    editpanel->add_value("Units", "Units");
    editpanel->add_list("Constraints", "Constraints");

    // "Properties"

    editpanel->do_layout();

    editpanel->entries[1].widget->deactivate(); //->hide();

    // End of course editor form

    // End of mainview
    mainview->fixed(editpanel, 300);
    mainview->end();

    // End of CoursesGui
    end();
}

void CourseTable::_row_cb(
    void* table)
{
    //TODO
    std::cout << "§§§ " << ((CourseTable*) table)->_current_row << std::endl;
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
    //col_width_all(150);    // default width of columns
    col_resize(0); // disable column resizing
    col_header_color(fl_rgb_color(230, 230, 255));
    type(Fl_Table_Row::SELECT_SINGLE);
    end(); // end the Fl_Table group

    // Temporary (test) data
    headers = {"Subject", "Groups", "Teachers", "Rooms", "Units", "Properties"};
    data = {{"Fr", "11.A", "PM", "fs1", "1, 1, 1", ""},
            {"Fr", "11.B", "PM", "fs1", "1, 1, 1", ""},
            {"Ma", "11", "AH", "k11", "3, 3", "HU"},
            {"Ma", "11.A", "AH", "k11", "1, 1", ""},
            {"Fr", "11.B", "AH", "k11", "1, 1", ""}};
}

void CourseTable::draw_cell(
    TableContext context, int ROW, int COL, int X, int Y, int W, int H)
{
    static char s[40];
    switch (context) {
    case CONTEXT_STARTPAGE: // before page is drawn..
        //fl_font(FL_HELVETICA, 16); // set the font for our drawing operations

        // Adjust column widths
        size_columns();

        // Handle change of selected row
        if (Fl_Table::select_row != _current_row) {
            select_row(Fl_Table::select_row);
            _current_row = Fl_Table::select_row;
            Fl::add_timeout(0.0, _row_cb, this);
        }
        return;
    case CONTEXT_COL_HEADER: // Draw column headers
        fl_push_clip(X, Y, W, H);
        fl_draw_box(FL_THIN_UP_BOX, X, Y, W, H, col_header_color());
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
            //if (ROW == _current_row) {
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

//TODO: Do I need to set font (fl_font()) before using fl_measure()?
void CourseTable::size_columns()
{
    struct colwidth
    {
        int col, width;
    };
    std::vector<colwidth> colwidths;

    int ncols = cols();
    int nrows = rows();
    int w, h, wmax;
    for (int c = 0; c < ncols; ++c) {
        w = 0;
        fl_measure(headers[c].c_str(), w, h, 0);
        wmax = w;
        for (int r = 0; r < nrows; ++r) {
            w = 0;
            fl_measure(data[r][c].c_str(), w, h, 0);
            if (w > wmax)
                wmax = w;
        }
        colwidths.push_back(colwidth{c, wmax});
    }
    std::sort(colwidths.begin(), colwidths.end(), [](colwidth a, colwidth b) {
        return a.width < b.width;
    });

    //for (auto i : colwidths)
    //    std::cout << "$ " << i.col << ": " << i.width << std::endl;

    int restwid = wiw;
    int cols = ncols;
    const int padwidth = 4;
    for (colwidth cw : colwidths) {
        if (cols == 1) {
            if (cw.width + padwidth < restwid) {
                col_width(cw.col, restwid);
            } else {
                col_width(cw.col, cw.width + padwidth);
            }
        } else {
            int defwid = restwid / cols;
            --cols;
            if (cw.width + padwidth < defwid) {
                col_width(cw.col, defwid);
                restwid -= defwid;
            } else {
                col_width(cw.col, cw.width + padwidth);
                restwid -= cw.width + padwidth;
            }
        }
    }
}
