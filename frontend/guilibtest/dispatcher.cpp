#include "fltk_json.h"

// *** Dispatch table for widget-creation functions
// First create a function template
using new_widget_func
    = function<void(string_view name, string_view parent, json data)>;
// Create the map object
unordered_map<string_view, new_widget_func> function_map{{"Window", new_window},
                                                         {"Flex", new_flex},
                                                         {"Grid", new_grid},
                                                         {"Box", new_box}};

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
        auto fn = function_map.at(widget_type);
        fn(name, parent, data);
    } catch (const std::out_of_range& e) {
        throw fmt::format("Unknown widget type: {} ({})", widget_type, e.what());
    }
}

// NEW ...
using _new_widget_func = function<void(json data)>;
unordered_map<string_view, _new_widget_func> _function_map{};

void widget_method(
    Fl_Widget* w, json obj)
{
    string m;
    if (get_json_string(obj, "M", m)) {
        try {
            auto wd{static_cast<WidgetData*>(w->user_data())};
            //            auto fn = _function_map.at(fw);
            //fn(obj);
        } catch (const std::out_of_range& e) {
            //            throw fmt::format("Unknown function: {} ({})", fw, e.what());
        }
    }
}

void GUI(
    json obj)
{
    string fw;
    if (get_json_string(obj, "F", fw)) {
        try {
            auto fn = _function_map.at(fw);
            fn(obj);
        } catch (const std::out_of_range& e) {
            throw fmt::format("Unknown function: {} ({})", fw, e.what());
        }
    } else if (get_json_string(obj, "W", fw)) {
        auto w = get_widget(fw);
        json::array_t mlist = obj.at("DO");
        for (const auto& m : mlist) {
            //TODO: widget_method(w, m);
        }
    } else {
        throw fmt::format("Invalid GUI data: {}", obj);
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

//TODO
json message(
    json data)
{
    return data;
}

json to_back_end(
    json data)
{
    json result = message(data);
    auto it = result.find("DO");
    if (it != data.end()) {
        for (const auto& cmd : it.value()) {
            GUI(cmd);
        }
    }

    return result;
}
