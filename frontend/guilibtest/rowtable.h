#ifndef ROWTABLE_H
#define ROWTABLE_H

#include <FL/Fl_Table_Row.H>
#include <string>
#include <vector>

class RowTable : public Fl_Table_Row
{
    // Handle drawing table's cells
    //     Fl_Table calls this function to draw each visible cell in the
    //     table. It's up to us to use FLTK's drawing functions to draw
    //     the cells the way we want.
    //
    void draw_cell(TableContext context,
                   int ROW = 0,
                   int COL = 0,
                   int X = 0,
                   int Y = 0,
                   int W = 0,
                   int H = 0) FL_OVERRIDE;

    static void _row_cb(void* table);

    int _current_row = -1;

    void size_columns();

public:
    RowTable();

    void set_cols(int n);
    void set_rows(int n);

    Fl_Color bg{0xf0f0f000};
    Fl_Color header_bg{0xe0e0e000};

    std::vector<std::vector<std::string>> data; // data array for cells
    std::vector<std::string> row_headers;
    std::vector<std::string> col_headers;
};

#endif // ROWTABLE_H
