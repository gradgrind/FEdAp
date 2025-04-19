#include "layout.h"
#include "widget_methods.h"
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
using namespace std;
#include "minion.h"
#include <fmt/format.h>
using mlist = minion::MinionList;
using mmap = minion::MinionMap;

// *** "Group" widgets ***

void widget_method(
    Fl_Widget *w, string_view c, mlist m)
{
    int ww, wh;
    if (c == "SIZE") {
        ww = stoi(get<string>(m.at(1))); // width
        wh = stoi(get<string>(m.at(2))); // height
        w->size(ww, wh);
    } else if (c == "COLOUR") {
        auto clr = get_colour(get<string>(m.at(1)));
        w->color(clr);
    } else if (c == "BOXTYPE") {
        auto bxt = get_boxtype(get<string>(m.at(1)));
        w->box(bxt);
    } else if (c == "LABEL") {
        auto lbl = get<string>(m.at(1));
        w->copy_label(lbl.c_str());
    } else {
        throw fmt::format("Unknown widget method: {}", c);
    }
}

void group_method(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "???") {
        // TODO
    } else {
        widget_method(w, c, m);
    }
}

void flex_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "FIXED") {
        auto parent = static_cast<Fl_Flex *>(w->parent());
        int sz = stoi(get<string>(m.at(1)));
        parent->fixed(w, sz);
    } else {
        group_method(w, c, m);
    }
}

void grid_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "???") {
        // TODO
    } else {
        group_method(w, c, m);
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
