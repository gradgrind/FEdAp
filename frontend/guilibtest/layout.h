#ifndef LAYOUT_H
#define LAYOUT_H

#include "minion.h"
#include <FL/Fl_Widget.H>
#include <functional>
#include <string_view>

Fl_Widget *NEW_Window(std::string_view name, minion::MinionList do_list);
Fl_Widget *NEW_Vlayout(std::string_view name, minion::MinionList do_list);
Fl_Widget *NEW_Hlayout(std::string_view name, minion::MinionList do_list);
Fl_Widget *NEW_Grid(std::string_view name, minion::MinionList do_list);

using method_handler = std::function<void(Fl_Widget *, std::string_view, minion::MinionList)>;
void grid_methods(Fl_Widget *w, std::string_view c, minion::MinionList m);
void flex_methods(Fl_Widget *w, std::string_view c, minion::MinionList m);
void group_methods(Fl_Widget *w, std::string_view c, minion::MinionList m);
void widget_methods(Fl_Widget *w, std::string_view c, minion::MinionList m);

#endif // LAYOUT_H
