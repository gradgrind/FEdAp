#ifndef FLTK_JSON_H
#define FLTK_JSON_H

#include <FL/Fl_Widget.H>

#include <fmt/format.h>
#include <string>
#include <string_view>

#include <json.hpp>
using json = nlohmann::json;

using namespace std;

// The user_data in Fl_Widget is used to store additional fields,
// primarily the widget's name and type. As this feature is now no
// longer directly available to the programmer, another field (user_data)
// is provided within the WidgetData subclass of Fl_Callback_User_Data
// which is referred to by the original user data field.

//extern unordered_map<string_view, Fl_Widget *> widget_map;
Fl_Widget *get_widget(string_view name);
json list_widgets();

class WidgetData : public Fl_Callback_User_Data
{
    // Widget name, used for look-up, etc.
    string wname;
    // Widget type, which can be used to access a type's member
    // functions, also the name of the type.
    int wtype;
    // Substitute for Fl_Widget's user_data
    void *user_data = nullptr;
    bool auto_delete_user_data = false;

public:
    WidgetData(std::string_view type, std::string_view name, Fl_Widget *widget);
    ~WidgetData() override;

    void add_widget(Fl_Widget *w);
    void remove_widget(std::string_view name);

    string_view widget_name();
    int widget_type();
    string_view widget_type_name();
};

void gui_new(string_view name,
             string_view widget_type,
             string_view parent,
             json data);
void new_window(string_view name, string_view parent, json data);
void new_flex(string_view name, string_view parent, json data);
void new_grid(string_view name, string_view parent, json data);

#endif // FLTK_JSON_H
