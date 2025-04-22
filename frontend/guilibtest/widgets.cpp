#include "layout.h"
#include "minion.h"
#include "widget_methods.h"
#include <FL/Fl_Box.H>
#include <FL/Fl_Choice.H>
#include <FL/Fl_Flex.H>
#include <FL/Fl_Output.H>
using namespace std;
using mmap = minion::MinionMap;
using mlist = minion::MinionList;

void choice_method(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "ADD") {
        for (int i = 1; i < m.size(); ++i) {
            //TODO: Do I need to store the string somewhere, or is that
            // handled by the widget?
            static_cast<Fl_Choice *>(w)->add(get<string>(m.at(i)).c_str());
        }
    } else if (c == "LABEL") {
        left_label(w, m);
    } else {
        widget_method(w, c, m);
    }
}

void input_method(
    Fl_Widget *w, string_view c, mlist m)
{
    if (c == "VALUE") {
        //TODO: Do I need to store the string somewhere, or is that
        // handled by the widget?
        static_cast<Fl_Input *>(w)->value(get<string>(m.at(1)).c_str());
    } else if (c == "LABEL") {
        left_label(w, m);
    } else {
        widget_method(w, c, m);
    }
}

// *** Non-layout widgets ***

Fl_Widget *NEW_Box(
    mmap param)
{
    return new Fl_Box(0, 0, 0, 0);
}

Fl_Widget *NEW_Choice(
    mmap param)
{
    return new Fl_Choice(0, 0, 0, 0);
}

Fl_Widget *NEW_Output(
    mmap param)
{
    return new Fl_Output(0, 0, 0, 0);
}
