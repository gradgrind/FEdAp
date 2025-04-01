#include "fltk_json.h"

// *** Dispatch table for widget-creation functions
// First create a function template
using new_widget_func
    = function<void(string_view name, string_view parent, json data)>;
// Create the map object
unordered_map<string_view, new_widget_func> new_widget_map{{"Window",
                                                            new_window},
                                                           {"Flex", new_flex},
                                                           {"Grid", new_grid}};

// Call a method of a given widget, passing any parameters as a JSON object.
void gui(
    string_view widget, string_view method, json data)
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
    string_view name, string_view widget_type, string_view parent, json data)
{
    try {
        auto fn = new_widget_map.at(widget_type);
        fn(name, parent, data);
    } catch (const std::out_of_range& e) {
        throw fmt::format("Unknown widget type: {} ({})", widget_type, e.what());
    }
}

// Pass a message to the back-end. This can be an event/callback, the
// reponse to a query, or whatever.
void message(
    json data)
{}
