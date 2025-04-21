#include "layout.h"
#include "minion.h"
#include "widget_base.h"
#include "widget_methods.h"
#include <FL/Fl_Double_Window.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Grid.H>
#include <fmt/format.h>
#include <iostream>
using namespace std;
using mlist = minion::MinionList;
using mmap = minion::MinionMap;

// *** "Group" widgets ***

void w_widget_method(
    Widget wd, string_view c, mlist m)
{
    int ww, wh;
    if (c == "SIZE") {
        ww = int_param(m, 1); // width
        wh = int_param(m, 2); // height
        cout << "?SIZE? " << ww << " ? " << wh << endl;
        wd.widget->size(ww, wh);
    } else if (c == "COLOUR") {
        auto clr = get_colour(get<string>(m.at(1)));
        wd.widget->color(clr);
    } else if (c == "BOXTYPE") {
        auto bxt = get_boxtype(get<string>(m.at(1)));
        wd.widget->box(bxt);
    } else if (c == "LABEL") {
        auto lbl = get<string>(m.at(1));
        wd.widget->copy_label(lbl.c_str());
    } else if (c == "CALLBACK") {
        auto cb = get<string>(m.at(1));
        wd.widget->callback(do_callback);
    } else if (c == "SHOW") {
        wd.widget->show();
    } else if (c == "FIXED") {
        auto parent = dynamic_cast<Fl_Flex *>(wd.widget->parent());
        if (parent) {
            int sz = int_param(m, 1);
            parent->fixed(wd.widget, sz);
        } else {
            //TODO: how to get widget name?
            //throw fmt::format("Widget ({}) method FIXED: parent not VLayout/Hlayout",
            //                  WidgetData::get_widget_name(w));
            throw "Widget method FIXED: parent not VLayout/Hlayout";
        }
    } else if (c == "clear_visible_focus") {
        wd.widget->clear_visible_focus();
    } else if (c == "measure_label") {
        int wl, hl;
        wd.widget->measure_label(wl, hl);
        //TODO ... how to get widget name?
        //cout << "Measure " << WidgetData::get_widget_name(w) << " label: " << wl << ", " << hl
        //     << endl;
        cout << "Measure label: " << wl << ", " << hl << endl;
    } else {
        throw fmt::format("Unknown widget method: {}", c);
    }
}

void w_group_method(
    Widget wd, string_view c, mlist m)
{
    if (c == "RESIZABLE") {
        auto rsw = widget_info.at(get<string>(m.at(1)));
        static_cast<Fl_Group *>(wd.widget)->resizable(rsw.widget);
    } else {
        w_widget_method(wd, c, m);
    }
}

void w_flex_method(
    Widget wd, string_view c, mlist m)
{
    if (c == "MARGIN") {
        int sz = int_param(m, 1);
        static_cast<Fl_Flex *>(wd.widget)->margin(sz);
    } else {
        w_group_method(wd, c, m);
    }
}

void w_grid_method(
    Widget wd, string_view c, mlist m)
{
    if (c == "???") {
        // TODO
    } else {
        w_group_method(wd, c, m);
    }
}

void w_callback_no_esc_closes(
    Fl_Widget *w, void *x)
{
    if (Fl::event() == FL_SHORTCUT && Fl::event_key() == FL_Escape)
        return; // ignore Escape
    //TODO: message to backend? How to get widget name?
    cout << "Closing Window" << endl;
    //TODO--
    exit(0);
}

Fl_Widget *w_NEW_Window(
    mmap param)
{
    int w = 800;
    int h = 600;
    param.get_int("WIDTH", w);
    param.get_int("HEIGHT", h);
    auto widg = new Fl_Double_Window(w, h);
    int esc_closes{0};
    param.get_int("ESC_CLOSES", esc_closes);
    if (!esc_closes)
        widg->callback(w_callback_no_esc_closes);

    Fl_Group::current(0); // disable "auto-grouping"
    return widg;
}
