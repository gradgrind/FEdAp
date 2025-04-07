#ifndef WIDGET_METHODS_H
#define WIDGET_METHODS_H

#include "fltk_json.h"

void widget_set_size(std::string_view name, json data);
void widget_set_box(std::string_view name, json data);
void widget_set_color(std::string_view name, json data);

#endif // WIDGET_METHODS_H
