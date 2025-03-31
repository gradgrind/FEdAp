#include "fltk_json.h"
#include <iostream>

typedef map<string_view, void *> member_map; // TODO

unordered_map<string_view, int> widget_type_map;
vector<member_map> widget_type_vec;
unordered_map<string_view, Fl_Widget *> widget_name_map;

WidgetData::WidgetData(
    string_view w_type, string_view w_name, Fl_Widget *widget)
    : Fl_Callback_User_Data()
    , wname{w_name}
{
    int i;
    if (widget_type_map.contains(w_type)) {
        i = widget_type_map.at(w_type);
    } else {
        i = widget_type_vec.size();
        widget_type_vec.push_back(member_map{});
        widget_type_map.emplace(w_type, i);
    }
    wtype = i;
    add_widget(widget);
    widget->user_data(this, true);
}

void WidgetData::add_widget(
    Fl_Widget *w)
{
    // Allow unnamed widgets. These are not placed in the map.
    if (wname.empty())
        return;
    if (widget_name_map.contains(wname)) {
        throw fmt::format("Widget name already exists: %s", wname);
    }
    widget_name_map.emplace(wname, w);
}

void WidgetData::remove_widget(
    std::string_view name)
{
    // Allow unnamed widgets. These are not placed in the map.
    if (wname.empty())
        return;
    if (widget_name_map.erase(wname) == 0) {
        throw fmt::format("Can't remove widget '%s', it doesn't exist", wname);
    }
}

// The user data might need deleting
WidgetData::~WidgetData()
{
    if (auto_delete_user_data && user_data)
        delete (Fl_Callback_User_Data *) user_data;
}
