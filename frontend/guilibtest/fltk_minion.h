#ifndef FLTK_MINION_H
#define FLTK_MINION_H

#include <FL/Fl_Widget.H>

#include <fmt/format.h>
#include <string>
#include <string_view>

#include "minion.h"
using mvalue = minion::MinionValue;
using mmap = minion::MinionMap;
using mlist = minion::MinionList;
using minion::read_minion;

#define MAGIC_ENUM_RANGE_MIN 0
#define MAGIC_ENUM_RANGE_MAX 255
#include "magic_enum/magic_enum.hpp"

//using method = std::function<void(Fl_Widget *, mmap)>;
//using member_map = std::map<std::string_view, method>;

// The user_data in Fl_Widget is used to store additional fields,
// primarily the widget's name and type. As the original user-data feature
// is now no longer directly available to the programmer, another field
// (user_data) is provided within the WidgetData subclass of
// Fl_Callback_User_Data which is referred to by the original user data field.

//extern unordered_map<string_view, Fl_Widget *> widget_map;
Fl_Widget *get_widget(std::string_view name);
mlist list_widgets();

class WidgetData : public Fl_Callback_User_Data
{
    // Widget name, used for look-up, etc.
    std::string wname;
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

    std::string_view widget_name();
    int widget_type();
    std::string_view widget_type_name();
    void do_method(Fl_Widget *widget, std::string_view name, mmap data);
};

bool get_minion_string(mmap data, std::string_view key, std::string &value);
std::string get_minion_string(mmap data, std::string_view key);
bool get_minion_int(mmap data, std::string_view key, int &value);
int get_minion_int(mmap data, std::string_view key);

void gui_new(std::string_view name,
             std::string_view widget_type,
             std::string_view parent,
             mmap data);

void new_window(std::string_view name, std::string_view parent, mmap data);
void new_flex(std::string_view name, std::string_view parent, mmap data);
void new_grid(std::string_view name, std::string_view parent, mmap data);
void new_box(std::string_view name, std::string_view parent, mmap data);

#endif // FLTK_MINION_H
