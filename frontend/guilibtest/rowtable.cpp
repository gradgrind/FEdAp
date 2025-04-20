#include "rowtable.h"
#include "layout.h"
#include "widget_methods.h"
#include "widgets.h"
#include <FL/fl_draw.H>
//#include <fmt/format.h>
#include <iostream>
using namespace std;

RowTable::RowTable()
    : Fl_Table_Row(0, 0, 0, 0)
{
    col_resize(0); // disable manual column resizing
    row_resize(0); // disable manual row resizing
    row_header(0); // disable row headers (along left)
    col_header(0); // disable column headers (along top)
    type(Fl_Table_Row::SELECT_SINGLE);
}
/*// Rows
    rows(5);            // how many rows
    row_header(0);      // disable row headers (along left)
    row_height_all(30); // default height of rows
    // Cols
    cols(6);               // how many columns
    col_header(1);         // enable column headers (along top)
    col_header_height(30); // enable column headers (along top)
    //col_width_all(150);    // default width of columns
    col_header_color(fl_rgb_color(230, 230, 255));
    end(); // end the Fl_Table group
    
    // Temporary (test) data
    headers = {"Subject", "Groups", "Teachers", "Rooms", "Units", "Properties"};
    data = {{"Fr", "11.A", "PM", "fs1", "1, 1, 1", ""},
            {"Fr", "11.B", "PM", "fs1", "1, 1, 1", ""},
            {"Ma", "11", "AH", "k11", "3, 3", "HU"},
            {"Ma", "11.A", "AH", "k11", "1, 1", ""},
            {"Fr", "11.B", "AH", "k11", "1, 1", ""}};
*/

void rowtable_method(
    Fl_Widget *w, std::string_view c, minion::MinionList m)
{
    if (c == "rows") {
        static_cast<RowTable *>(w)->set_rows(int_param(m, 1));
    } else if (c == "cols") {
        static_cast<RowTable *>(w)->set_cols(int_param(m, 1));
    } else if (c == "row_header_width") {
        int rhw = int_param(m, 1);
        if (rhw) {
            static_cast<RowTable *>(w)->row_header(1);
            static_cast<RowTable *>(w)->row_header_width(rhw);
        } else {
            static_cast<RowTable *>(w)->row_header(0);
        }
    } else if (c == "col_header_height") {
        int chh = int_param(m, 1);
        if (chh) {
            static_cast<RowTable *>(w)->col_header(1);
            static_cast<RowTable *>(w)->col_header_height(chh);
        } else {
            static_cast<RowTable *>(w)->col_header(0);
        }
    } else if (c == "col_header_color") {
        static_cast<RowTable *>(w)->col_header_color(colour_param(m, 1));
    } else if (c == "row_header_color") {
        static_cast<RowTable *>(w)->row_header_color(colour_param(m, 1));
    } else if (c == "row_height_all") {
        static_cast<RowTable *>(w)->row_height_all(int_param(m, 1));
    } else if (c == "col_width_all") {
        static_cast<RowTable *>(w)->col_width_all(int_param(m, 1));
    } else if (c == "col_headers") {
        static_cast<RowTable *>(w)->col_headers.clear();
        int n = m.size() - 1;
        static_cast<RowTable *>(w)->set_cols(n);
        for (int i = 0; i < n; ++i) {
            static_cast<RowTable *>(w)->col_headers[i] = get<string>(m.at(i + 1));
        }
    } else if (c == "row_headers") {
        static_cast<RowTable *>(w)->row_headers.clear();
        int n = m.size() - 1;
        static_cast<RowTable *>(w)->set_rows(n);
        for (int i = 0; i < n; ++i) {
            static_cast<RowTable *>(w)->row_headers[i] = get<string>(m.at(i + 1));
        }
    } else {
        widget_method(w, c, m);
    }
}

// Need to handle the effect of column changes on data stores.
void RowTable::set_cols(
    int n)
{
    int nr = rows();
    cols(n);
    col_headers.resize(n);
    if (n && nr) {
        for (int i = 0; i < nr; ++i) {
            data.at(i).resize(n);
        }
    }
}

// Need to handle the effect of row changes on data stores.
void RowTable::set_rows(
    int n)
{
    int nc = cols();
    rows(n);
    row_headers.resize(n);
    if (n && nc) {
        data.resize(n, vector<string>(nc));
    }
}

Fl_Widget *NEW_RowTable(
    minion::MinionMap param)
{
    return new RowTable();
}

void RowTable::draw_cell(
    TableContext context, int ROW, int COL, int X, int Y, int W, int H)
{
    switch (context) {
    case CONTEXT_STARTPAGE: // before page is drawn..
        //fl_font(FL_HELVETICA, 16); // set the font for our drawing operations

        // Adjust column widths
        size_columns();
        //TODO: Adjust width of row headers?

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
        fl_draw(col_headers[COL].c_str(), X, Y, W, H, FL_ALIGN_CENTER);
        fl_pop_clip();
        return;
    case CONTEXT_ROW_HEADER: // Draw row headers
        fl_push_clip(X, Y, W, H);
        fl_draw_box(FL_THIN_UP_BOX, X, Y, W, H, row_header_color());
        fl_color(FL_BLACK);
        fl_draw(row_headers[ROW].c_str(), X, Y, W, H, FL_ALIGN_CENTER);
        fl_pop_clip();
        return;
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
void RowTable::size_columns()
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
        fl_measure(col_headers[c].c_str(), w, h, 0);
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

void RowTable::_row_cb(
    void *table)
{
    //TODO
    cout << "§§§ " << static_cast<RowTable *>(table)->_current_row << endl;
}
