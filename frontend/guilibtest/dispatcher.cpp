#include "fltk_minion.h"
#include "layout.h"
#include <FL/Fl_Group.H>
using namespace std;

void Handle_NEW(
    string_view wtype, mmap m)
{
    string name;
    Fl_Widget* w;
    method_handler h;
    if (get_minion_string(m, "NAME", name)) {
        mlist do_list = get<mlist>(m.get("DO"));
        if (wtype == "Window") {
            w = NEW_Window(name, do_list);
            h = group_methods;
            //new WidgetData("Group:Window:Double", name, w);
        } else if (wtype == "Vlayout") {
            w = NEW_Vlayout(name, do_list);
            h = flex_methods;
            //new WidgetData("Group:Flex", name, w);
        } else if (wtype == "Hlayout") {
            w = NEW_Hlayout(name, do_list);
            h = flex_methods;
            //new WidgetData("Group:Flex", name, w);
        } else if (wtype == "Grid") {
            w = NEW_Grid(name, do_list);
            h = grid_methods;
            //new WidgetData("Group:Grid", name, w);
        } else {
            throw fmt::format("Unknown widget type: {}", wtype);
        }
        string parent;
        if (get_minion_string(m, "PARENT", parent)) {
            static_cast<Fl_Group*>(WidgetData::get_widget(parent))->add(w);
        }
        //TODO
        WidgetData::add_widget(name, w, h);
        //TODO: handle methods

        return;
    }
    throw fmt::format("Bad NEW command: {}", minion::dump_map_items(m, -1));
}

void Handle_WIDGET(
    Fl_Widget* w, mlist param)
{
    //TODO: This needs to find the correct starting point in the method-
    // handler hierarchy ...
}

// *** Dispatch table for widget-creation functions
// First create a function template
using new_widget_func = function<void(string_view name, string_view parent, mmap data)>;
// Create the map object
using function_map = unordered_map<string_view, new_widget_func>;
function_map fmap{{"Window", new_window}, {"Flex", new_flex}, {"Grid", new_grid}, {"Box", new_box}};

// Call a method of a given widget, passing any parameters as a mmap object.
void gui(
    string_view widget, string_view method, mmap data)
{}

// The widget-creation functions could be treated as methods
// on a group parent. That means there would need to be a sort of
// top-level set of functions (i.e. not really methods) which are
// accessible to the group widgets. The newly-created widgets would
// be added to the parent group. This could look like:
//    gui("MyParent", "Box", R"({"name": "MyWidget")")
// As a special case, parent = "" could create unattached widgets.
// A main window would be one of these.

// On the other hand, it might be nice to have a widget definition
// of the form:
//    gui("MyWidget", "Box", ...)
// That might require placing its parent in the data:
//    gui("MyWidget", "Box", R"({"parent": "MyParent"})")
// In this case, the widget-creation functions should not be available
// as methods on any particular widget (type), but only accessible if
// the widget (name) is undefined, including "" (in which case there
// must be a parent?).

// Another possibility might be a distinct creation function, say:
//    gui_new("MyWidget", "Box", "MyParent", ...)
// That might be the most "logical", so let's do it for now.

void gui_new(
    string_view name, string_view widget_type, string_view parent, mmap data)
{
    try {
        auto fn = fmap.at(widget_type);
        fn(name, parent, data);
    } catch (const std::out_of_range& e) {
        throw fmt::format("Unknown widget type: {} ({})", widget_type, e.what());
    }
}

// NEW ...
using _new_widget_func = function<void(mmap data)>;
using _function_map = unordered_map<string_view, _new_widget_func>;
_function_map fn_map{};

void widget_method(
    Fl_Widget* w, mmap obj)
{
    auto wd{static_cast<WidgetData*>(w->user_data())};
    string m;
    if (get_minion_string(obj, "M", m)) {
        wd->do_method(w, m, obj);
    } else {
        throw fmt::format("Invalid method on {}: {}", wd->widget_name(), dump_map_items(obj, -1));
    }
}

void GUI(
    mmap obj)
{
    string w;
    if (get_minion_string(obj, "NEW", w)) {
        Handle_NEW(w, obj);
    } else if (get_minion_string(obj, "WIDGET", w)) {
        mlist do_list = get<mlist>(obj.get("DO"));
        auto widg = get_widget(w);
        Handle_WIDGET(widg, do_list);
    } else if (get_minion_string(obj, "FUNCTION", w)) {
        //TODO
        //auto f = function_map.at(w);
        //f(obj);
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
