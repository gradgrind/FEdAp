#include "minion.h"
#include "widget_base.h"
#include <FL/Fl_Box.H>
#include <FL/Fl_Choice.H>
#include <FL/Fl_Output.H>
using namespace std;
using mmap = minion::MinionMap;
using mlist = minion::MinionList;

void w_choice_method(
    Widget wd, string_view c, mlist m)
{
    if (c == "ADD") {
        for (int i = 1; i < m.size(); ++i) {
            //TODO: Do I need to store the string somewhere, or is that
            // handled by the widget?
            static_cast<Fl_Choice *>(wd.widget)->add(get<string>(m.at(i)).c_str());
        }
    } else {
        w_widget_method(wd, c, m);
    }
}

void w_input_method(
    Widget wd, string_view c, mlist m)
{
    if (c == "VALUE") {
        //TODO: Do I need to store the string somewhere, or is that
        // handled by the widget?
        static_cast<Fl_Input *>(wd.widget)->value(get<string>(m.at(1)).c_str());
    } else if (c == "???") {
        //TODO
    } else {
        w_widget_method(wd, c, m);
    }
}
