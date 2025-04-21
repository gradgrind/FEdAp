#ifndef WIDGET_BASE_H
#define WIDGET_BASE_H

#include "layout.h"
#include <FL/Fl_Widget.H>
#include <string>
#include <unordered_map>

//TODO: This is test code trying to trace the problems with table layout
// missing scrollbars).

struct Widget;
using w_method_handler = std::function<void(Widget, std::string_view, minion::MinionList)>;

struct Widget
{
    Fl_Widget *widget;
    std::string widget_type; //?
    w_method_handler handle_method;
};

extern std::unordered_map<std::string, Widget> widget_info;

void do_GUI(minion::MinionMap m);

Fl_Widget *w_NEW_Window(minion::MinionMap param);

void w_grid_method(Widget wd, std::string_view c, minion::MinionList m);
void w_flex_method(Widget wd, std::string_view c, minion::MinionList m);
void w_group_method(Widget wd, std::string_view c, minion::MinionList m);
void w_widget_method(Widget wd, std::string_view c, minion::MinionList m);

void w_choice_method(Widget wd, std::string_view c, minion::MinionList m);
void w_input_method(Widget wd, std::string_view c, minion::MinionList m);
void w_rowtable_method(Widget wd, std::string_view c, minion::MinionList m);

#endif // WIDGET_BASE_H
