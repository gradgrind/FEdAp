#include "layout.h"
#include <fmt/format.h>
#include <iostream>
using namespace std;
using mlist = minion::MinionList;

//static member variable
std::unordered_map<std::string_view, Fl_Widget *> WidgetData::widget_map;

// static
mlist WidgetData::list_widgets()
{
    mlist keys;
    for (const auto &kv : widget_map) {
        keys.emplace_back(string{kv.first});
    }
    return keys;
}

// static
Fl_Widget *WidgetData::get_widget(
    string_view name)
{
    try {
        return widget_map.at(name);
    } catch (const std::out_of_range &e) {
        throw fmt::format("Unknown widget: {} ({})", name, e.what());
    }
}

// static
void WidgetData::add_widget(
    string_view name, Fl_Widget *w, method_handler h)
{
    if (widget_map.contains(name)) {
        throw fmt::format("Widget name already exists: {}", name);
    }
    auto wd = new WidgetData(name, h);
    widget_map.emplace(wd->w_name, w);
    w->user_data(wd, true); // auto-free = true
}

// static
string_view WidgetData::get_widget_name(
    Fl_Widget *w)
{
    auto wd{static_cast<WidgetData *>(w->user_data())};
    return wd->w_name;
}

WidgetData::WidgetData(
    string_view wname, method_handler h)
    : Fl_Callback_User_Data()
    , w_name{wname}
    , handle_method{h}
{
    //TODO: widget type?
}

// The destructor is called when the widget to which this WidgetData is
// attached is deleted – if the widget's user-data was set with auto-free
// true. The widget reference must be removed from widget_map and the
// user-data associated with the WidgetData might need deleting.
WidgetData::~WidgetData()
{
    if (widget_map.erase(w_name) == 0) {
        cerr << fmt::format("Can't remove widget '{}', it doesn't exist", w_name) << endl;
    }
    if (auto_delete_user_data && user_data)
        delete static_cast<Fl_Callback_User_Data *>(user_data);
}

string_view WidgetData::widget_name()
{
    return w_name;
}

/*
int WidgetData::widget_type()
{
    return wtype;
}

string_view WidgetData::widget_type_name()
{
    auto i = widget_type();
    return widget_type_names[i];
}
*/
