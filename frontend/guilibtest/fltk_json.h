#ifndef FLTK_JSON_H
#define FLTK_JSON_H

#include <FL/Fl_Widget.H>

#include <fmt/format.h>
#include <string>
#include <string_view>

#include <json.hpp>
using json = nlohmann::json;
using jobj = json::object_t;

#define MAGIC_ENUM_RANGE_MIN 0
#define MAGIC_ENUM_RANGE_MAX 255
#include "magic_enum/magic_enum.hpp"

using namespace std;

using method = function<void(Fl_Widget *, json)>;
using member_map = map<string_view, method>;

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
    WidgetData(string_view type, string_view name, Fl_Widget *widget);
    ~WidgetData() override;

    void add_widget(Fl_Widget *w);
    void remove_widget(string_view name);

    string_view widget_name();
    int widget_type();
    string_view widget_type_name();
    void do_method(Fl_Widget *widget, string_view name, json data);
};

bool get_json_string(json data, string_view key, string &value);
string get_json_string(json data, string_view key);
int get_json_int(json data, string_view key);

void gui_new(string_view name,
             string_view widget_type,
             string_view parent,
             json data);

void new_window(string_view name, string_view parent, json data);
void new_flex(string_view name, string_view parent, json data);
void new_grid(string_view name, string_view parent, json data);
void new_box(string_view name, string_view parent, json data);

#endif // FLTK_JSON_H
