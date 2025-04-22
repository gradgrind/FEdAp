#ifndef WIDGET_METHODS_H
#define WIDGET_METHODS_H

#include "minion.h"
#include <FL/Enumerations.H>
#include <FL/Fl_Widget.H>
#include <string>

Fl_Color get_colour(std::string &colour);
Fl_Boxtype get_boxtype(std::string &boxtype);

int int_param(minion::MinionList m, int i);
Fl_Color colour_param(minion::MinionList m, int i);

void left_label(Fl_Widget *w, minion::MinionList m);

#endif // WIDGET_METHODS_H
