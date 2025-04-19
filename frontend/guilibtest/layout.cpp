#include "layout.h"
//#include "fltk_minion.h"
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
using namespace std;
#include "minion.h"
#include <fmt/format.h>
using mlist = minion::MinionList;
using mmap = minion::MinionMap;

// *** "Group" widgets ***

void widget_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    int ww, wh;
    if (c == "SIZE") {
        ww = stoi(get<string>(m.at(1))); // width
        wh = stoi(get<string>(m.at(2))); // height
        w->size(ww, wh);
    } else {
        throw fmt::format("Unknown widget method: {}", c);
    }
}

void group_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "???") {
        // TODO
    } else {
        widget_methods(w, c, m);
    }
}

void flex_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "???") {
        // TODO
    } else {
        group_methods(w, c, m);
    }
}

void grid_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "???") {
        // TODO
    } else {
        group_methods(w, c, m);
    }
}

Fl_Widget *NEW_Window(
    mmap param)
{
    int w = 800;
    int h = 600;
    param.get_int("WIDTH", w);
    param.get_int("HEIGHT", h);
    auto widg = new Fl_Double_Window(w, h);
    Fl_Group::current(0); // disable "auto-grouping"
    return widg;
}

Fl_Widget *NEW_Vlayout(
    mmap param)
{
    auto widg = new Fl_Flex(Fl_Flex::COLUMN);
    Fl_Group::current(0); // disable "auto-grouping"
    return widg;
}

Fl_Widget *NEW_Hlayout(
    mmap param)
{
    auto widg = new Fl_Flex(Fl_Flex::ROW);
    Fl_Group::current(0); // disable "auto-grouping"
    return widg;
}

Fl_Widget *NEW_Grid(
    mmap param)
{
    auto widg = new Fl_Grid(0, 0, 0, 0);
    Fl_Group::current(0); // disable "auto-grouping"
    return widg;
}
