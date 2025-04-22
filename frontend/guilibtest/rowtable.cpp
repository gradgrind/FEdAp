#include "rowtable.h"
#include "layout.h"
#include "widget_methods.h"
#include "widgets.h"
#include <FL/fl_draw.H>
//#include <fmt/format.h>
#include <iostream>
using namespace std;
using namespace minion;

//TODO: When the table doesn't fit in the space, scrollbars should appear.
// It is not clear to me why they don't!

//TODO: After changes I probably need to force a redraw (redraw()).

RowTable::RowTable()
    : Fl_Table_Row(0, 0, 0, 0)
{
    col_resize(0); // disable manual column resizing
    row_resize(0); // disable manual row resizing
    row_header(0); // disable row headers (along left)
    col_header(0); // disable column headers (along top)
    type(Fl_Table_Row::SELECT_SINGLE);
}

void rowtable_method(
    Fl_Widget *w, string_view c, MinionList m)
{
    auto t = static_cast<RowTable *>(w);

    //TODO: It would probably be better to use an enum and switch!
    if (c == "rows") {
        t->set_rows(int_param(m, 1));
    } else if (c == "cols") {
        t->set_cols(int_param(m, 1));
    } else if (c == "row_header_width") {
        int rhw = int_param(m, 1);
        if (rhw) {
            t->row_header(1);
            t->row_header_width(rhw);
        } else {
            t->row_header(0);
        }
    } else if (c == "col_header_height") {
        int chh = int_param(m, 1);
        if (chh) {
            t->col_header(1);
            t->col_header_height(chh);
        } else {
            t->col_header(0);
        }
    } else if (c == "col_header_color") {
        t->col_header_color(colour_param(m, 1));
    } else if (c == "row_header_color") {
        t->row_header_color(colour_param(m, 1));
    } else if (c == "row_height_all") {
        t->row_height_all(int_param(m, 1));
    } else if (c == "col_width_all") {
        t->col_width_all(int_param(m, 1));
    } else if (c == "col_headers") {
        t->col_headers.clear();
        int n = m.size() - 1;
        t->set_cols(n);
        for (int i = 0; i < n; ++i) {
            t->col_headers[i] = get<string>(m.at(i + 1));
        }
    } else if (c == "row_headers") {
        t->row_headers.clear();
        int n = m.size() - 1;
        t->set_rows(n);
        for (int i = 0; i < n; ++i) {
            t->row_headers[i] = get<string>(m.at(i + 1));
        }
    } else if (c == "add_row") {
        int n = m.size() - 2;
        if (n != t->cols()) {
            throw "RowTable: add_row with wrong length";
        }
        t->row_headers.emplace_back(get<string>(m.at(1)));
        vector<string> r(n);
        for (int i = 0; i < n; ++i) {
            r[i] = get<string>(m.at(i + 2));
        }
        t->data.emplace_back(r);
        t->rows(t->rows() + 1);
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
    MinionMap param)
{
    auto widg = new RowTable();
    widg->color(widg->bg);
    widg->col_header_color(widg->header_bg);
    widg->row_header_color(widg->header_bg);
    return widg;
}

void RowTable::draw_cell(
    TableContext context, int ROW, int COL, int X, int Y, int W, int H)
{
    switch (context) {
    case CONTEXT_STARTPAGE: // before page is drawn..
        //fl_font(FL_HELVETICA, 16); // set the font for our drawing operations

        //TODO: Adjust width of row headers?
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
        fl_color(0xb0b0b000);
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
    int ncols = cols();
    int nrows = rows();
    int cw, ch, wmax;

    // Deal with row headers
    wmax = 0;
    if (row_header()) {
        for (int r = 0; r < nrows; ++r) {
            cw = 0;
            fl_measure(row_headers[r].c_str(), cw, ch, 0);
            if (cw > wmax)
                wmax = cw;
        }
    }
    int rhw = wmax + 10;
    row_header_width(rhw);

    // Get widest column entries
    struct colwidth
    {
        int col, width;
    };
    std::vector<colwidth> colwidths;

    for (int c = 0; c < ncols; ++c) {
        cw = 0;
        fl_measure(col_headers[c].c_str(), cw, ch, 0);
        wmax = cw;
        for (int r = 0; r < nrows; ++r) {
            cw = 0;
            fl_measure(data[r][c].c_str(), cw, ch, 0);
            if (cw > wmax)
                wmax = cw;
        }
        colwidths.emplace_back(colwidth{c, wmax});
    }
    std::sort(colwidths.begin(), colwidths.end(), [](colwidth a, colwidth b) {
        return a.width > b.width;
    });

    //for (const auto &i : colwidths)
    //    cout << "$ " << i.col << ": " << i.width << endl;

    int restwid = tiw;
    int icols = ncols;
    const int padwidth = 4;
    for (colwidth cw : colwidths) {
        if (icols == 1) {
            if (cw.width + padwidth < restwid) {
                col_width(cw.col, restwid);
            } else {
                col_width(cw.col, cw.width + padwidth);
            }
        } else {
            int defwid = restwid / icols;
            --icols;
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
