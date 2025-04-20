#ifndef LAYOUT_H
#define LAYOUT_H

#include "minion.h"
#include <FL/Fl_Widget.H>
#include <functional>
#include <string_view>

//TODO ...
void tmp_run(minion::MinionMap data);

Fl_Widget *NEW_Window(minion::MinionMap param);
Fl_Widget *NEW_Vlayout(minion::MinionMap param);
Fl_Widget *NEW_Hlayout(minion::MinionMap param);
Fl_Widget *NEW_Grid(minion::MinionMap param);

using method_handler = std::function<void(Fl_Widget *, std::string_view, minion::MinionList)>;
void grid_method(Fl_Widget *w, std::string_view c, minion::MinionList m);
void flex_method(Fl_Widget *w, std::string_view c, minion::MinionList m);
void group_method(Fl_Widget *w, std::string_view c, minion::MinionList m);
void widget_method(Fl_Widget *w, std::string_view c, minion::MinionList m);

void do_callback(Fl_Widget *w, void *x);

class WidgetData : public Fl_Callback_User_Data
{
    static std::unordered_map<std::string_view, Fl_Widget *> widget_map;

    // Widget name, used for look-up, etc.
    std::string w_name;
    // Widget type, which can be used to access a type's member
    // functions, also the name of the type.
    //??? int wtype;
    // Substitute for Fl_Widget's user_data
    void *user_data = nullptr;
    bool auto_delete_user_data = false;

    WidgetData(std::string_view w_name, method_handler h);

public:
    ~WidgetData() override;

    static void add_widget(std::string_view name, Fl_Widget *w, method_handler h);
    static Fl_Widget *get_widget(std::string_view name);
    static minion::MinionList list_widgets();
    static std::string_view get_widget_name(Fl_Widget *w);

    method_handler handle_method;

    void remove_widget(std::string_view name);

    std::string_view widget_name();
    //int widget_type();
    //std::string_view widget_type_name();
};

#endif // LAYOUT_H
