#ifndef WIDGET_METHODS_H
#define WIDGET_METHODS_H

#include "fltk_json.h"

void widget_set_size(string_view name, json data);
void widget_set_box(string_view name, json data);
void widget_set_color(string_view name, json data);

#endif // WIDGET_METHODS_H
