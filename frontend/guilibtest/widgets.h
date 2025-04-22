#ifndef WIDGETS_H
#define WIDGETS_H

#include "minion.h"
#include <FL/Fl_Widget.H>

void choice_method(Fl_Widget *w, std::string_view c, minion::MinionList m);
void input_method(Fl_Widget *w, std::string_view c, minion::MinionList m);
void rowtable_method(Fl_Widget *w, std::string_view c, minion::MinionList m);
void editform_method(Fl_Widget *w, std::string_view c, minion::MinionList m);

Fl_Widget *NEW_Box(minion::MinionMap param);
Fl_Widget *NEW_Choice(minion::MinionMap param);
Fl_Widget *NEW_Output(minion::MinionMap param);
Fl_Widget *NEW_RowTable(minion::MinionMap param);
Fl_Widget *NEW_EditForm(minion::MinionMap param);

#endif // WIDGETS_H
