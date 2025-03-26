#ifndef COURSES_GUI_H
#define COURSES_GUI_H

#include <FL/Fl_Box.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Table_Row.H>
#include <map>
#include <string>
#include <vector>

class CourseTable : public Fl_Table_Row
{
    std::vector<std::vector<std::string>> data; // data array for cells
    std::vector<std::string> headers;

    // Draw the (row/)col headings
    void DrawHeader(const char *s, int X, int Y, int W, int H);

    // Draw the cell data
    void DrawData(const char *s, int X, int Y, int W, int H);

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

public:
    CourseTable();
};

class CoursesGui : public Fl_Flex
{
public:
    CoursesGui();

    std::map<std::string, Fl_Widget *> Widgets;
};

#endif // COURSES_GUI_H
