#include "functions.h"
#include "layout.h"
#include "minion.h"
#include "widgets.h"
#include <FL/Fl_Group.H>
#include <fmt/format.h>
#include <iostream>
using mmap = minion::MinionMap;
using mlist = minion::MinionList;
using namespace std;

void Handle_methods(
    Fl_Widget* w, mmap m, method_handler h)
{
    auto dolist = m.get("DO");
    if (holds_alternative<mlist>(dolist)) {
        mlist do_list = get<mlist>(dolist);
        for (const auto& cmd : do_list) {
            mlist m = get<mlist>(cmd);
            string_view c = get<string>(m.at(0));
            h(w, c, m);
        }
    } else if (dolist.index() != 0) {
        string s;
        minion::dump(s, dolist, 0);
        throw string{"Invalid DO list: "} + s;
    }
}

void Handle_NEW(
    string_view wtype, mmap m)
{
    //cout << "Handle_NEW " << wtype << ":" << minion::dump_map_items(m, 1) << endl;
    string name;
    Fl_Widget* w;
    method_handler h;
    if (m.get_string("NAME", name)) {
        if (wtype == "Window") {
            w = NEW_Window(m);
            h = group_method;
        } else if (wtype == "Vlayout") {
            w = NEW_Vlayout(m);
            h = flex_method;
        } else if (wtype == "Hlayout") {
            w = NEW_Hlayout(m);
            h = flex_method;
        } else if (wtype == "Grid") {
            w = NEW_Grid(m);
            h = grid_method;
            // *** End of layouts, start of other widgets
        } else if (wtype == "Box") {
            w = NEW_Box(m);
            h = widget_method;
        } else if (wtype == "Choice") {
            w = NEW_Choice(m);
            h = choice_method;
        } else if (wtype == "Output") {
            w = NEW_Output(m);
            h = input_method;
        } else if (wtype == "RowTable") {
            w = NEW_RowTable(m);
            h = rowtable_method;
        } else if (wtype == "EditForm") {
            w = NEW_EditForm(m);
            h = editform_method;
        } else {
            throw fmt::format("Unknown widget type: {}", wtype);
        }
        string parent;
        if (m.get_string("PARENT", parent) && !parent.empty()) {
            static_cast<Fl_Group*>(WidgetData::get_widget(parent))->add(w);
        }
        // Add a WidgetData as "user data" to the widget
        WidgetData::add_widget(name, w, h);
        // Handle methods
        Handle_methods(w, m, h);
        return;
    }
    throw fmt::format("Bad NEW command: {}", minion::dump_map_items(m, -1));
}

void GUI(
    mmap obj)
{
    string w;
    if (obj.get_string("NEW", w)) {
        Handle_NEW(w, obj);
    } else if (obj.get_string("WIDGET", w)) {
        // Handle methods
        auto widg = WidgetData::get_widget(w);
        auto wd{static_cast<WidgetData*>(widg->user_data())};
        Handle_methods(widg, obj, wd->handle_method);
    } else if (obj.get_string("FUNCTION", w)) {
        auto f = function_map.at(w);
        f(obj);
    } else {
        throw fmt::format("Invalid GUI parameters: {}", dump_map_items(obj, -1));
    }
}

//TODO ...
void tmp_run(
    mmap data)
{
    auto dolist0 = data.get("GUI");
    if (holds_alternative<mlist>(dolist0)) {
        auto dolist = get<mlist>(dolist0);
        for (const auto& cmd : dolist) {
            GUI(get<mmap>(cmd));
        }
    } else {
        cerr << "Input data not a GUI command list" << endl;
    }
}

// Pass a message to the back-end. This can be an event/callback, the
// reponse to a query, or whatever.

// There need to be two kinds of "message":
// 1) Let's call this a virtual override. It is a call to the back-end,
//    perhaps with a result (like 0 or 1 for event handlers), and is
//    blocking – so it should execute quickly. Unfortunately this seems
//    very difficult to implement, because it might also need to query
//    the front-end or perform other gui operations. Thus it entails
//    a calling back and forth between back-end and front-end.
// 2) Let's call this a trigger. It sets an operation in the back-end
//    going, but doesn't wait for it to finish. Any resulting calls to
//    the front-end could be picked up by an idle function.
// For the moment I would like to implement just normal callbacks, i.e.
// asynchronous calls. Where event handlers are necessary, I would first
// consider extending the C++ widgets.

//TODO: ???
void do_callback(
    Fl_Widget* w, void* x)
{
    auto wd{static_cast<WidgetData*>(w->user_data())};
    cout << "Callback: " << wd->widget_name() << endl;
}

//TODO
mmap message(
    mmap data)
{
    return data;
}

void to_back_end(
    mmap data)
{
    mmap result = message(data);
    auto dolist0 = data.get("DO");
    if (holds_alternative<mlist>(dolist0)) {
        auto dolist = get<mlist>(dolist0);
        for (const auto& cmd : dolist) {
            GUI(get<mmap>(cmd));
        }
    }
    // Any back-end function which can take more than about 100ms should
    // initiate a timeout leading to a modal "progress" dialog.
    // Any data generated while such a callback is operating (i.e. before
    // it returns a completion code) should be fetched and run by an idle
    // handler. Any data generated outside of this period is probably an
    // error – the back-end should not be doing anything then!
}
