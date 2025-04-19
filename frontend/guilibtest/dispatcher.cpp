#include "fltk_minion.h"
#include "layout.h"
#include <FL/Fl_Group.H>
using namespace std;

void Handle_methods(
    Fl_Widget* w, mmap m, method_handler h)
{
    mlist do_list = get<mlist>(m.get("DO"));
    for (const auto& cmd : do_list) {
        mlist m = get<mlist>(cmd);
        string_view c = get<string>(m.at(0));
        h(w, c, m);
    }
}

void Handle_NEW(
    string_view wtype, mmap m)
{
    string name;
    Fl_Widget* w;
    method_handler h;
    if (get_minion_string(m, "NAME", name)) {
        if (wtype == "Window") {
            w = NEW_Window(name, m);
            h = group_methods;
        } else if (wtype == "Vlayout") {
            w = NEW_Vlayout(name, m);
            h = flex_methods;
        } else if (wtype == "Hlayout") {
            w = NEW_Hlayout(name, m);
            h = flex_methods;
        } else if (wtype == "Grid") {
            w = NEW_Grid(name, m);
            h = grid_methods;
        } else {
            throw fmt::format("Unknown widget type: {}", wtype);
        }
        string parent;
        if (get_minion_string(m, "PARENT", parent)) {
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

using function_handler = std::function<void(minion::MinionMap)>;
unordered_map<string, function_handler> function_map;

void GUI(
    mmap obj)
{
    string w;
    if (get_minion_string(obj, "NEW", w)) {
        Handle_NEW(w, obj);
    } else if (get_minion_string(obj, "WIDGET", w)) {
        // Handle methods
        auto widg = WidgetData::get_widget(w);
        auto wd{static_cast<WidgetData*>(widg->user_data())};
        Handle_methods(widg, obj, wd->handle_method);
    } else if (get_minion_string(obj, "FUNCTION", w)) {
        auto f = function_map.at(w);
        f(obj);
    } else {
        throw fmt::format("Invalid GUI parameters: {}", dump_map_items(obj, -1));
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
