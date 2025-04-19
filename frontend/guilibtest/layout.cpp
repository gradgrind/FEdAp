#include "layout.h"
//#include "fltk_minion.h"
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
using namespace std;
#include "minion.h"
#include <fmt/format.h>
using mlist = minion::MinionList;

// *** "Group" widgets ***

void widget_methods(
    Fl_Widget *w, string_view c, mlist m)
{
    int ww, wh;
    if (c == "SIZE") {
        ww = stoi(get<string>(m.at(1)));
        wh = stoi(get<string>(m.at(2)));
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
    string_view name, mlist do_list)
{
    int w = 800;
    int h = 600;
    mlist do_stripped;
    for (const auto &cmd : do_list) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        if (c == "WIDTH") {
            w = stoi(get<string>(m.at(1)));
        } else if (c == "HEIGHT") {
            h = stoi(get<string>(m.at(1)));
        } else {
            do_stripped.emplace_back(m);
        }
    }

    auto widg = new Fl_Double_Window(w, h);
    Fl_Group::current(0); // disable "auto-grouping"

    for (const auto &cmd : do_stripped) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        group_methods(widg, c, m);
    }
    return widg;
}

Fl_Widget *NEW_Vlayout(
    string_view name, mlist do_list)
{
    auto widg = new Fl_Flex(Fl_Flex::COLUMN);
    Fl_Group::current(0); // disable "auto-grouping"
    for (const auto &cmd : do_list) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        flex_methods(widg, c, m);
    }
    return widg;
}

Fl_Widget *NEW_Hlayout(
    string_view name, mlist do_list)
{
    auto widg = new Fl_Flex(Fl_Flex::ROW);
    Fl_Group::current(0); // disable "auto-grouping"
    for (const auto &cmd : do_list) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        flex_methods(widg, c, m);
    }
    return widg;
}

Fl_Widget *NEW_Grid(
    string_view name, mlist do_list)
{
    auto widg = new Fl_Grid(0, 0, 0, 0);
    Fl_Group::current(0); // disable "auto-grouping"
    for (const auto &cmd : do_list) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        grid_methods(widg, c, m);
    }
    return widg;
}
