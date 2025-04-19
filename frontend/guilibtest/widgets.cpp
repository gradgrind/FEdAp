#include "fltk_minion.h"
#include <FL/Fl_Box.H>
using namespace std;

// *** Non-layout widgets ***

Fl_Widget *NEW_Box(
    string_view name, mlist do_list)
{
    auto widg = new Fl_Box(0, 0, 0, 0);
    for (const auto &cmd : do_list) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        //_widget_methods(widg, c, m);
    }
    return widg;
}
